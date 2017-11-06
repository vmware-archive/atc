port module SubPage exposing (Model(..), Msg(..), init, subscriptions, update, urlUpdate, view)

import Autoscroll
import Build
import Concourse
import Concourse.Pipeline
import Html exposing (Html)
import Http
import Job
import Json.Encode
import Login
import Resource
import BetaResource
import Build
import NoPipeline
import NotFound
import Pipeline
import QueryString
import Resource
import Routes
import String
import Task
import TeamSelection
import UpdateMsg exposing (UpdateMsg)


-- TODO: move ports somewhere else


port renderPipeline : ( Json.Encode.Value, Json.Encode.Value ) -> Cmd msg


port setTitle : String -> Cmd msg


type Model
    = WaitingModel Routes.ConcourseRoute
    | NoPipelineModel
    | BuildModel (Autoscroll.Model Build.Model)
    | JobModel Job.Model
    | ResourceModel Resource.Model
    | BetaResourceModel BetaResource.Model
    | LoginModel Login.Model
    | PipelineModel Pipeline.Model
    | SelectTeamModel TeamSelection.Model
    | NotFoundModel NotFound.Model


type Msg
    = PipelinesFetched (Result Http.Error (List Concourse.Pipeline))
    | DefaultPipelineFetched (Maybe Concourse.Pipeline)
    | NoPipelineMsg NoPipeline.Msg
    | BuildMsg (Autoscroll.Msg Build.Msg)
    | JobMsg Job.Msg
    | ResourceMsg Resource.Msg
    | BetaResourceMsg BetaResource.Msg
    | LoginMsg Login.Msg
    | PipelineMsg Pipeline.Msg
    | SelectTeamMsg TeamSelection.Msg
    | NewCSRFToken String


superDupleWrap : ( a -> b, c -> d ) -> ( a, Cmd c ) -> ( b, Cmd d )
superDupleWrap ( modelFunc, msgFunc ) ( model, msg ) =
    ( modelFunc model, Cmd.map msgFunc msg )


queryGroupsForRoute : Routes.ConcourseRoute -> List String
queryGroupsForRoute route =
    QueryString.all "groups" route.queries


init : String -> Routes.ConcourseRoute -> ( Model, Cmd Msg )
init turbulencePath route =
    case route.logical of
        Routes.Build teamName pipelineName jobName buildName ->
            superDupleWrap ( BuildModel, BuildMsg ) <|
                Autoscroll.init
                    Build.getScrollBehavior
                    << Build.init
                        { title = setTitle }
                        { csrfToken = "", hash = route.hash }
                <|
                    Build.JobBuildPage
                        { teamName = teamName
                        , pipelineName = pipelineName
                        , jobName = jobName
                        , buildName = buildName
                        }

        Routes.OneOffBuild buildId ->
            superDupleWrap ( BuildModel, BuildMsg ) <|
                Autoscroll.init
                    Build.getScrollBehavior
                    << Build.init
                        { title = setTitle }
                        { csrfToken = "", hash = route.hash }
                <|
                    Build.BuildPage <|
                        Result.withDefault 0 (String.toInt buildId)

        Routes.Resource teamName pipelineName resourceName ->
            superDupleWrap ( ResourceModel, ResourceMsg ) <|
                Resource.init
                    { title = setTitle }
                    { resourceName = resourceName
                    , teamName = teamName
                    , pipelineName = pipelineName
                    , paging = route.page
                    , csrfToken = ""
                    }

        Routes.BetaResource teamName pipelineName resourceName ->
            superDupleWrap ( BetaResourceModel, BetaResourceMsg ) <|
                BetaResource.init
                    { title = setTitle }
                    { resourceName = resourceName
                    , teamName = teamName
                    , pipelineName = pipelineName
                    , paging = route.page
                    , csrfToken = ""
                    }

        Routes.Job teamName pipelineName jobName ->
            superDupleWrap ( JobModel, JobMsg ) <|
                Job.init
                    { title = setTitle }
                    { jobName = jobName
                    , teamName = teamName
                    , pipelineName = pipelineName
                    , paging = route.page
                    , csrfToken = ""
                    }

        Routes.SelectTeam ->
            let
                redirect =
                    Maybe.withDefault "" <| QueryString.one QueryString.string "redirect" route.queries
            in
                superDupleWrap ( SelectTeamModel, SelectTeamMsg ) <|
                    TeamSelection.init { title = setTitle } redirect

        Routes.TeamLogin teamName ->
            superDupleWrap ( LoginModel, LoginMsg ) <|
                Login.init { title = setTitle } teamName (QueryString.one QueryString.string "redirect" route.queries)

        Routes.Pipeline teamName pipelineName ->
            superDupleWrap ( PipelineModel, PipelineMsg ) <|
                Pipeline.init
                    { render = renderPipeline
                    , title = setTitle
                    }
                    { teamName = teamName
                    , pipelineName = pipelineName
                    , turbulenceImgSrc = turbulencePath
                    , route = route
                    }

        Routes.Home ->
            ( WaitingModel route
            , Cmd.batch
                [ fetchPipelines
                , setTitle ""
                ]
            )


