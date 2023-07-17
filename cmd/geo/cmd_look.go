package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/metacubex/geo/geoip"

	F "github.com/sagernet/sing/common/format"
	"github.com/spf13/cobra"
)

func init() {
	commandLook.PersistentFlags().StringVarP(&dbType, "type", "t", "", "specify database type")
	commandLook.PersistentFlags().StringVarP(&dbPath, "file", "f", "", "specify database file path")
	commandLook.PersistentFlags().BoolVarP(&immediate, "immediate", "i", false, "return immediately as soon as a result is found")
	mainCommand.AddCommand(commandLook)
}

var commandLook = &cobra.Command{
	Use:   "look",
	Short: "Query geo information from databases",
	RunE:  look,
	Args:  cobra.ExactArgs(1),
}

var immediate bool

func look(cmd *cobra.Command, args []string) error {
	var (
		paths []string
		err   error
	)
	if dbPath == "" {
		paths, err = find()
		if err != nil {
			return err
		}
	} else {
		paths = []string{dbPath}
	}

	fmt.Println("🔎Querying from", paths)
	result := make(map[string]struct{})
	startTime := time.Now()
	for _, filePath := range paths {
		var db *geoip.Database
		db, err = geoip.FromFile(filePath)
		if err != nil {
			return err
		}

		codes := db.LookupCode(net.ParseIP(args[0]))
		for _, code := range codes {
			result[code] = struct{}{}
		}

		if immediate && len(codes) > 0 {
			break
		}
	}

	os.Stdout.WriteString(F.ToString("🎉Query finished in ", time.Now().Sub(startTime), "!\n"))
	fmt.Print("Total ", len(result), " results (GeoIP codes):\n  ")
	for code := range result {
		os.Stdout.WriteString(code)
		os.Stdout.WriteString(" ")
	}
	return nil
}
