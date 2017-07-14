package commands

type CredhubCommand struct {
	Api ApiCommand `command:"api" alias:"a" description:"Set the CredHub API target to be used for subsequent commands"`

	Get        GetCommand        `command:"get" alias:"g" description:"Get a credential value"`
	Set        SetCommand        `command:"set" alias:"s" description:"Set a credential with a provided value"`
	Generate   GenerateCommand   `command:"generate" alias:"n" description:"Set a credential with a generated value"`
	Regenerate RegenerateCommand `command:"regenerate" alias:"r" description:"Set a credential with a generated value using the same attributes as the stored value"`
	Delete     DeleteCommand     `command:"delete" alias:"d" description:"Delete a credential value"`
	Login      LoginCommand      `command:"login" alias:"l" description:"Authenticate user with CredHub"`
	Logout     LogoutCommand     `command:"logout" alias:"o" description:"Discard authenticated user session"`
	Find       FindCommand       `command:"find" alias:"f" description:"Find stored credentials based on query parameters"`
	Import     ImportCommand     `command:"import" alias:"i" description:"Set multiple credential values"`
	Version    func()            `long:"version" description:"Version of CLI and targeted CredHub API"`
	Token      func()            `long:"token" description:"Return your current CredHub authorization token"`
}

var CredHub CredhubCommand
