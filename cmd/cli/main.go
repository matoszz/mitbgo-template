package main

import (
	template "github.com/datumforge/go-template/cmd/cli/cmd"

	// since the cmds are no longer part of the same package
	// they must all be imported in main
	_ "github.com/datumforge/go-template/cmd/cli/cmd/version"
)

func main() {
	template.Execute()
}
