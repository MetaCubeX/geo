package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/metacubex/geo"

	F "github.com/sagernet/sing/common/format"
	"github.com/spf13/cobra"
)

var (
	workingDir string
	currentDir string
)

var mainCommand = &cobra.Command{
	Use:               "geo",
	PersistentPreRunE: preRun,
	Long: F.ToString("geo v", geo.Version,
		" (", runtime.Version(), ", ", runtime.GOOS, "/", runtime.GOARCH, ")"),
}

func init() {
	mainCommand.PersistentFlags().StringVarP(&workingDir, "directory", "D", "", "set working directory")
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func preRun(cmd *cobra.Command, args []string) error {
	if workingDir == "" {
		var err error
		workingDir, err = os.UserHomeDir()
		if err != nil {
			workingDir, _ = os.Getwd()
		}
		workingDir = path.Join(workingDir, ".geo")
		return os.MkdirAll(workingDir, 0o750)
	} else {
		if !filepath.IsAbs(workingDir) {
			currentDir, _ := os.Getwd()
			workingDir = filepath.Join(currentDir, workingDir)
		}
		return os.Chdir(workingDir)
	}
}
