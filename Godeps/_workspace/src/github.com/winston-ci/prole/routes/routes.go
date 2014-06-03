package routes

import "github.com/tedsuo/router"

const (
	ExecuteBuild = "ExecuteBuild"
	AbortBuild   = "AbortBuild"
	CheckInput   = "CheckInput"
)

var Routes = router.Routes{
	{Path: "/builds", Method: "POST", Handler: ExecuteBuild},
	{Path: "/builds/:guid/abort", Method: "POST", Handler: AbortBuild},
	{Path: "/checks", Method: "POST", Handler: CheckInput},
}
