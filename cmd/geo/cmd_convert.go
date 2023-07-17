package main

import (
	"github.com/metacubex/geo/cmd/geo/internal/convert"

	"github.com/spf13/cobra"
)

func init() {
	commandConvert.AddCommand(convert.CommandIP)
	commandConvert.AddCommand(convert.CommandSite)
	mainCommand.AddCommand(commandConvert)
}

var commandConvert = &cobra.Command{
	Use:   "convert",
	Short: "Convert geo resource encodings",
}
