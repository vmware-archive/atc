package models

type infoApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type authServer struct {
	Url    string `json:"url"`
	Client string `json:"client"`
}

type Info struct {
	App        infoApp    `json:"app"`
	AuthServer authServer `json:"auth-server"`
}
