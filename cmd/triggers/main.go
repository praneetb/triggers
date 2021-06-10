package main

import (
	"github.com/praneetb/triggers/pkgs/alcon"
	"github.com/spf13/cobra"
)

// Execute - start the controller
func main() {
	cobra.CheckErr(alcon.Triggers.Execute())
}
