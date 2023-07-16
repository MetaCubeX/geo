package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/metacubex/geo/geoip"

	"github.com/spf13/cobra"
)

func init() {
	commandCheck.PersistentFlags().StringVarP(&dbType, "type", "t", "", "specific database type")
	commandCheck.PersistentFlags().StringVarP(&dbPath, "file", "f", "", "specific database file path")
	mainCommand.AddCommand(commandCheck)
}

var commandCheck = &cobra.Command{
	Use:   "check",
	Short: "Check geo resources availability",
	RunE:  check,
}

var descriptionPlaceholder = map[string]string{"PLACEHOLDER": "geo"}

func check(cmd *cobra.Command, args []string) error {
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

	for _, filePath := range paths {
		fmt.Println("🔎Checking", filePath)

		var db *geoip.Database
		db, err = geoip.FromFile(filePath)
		if err != nil {
			return err
		}
		mmdb := db.Reader()
		if len(mmdb.Metadata.Description) == 0 {
			// fix verify failed when description is empty
			mmdb.Metadata.Description = descriptionPlaceholder
		}

		err = mmdb.Verify()
		if err == nil {
			fmt.Println("👌Successfully verified geo database!")
		} else {
			fmt.Println("❌Failed to verify geo database!")
			fmt.Println("Error:", err)
		}
		if mmdb.Metadata.Description["PLACEHOLDER"] == "geo" {
			mmdb.Metadata.Description = nil
		}

		fmt.Print("🔢Type: ")
		switch db.SourceType {
		case geoip.TypeMaxmind, geoip.TypeSing, geoip.TypeMetaV0:
			fmt.Println(db.SourceType, "MMDB database")
			fmt.Println("⏰MMDB build time:", time.Unix(int64(mmdb.Metadata.BuildEpoch), 0))
			fmt.Println("📃MMDB metadata:")
			mmdbJSON, _ := json.MarshalIndent(mmdb.Metadata, "", "    ")
			os.Stdout.Write(mmdbJSON)
			os.Stdout.WriteString("\n")
		case geoip.TypeV2Ray:
			fmt.Println("V2Ray GeoIP database")
			fmt.Println("📒Total nodes:", mmdb.Metadata.NodeCount)
		default:
			fmt.Println("unknown:", mmdb.Metadata.DatabaseType)
		}
		os.Stdout.WriteString("\n")
	}
	return nil
}
