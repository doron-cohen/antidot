package main

import (
	"github.com/doron-cohen/antidot/cmd"
)

/*
We are injecting the app's version to this variable.
This is done using: go build -ldflags "-X 'main.version=$VERSION".
It is only used in the cmd package but I couldn't make it inject there.
*/
var version = "dev"

func main() {
	cmd.Execute(version)
}
