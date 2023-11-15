package convert

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/metacubex/geo/convert"
	"github.com/metacubex/geo/encoding/v2raygeo"

	E "github.com/sagernet/sing/common/exceptions"
	"github.com/spf13/cobra"
)

func init() {
	CommandSite.PersistentFlags().StringVarP(&fromType, "from-type", "i", "", "specify input database type")
	CommandSite.PersistentFlags().StringVarP(&toType, "to-type", "o", "meta", "set output database type")
	CommandSite.PersistentFlags().StringVarP(&output, "output-name", "f", "", "specify output filename")
	CommandSite.PersistentFlags().StringVarP(&code, "code", "c", "", "specify output code")
}

var CommandSite = &cobra.Command{
	Use:   "site",
	Short: "Convert GeoSite resources",
	RunE:  site,
	Args:  cobra.ExactArgs(1),
}

func site(cmd *cobra.Command, args []string) error {
	var (
		buffer   bytes.Buffer
		filename = time.Now().Format("2006-01-02 15-04-05 -07 MST")
	)
	fmt.Println("➕Loading file:", args[0])
	fileContent, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	buffer.Grow(8 * 1024 * 1024) // 8 MiB
	fmt.Println("🔁Converting GeoSite database:", fromType, "->", toType)
	startTime := time.Now()

	switch strings.ToLower(fromType) {
	case "v2ray":
		var geositeList []*v2raygeo.GeoSite
		geositeList, err = v2raygeo.LoadSite(fileContent)
		if err != nil {
			return err
		}
		switch strings.ToLower(toType) {
		case "sing", "sing-geosite":
			err = convert.V2RaySiteToSing(geositeList, &buffer)
			if err != nil {
				return err
			}
			filename += ".db"
		case "clash":
			err = convert.V2RayToYamlByCode(geositeList, &buffer, code)
			if err != nil {
				return err
			}
			filename += ".yaml"
		default:
			return E.New("unsupported output GeoSite database type: ", toType)
		}

	default:
		return E.New("unsupported input GeoSite database type: ", toType)
	}

	if output != "" {
		filename = output
	}
	err = os.WriteFile(filename, buffer.Bytes(), 0o666)
	if err != nil {
		return err
	}
	fmt.Println("🎉Successfully converted to", filename, "in", time.Now().Sub(startTime))
	return nil
}