handleNotFound : String -> ( a -> Model, c -> Msg ) -> ( a, Cmd c, Maybe UpdateMsg ) -> ( Model, Cmd Msg )
handleNotFound notFound ( mdlFunc, msgFunc ) ( mdl, msg, outMessage ) =
    case outMessage of
        Just UpdateMsg.NotFound ->
            ( NotFoundModel { notFoundImgSrc = notFound }, setTitle "Not Found " )

        Nothing ->
            superDupleWrap ( mdlFunc, msgFunc ) <| ( mdl, msg )


update : String -> String -> Concourse.CSRFToken -> Msg -> Model -> ( Model, Cmd Msg )
update turbulence notFound csrfToken msg mdl =
    case ( msg, mdl ) of
        ( NoPipelineMsg msg, model ) ->
            ( model, fetchPipelines )

        ( NewCSRFToken c, BuildModel scrollModel ) ->
            let
                buildModel =
                    scrollModel.subModel

                ( newBuildModel, buildCmd ) =
                    Build.update (Build.NewCSRFToken c) buildModel
            in
                ( BuildModel { scrollModel | subModel = newBuildModel }, buildCmd |> Cmd.map (\buildMsg -> BuildMsg (Autoscroll.SubMsg buildMsg)) )

        ( BuildMsg message, BuildModel scrollModel ) ->
            let
                subModel =
                    scrollModel.subModel

                model =
                    { scrollModel | subModel = { subModel | csrfToken = csrfToken } }
            in
                handleNotFound notFound ( BuildModel, BuildMsg ) (Autoscroll.update Build.updateWithMessage message model)

        ( NewCSRFToken c, JobModel model ) ->
            ( JobModel { model | csrfToken = c }, Cmd.none )

        ( JobMsg message, JobModel model ) ->
            handleNotFound notFound ( JobModel, JobMsg ) (Job.updateWithMessage message { model | csrfToken = csrfToken })

        ( LoginMsg message, LoginModel model ) ->
            let
                ( mdl, msg ) =
                    Login.update message model

                --            superDupleWrap ( LoginModel, LoginMsg ) <| Login.update message model
            in
                ( LoginModel mdl, Cmd.map LoginMsg msg )

        ( PipelineMsg message, PipelineModel model ) ->
            handleNotFound notFound ( PipelineModel, PipelineMsg ) (Pipeline.updateWithMessage message model)

        ( NewCSRFToken c, ResourceModel model ) ->
            ( ResourceModel { model | csrfToken = c }, Cmd.none )

        ( ResourceMsg message, ResourceModel model ) ->
            handleNotFound notFound ( ResourceModel, ResourceMsg ) (Resource.updateWithMessage message { model | csrfToken = csrfToken })

        ( BetaResourceMsg message, BetaResourceModel model ) ->
            handleNotFound notFound ( BetaResourceModel, BetaResourceMsg ) (BetaResource.updateWithMessage message { model | csrfToken = csrfToken })

        ( SelectTeamMsg message, SelectTeamModel model ) ->
            superDupleWrap ( SelectTeamModel, SelectTeamMsg ) <| TeamSelection.update message model

        ( DefaultPipelineFetched pipeline, WaitingModel route ) ->
            case pipeline of
                Nothing ->
                    ( NoPipelineModel, setTitle "" )

                Just p ->
                    let
                        flags =
                            { teamName = p.teamName
                            , pipelineName = p.name
                            , turbulenceImgSrc = turbulence
                            , route = route
                            }
                    in
                        superDupleWrap ( PipelineModel, PipelineMsg ) <| Pipeline.init { render = renderPipeline, title = setTitle } flags

        ( DefaultPipelineFetched _, NoPipelineModel ) ->
            ( mdl, Cmd.none )

        ( NewCSRFToken _, _ ) ->
            ( mdl, Cmd.none )

        unknown ->
            flip always (Debug.log ("impossible combination") unknown) <|
                ( mdl, Cmd.none )


