package routes

import "github.com/tedsuo/router"

const (
	ExecuteBuild = "ExecuteBuild"
)

var Routes = router.Routes{
	{Path: "/builds", Method: "POST", Handler: ExecuteBuild},
}
