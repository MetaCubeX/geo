package main

import (
	"github.com/metacubex/geo/cmd/geo/internal/unpack"

	"github.com/spf13/cobra"
)

func init() {
	commandUnpack.AddCommand(unpack.CommandSite)
	mainCommand.AddCommand(commandUnpack)
}

var commandUnpack = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack geo resources",
}
