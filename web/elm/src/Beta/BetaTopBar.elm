port module BetaTopBar exposing (Model, Msg(..), fetchUser, init, subscriptions, update, urlUpdate, view)

import BetaRoutes
import Concourse
import Concourse.Pipeline
import Concourse.User
import Format exposing (prependBeta)
import Html exposing (Html)
import Html.Attributes exposing (attribute, class, classList, disabled, href, id, style)
import Html.Events exposing (onClick)
import Http
import List
import LoginRedirect
import Navigation exposing (Location)
import QueryString
import StrictEvents exposing (onLeftClickOrShiftLeftClick)
import String
import Task
import Time


type alias Model =
    { route : BetaRoutes.ConcourseRoute
    , selectedGroups : List String
    , pipeline : Maybe Concourse.Pipeline
    , userState : UserState
    , userMenuVisible : Bool
    }


type UserState
    = UserStateLoggedIn Concourse.User
    | UserStateLoggedOut
    | UserStateUnknown


type Msg
    = Noop
    | PipelineFetched (Result Http.Error Concourse.Pipeline)
    | UserFetched (Result Http.Error Concourse.User)
    | FetchPipeline Concourse.PipelineIdentifier
    | ToggleSidebar
    | ToggleGroup Concourse.PipelineGroup
    | SetGroups (List String)
    | LogOut
    | LogIn
    | NavTo String
    | LoggedOut (Result Http.Error ())
    | ToggleUserMenu


queryGroupsForRoute : BetaRoutes.ConcourseRoute -> List String
queryGroupsForRoute route =
    QueryString.all "groups" route.queries