urlUpdate : Routes.ConcourseRoute -> Model -> ( Model, Cmd Msg )
urlUpdate route model =
    case ( route.logical, model ) of
        ( Routes.Pipeline team pipeline, PipelineModel mdl ) ->
            superDupleWrap ( PipelineModel, PipelineMsg ) <|
                Pipeline.changeToPipelineAndGroups
                    { teamName = team
                    , pipelineName = pipeline
                    , turbulenceImgSrc = mdl.turbulenceImgSrc
                    , route = route
                    }
                    mdl

        ( Routes.Resource teamName pipelineName resourceName, ResourceModel mdl ) ->
            superDupleWrap ( ResourceModel, ResourceMsg ) <|
                Resource.changeToResource
                    { teamName = teamName
                    , pipelineName = pipelineName
                    , resourceName = resourceName
                    , paging = route.page
                    , csrfToken = mdl.csrfToken
                    }
                    mdl

        ( Routes.BetaResource teamName pipelineName resourceName, BetaResourceModel mdl ) ->
            superDupleWrap ( BetaResourceModel, BetaResourceMsg ) <|
                BetaResource.changeToResource
                    { teamName = teamName
                    , pipelineName = pipelineName
                    , resourceName = resourceName
                    , paging = route.page
                    , csrfToken = mdl.csrfToken
                    }
                    mdl

        ( Routes.Job teamName pipelineName jobName, JobModel mdl ) ->
            superDupleWrap ( JobModel, JobMsg ) <|
                Job.changeToJob
                    { teamName = teamName
                    , pipelineName = pipelineName
                    , jobName = jobName
                    , paging = route.page
                    , csrfToken = mdl.csrfToken
                    }
                    mdl

        ( Routes.Build teamName pipelineName jobName buildName, BuildModel scrollModel ) ->
            let
                ( submodel, subcmd ) =
                    Build.changeToBuild
                        (Build.JobBuildPage
                            { teamName = teamName
                            , pipelineName = pipelineName
                            , jobName = jobName
                            , buildName = buildName
                            }
                        )
                        scrollModel.subModel
            in
                ( BuildModel { scrollModel | subModel = submodel }
                , Cmd.map BuildMsg (Cmd.map Autoscroll.SubMsg subcmd)
                )

        _ ->
            ( model, Cmd.none )


view : Model -> Html Msg
view mdl =
    case mdl of
        BuildModel model ->
            Html.map BuildMsg <| Autoscroll.view Build.view model

        JobModel model ->
            Html.map JobMsg <| Job.view model

        LoginModel model ->
            Html.map LoginMsg <| Login.view model

        PipelineModel model ->
            Html.map PipelineMsg <| Pipeline.view model

        ResourceModel model ->
            Html.map ResourceMsg <| Resource.view model

        BetaResourceModel model ->
            Html.map BetaResourceMsg <| BetaResource.view model

        SelectTeamModel model ->
            Html.map SelectTeamMsg <| TeamSelection.view model

        WaitingModel _ ->
            Html.div [] []

        NoPipelineModel ->
            Html.map NoPipelineMsg <| NoPipeline.view

        NotFoundModel model ->
            NotFound.view model


subscriptions : Model -> Sub Msg
subscriptions mdl =
    case mdl of
        BuildModel model ->
            Sub.map BuildMsg <| Autoscroll.subscriptions Build.subscriptions model

        JobModel model ->
            Sub.map JobMsg <| Job.subscriptions model

        LoginModel model ->
            Sub.map LoginMsg <| Login.subscriptions model

        NoPipelineModel ->
            Sub.map NoPipelineMsg <| NoPipeline.subscriptions

        PipelineModel model ->
            Sub.map PipelineMsg <| Pipeline.subscriptions model

        ResourceModel model ->
            Sub.map ResourceMsg <| Resource.subscriptions model

        BetaResourceModel model ->
            Sub.map BetaResourceMsg <| BetaResource.subscriptions model

        SelectTeamModel model ->
            Sub.map SelectTeamMsg <| TeamSelection.subscriptions model

        WaitingModel _ ->
            Sub.none

        NotFoundModel _ ->
            Sub.none


fetchPipelines : Cmd Msg
fetchPipelines =
    Task.attempt PipelinesFetched Concourse.Pipeline.fetchPipelines
