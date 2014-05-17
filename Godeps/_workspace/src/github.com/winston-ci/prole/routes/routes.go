package routes

import "github.com/tedsuo/router"

const (
	ExecuteBuild = "ExecuteBuild"
	CheckInput   = "CheckInput"
)

var Routes = router.Routes{
	{Path: "/builds", Method: "POST", Handler: ExecuteBuild},
	{Path: "/checks", Method: "POST", Handler: CheckInput},
}