init : BetaRoutes.ConcourseRoute -> ( Model, Cmd Msg )
init route =
    let
        pid =
            extractPidFromRoute route.logical
    in
        ( { selectedGroups = queryGroupsForRoute route
          , route = route
          , pipeline = Nothing
          , userState = UserStateUnknown
          , userMenuVisible = False
          }
        , case pid of
            Nothing ->
                fetchUser

            Just pid ->
                Cmd.batch [ fetchPipeline pid, fetchUser ]
        )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Noop ->
            ( model, Cmd.none )

        FetchPipeline pid ->
            ( model, fetchPipeline pid )

        UserFetched (Ok user) ->
            ( { model | userState = UserStateLoggedIn user }
            , Cmd.none
            )

        UserFetched (Err _) ->
            ( { model | userState = UserStateLoggedOut }
            , Cmd.none
            )

        PipelineFetched (Ok pipeline) ->
            ( { model
                | pipeline = Just pipeline
              }
            , Cmd.none
            )

        PipelineFetched (Err err) ->
            case err of
                Http.BadStatus { status } ->
                    if status.code == 401 then
                        case model.route.logical of
                            BetaRoutes.BetaSelectTeam ->
                                ( model, Cmd.none )

                            BetaRoutes.BetaTeamLogin _ ->
                                ( model, Cmd.none )

                            _ ->
                                ( model, LoginRedirect.requestLoginRedirect "" )
                    else
                        ( model, Cmd.none )

                _ ->
                    ( model, Cmd.none )

        ToggleSidebar ->
            flip always (Debug.log "sidebar-toggle-incorrectly-handled" ()) <|
                ( model, Cmd.none )

        ToggleGroup group ->
            setGroups (toggleGroup group model.selectedGroups model.pipeline) model

        SetGroups groups ->
            setGroups groups model

        LogIn ->
            ( { model
                | selectedGroups = []
                , pipeline = Nothing
              }
            , Navigation.newUrl <| prependBeta "/login"
            )

        LogOut ->
            ( model, logOut )

        LoggedOut (Ok _) ->
            ( { model
                | userState = UserStateLoggedOut
                , pipeline = Nothing
                , selectedGroups = []
              }
            , Navigation.newUrl <| prependBeta "/"
            )

        NavTo url ->
            ( model, Navigation.newUrl url )

        LoggedOut (Err err) ->
            flip always (Debug.log "failed to log out" err) <|
                ( model, Cmd.none )

        ToggleUserMenu ->
            ( { model | userMenuVisible = not model.userMenuVisible }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    case pipelineIdentifierFromRouteOrModel model.route model of
        Nothing ->
            Sub.none

        Just pid ->
            Time.every (5 * Time.second) (always (FetchPipeline pid))


pipelineIdentifierFromRouteOrModel : BetaRoutes.ConcourseRoute -> Model -> Maybe Concourse.PipelineIdentifier
pipelineIdentifierFromRouteOrModel route model =
    case extractPidFromRoute route.logical of
        Nothing ->
            case model.pipeline of
                Nothing ->
                    Nothing

                Just pipeline ->
                    Just { teamName = pipeline.teamName, pipelineName = pipeline.name }

        Just pidFromRoute ->
            Just pidFromRoute


extractPidFromRoute : BetaRoutes.Route -> Maybe Concourse.PipelineIdentifier
extractPidFromRoute route =
    case route of
        BetaRoutes.BetaBuild teamName pipelineName jobName buildName ->
            Just { teamName = teamName, pipelineName = pipelineName }

        BetaRoutes.BetaJob teamName pipelineName jobName ->
            Just { teamName = teamName, pipelineName = pipelineName }

        BetaRoutes.BetaResource teamName pipelineName resourceName ->
            Just { teamName = teamName, pipelineName = pipelineName }

        BetaRoutes.BetaOneOffBuild buildId ->
            Nothing

        BetaRoutes.BetaPipeline teamName pipelineName ->
            Just { teamName = teamName, pipelineName = pipelineName }

        BetaRoutes.BetaSelectTeam ->
            Nothing

        BetaRoutes.BetaTeamLogin teamName ->
            Nothing

        BetaRoutes.BetaHome ->
            Nothing

        BetaRoutes.Dashboard ->
            Nothing


setGroups : List String -> Model -> ( Model, Cmd Msg )
setGroups newGroups model =
    let
        newUrl =
            pidToUrl (pipelineIdentifierFromRouteOrModel model.route model) <|
                setGroupsInLocation model.route newGroups
    in
        ( model, Navigation.newUrl newUrl )


urlUpdate : BetaRoutes.ConcourseRoute -> Model -> ( Model, Cmd Msg )
urlUpdate route model =
    let
        pipelineIdentifier =
            pipelineIdentifierFromRouteOrModel route model
    in
        ( { model
            | route = route
            , selectedGroups = queryGroupsForRoute route
          }
        , case pipelineIdentifier of
            Nothing ->
                fetchUser

            Just pid ->
                Cmd.batch [ fetchPipeline pid, fetchUser ]
        )


getDefaultSelectedGroups : Maybe Concourse.Pipeline -> List String
getDefaultSelectedGroups pipeline =
    case pipeline of
        Nothing ->
            []

        Just pipeline ->
            case List.head pipeline.groups of
                Nothing ->
                    []

                Just first ->
                    [ first.name ]


setGroupsInLocation : BetaRoutes.ConcourseRoute -> List String -> BetaRoutes.ConcourseRoute
setGroupsInLocation loc groups =
    let
        updatedUrl =
            if List.isEmpty groups then
                QueryString.remove "groups" loc.queries
            else
                List.foldr
                    (QueryString.add "groups")
                    QueryString.empty
                    groups
    in
        { loc
            | queries = updatedUrl
        }


pidToUrl : Maybe Concourse.PipelineIdentifier -> BetaRoutes.ConcourseRoute -> String
pidToUrl pid { queries } =
    case pid of
        Just { teamName, pipelineName } ->
            String.join ""
                [ String.join "/"
                    [ "/teams"
                    , teamName
                    , "pipelines"
                    , pipelineName
                    ]
                , QueryString.render queries
                ]

        Nothing ->
            ""


toggleGroup : Concourse.PipelineGroup -> List String -> Maybe Concourse.Pipeline -> List String
toggleGroup grp names mpipeline =
    if List.member grp.name names then
        List.filter ((/=) grp.name) names
    else if List.isEmpty names then
        grp.name :: getDefaultSelectedGroups mpipeline
    else
        grp.name :: names


getSelectedOrDefaultGroups : Model -> List String
getSelectedOrDefaultGroups model =
    if List.isEmpty model.selectedGroups then
        getDefaultSelectedGroups model.pipeline
    else
        model.selectedGroups


getSelectedGroupsForRoute : Model -> List String
getSelectedGroupsForRoute model =
    case model.route.logical of
        BetaRoutes.BetaBuild _ _ jobName _ ->
            getGroupsForJob jobName model.pipeline

        BetaRoutes.BetaJob _ _ jobName ->
            getGroupsForJob jobName model.pipeline

        _ ->
            getSelectedOrDefaultGroups model


getGroupsForJob : String -> Maybe Concourse.Pipeline -> List String
getGroupsForJob jobName pipeline =
    case pipeline of
        Nothing ->
            []

        Just pl ->
            List.filter (.jobs >> List.member jobName) pl.groups |> List.map .name


view : Model -> Html Msg
view model =
    Html.nav
        [ classList
            [ ( "top-bar", True )
            , ( "test", True )
            , ( "paused", isPaused model.pipeline )
            ]
        ]
        [ let
            ( groupList, pipelineUrl ) =
                case model.pipeline of
                    Nothing ->
                        ( [], "/" )

                    Just pipeline ->
                        ( List.map
                            (viewGroup (getSelectedGroupsForRoute model) pipeline.url)
                            pipeline.groups
                        , pipeline.url
                        )
          in
            Html.ul [ class "groups" ] <|
                [ Html.li [ class "main" ]
                    [ Html.span
                        [ class "sidebar-toggle test btn-hamburger"
                        , onClick ToggleSidebar
                        , Html.Attributes.attribute "aria-label" "Toggle List of Pipelines"
                        ]
                        [ Html.i [ class "fa fa-bars" ] []
                        ]
                    ]
                , Html.li [ class "main" ]
                    [ Html.a
                        [ StrictEvents.onLeftClick <| NavTo <| prependBeta pipelineUrl
                        , Html.Attributes.href <| prependBeta pipelineUrl
                        ]
                        [ Html.i [ class "fa fa-home" ] []
                        ]
                    ]
                ]
                    ++ groupList
        , Html.ul [ class "nav-right" ]
            [ Html.li [ class "nav-item" ]
                [ viewUserState model.userState model.userMenuVisible
                ]
            ]
        ]


isPaused : Maybe Concourse.Pipeline -> Bool
isPaused =
    Maybe.withDefault False << Maybe.map .paused


viewUserState : UserState -> Bool -> Html Msg
viewUserState userState userMenuVisible =
    case userState of
        UserStateUnknown ->
            Html.text ""

        UserStateLoggedOut ->
            Html.div [ class "user-info" ]
                [ Html.a
                    [ StrictEvents.onLeftClick <| LogIn
                    , href <| prependBeta "/login"
                    , Html.Attributes.attribute "aria-label" "Log In"
                    , class "login-button"
                    ]
                    [ Html.text "login"
                    ]
                ]

        UserStateLoggedIn { team } ->
            Html.div [ class "user-info" ]
                [ Html.div [ class "user-id", onClick ToggleUserMenu ]
                    [ Html.i [ class "fa fa-user" ] []
                    , Html.text " "
                    , Html.text team.name
                    , Html.text " "
                    , Html.i [ class "fa fa-caret-down" ] []
                    ]
                , Html.div [ classList [ ( "user-menu", True ), ( "hidden", not userMenuVisible ) ] ]
                    [ Html.a
                        [ Html.Attributes.attribute "aria-label" "Log Out"
                        , onClick LogOut
                        ]
                        [ Html.text "logout"
                        ]
                    ]
                ]


viewGroup : List String -> String -> Concourse.PipelineGroup -> Html Msg
viewGroup selectedGroups url grp =
    Html.li
        [ if List.member grp.name selectedGroups then
            class "main active"
          else
            class "main"
        ]
        [ Html.a
            [ Html.Attributes.href <| url ++ "?groups=" ++ grp.name
            , onLeftClickOrShiftLeftClick (SetGroups [ grp.name ]) (ToggleGroup grp)
            ]
            [ Html.text grp.name ]
        ]


fetchPipeline : Concourse.PipelineIdentifier -> Cmd Msg
fetchPipeline pipelineIdentifier =
    Task.attempt PipelineFetched <|
        Concourse.Pipeline.fetchPipeline pipelineIdentifier


fetchUser : Cmd Msg
fetchUser =
    Task.attempt UserFetched Concourse.User.fetchUser


logOut : Cmd Msg
logOut =
    Task.attempt LoggedOut Concourse.User.logOut
