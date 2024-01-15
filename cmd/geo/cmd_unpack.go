package main

import (
	"github.com/metacubex/geo/cmd/geo/internal/convert"
	"github.com/spf13/cobra"
)

func init() {
	// commandUnpack.AddCommand(convert.CommandIP)
	commandUnpack.AddCommand(convert.CommandUnpackSite)
	mainCommand.AddCommand(commandUnpack)
}

var commandUnpack = &cobra.Command{
	Use:   "unpack",
	Short: "unpack geo resource",
}
