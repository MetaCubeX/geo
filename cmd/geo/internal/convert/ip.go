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
	CommandIP.PersistentFlags().StringVarP(&fromType, "from-type", "i", "", "specify input database type")
	CommandIP.PersistentFlags().StringVarP(&toType, "to-type", "o", "meta", "set output database type")
	CommandIP.PersistentFlags().StringVarP(&output, "output-name", "f", "", "specify output filename")
}

var (
	fromType string
	toType   string
	output   string
)

var CommandIP = &cobra.Command{
	Use:   "ip",
	Short: "Convert GeoIP resources",
	RunE:  ip,
	Args:  cobra.ExactArgs(1),
}

func ip(cmd *cobra.Command, args []string) error {
	var (
		buffer   bytes.Buffer
		filename = time.Now().Format("2006-01-02 15-04-05 MST")
	)
	fmt.Println("➕Loading file:", args[0])
	fileContent, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	buffer.Grow(8 * 1024 * 1024) // 8 MiB
	fmt.Println("🔁Converting GeoIP database:", fromType, "->", toType)
	startTime := time.Now()

	switch strings.ToLower(fromType) {
	case "maxmind":
		switch strings.ToLower(toType) {
		case "sing", "sing-geoip":
			err = convert.MaxMindIPToSing(fileContent, &buffer)
			if err != nil {
				return err
			}
			filename += ".db"

		//case "meta", "meta0", "meta-geoip0":
		//	err = convert.MaxMindIPToMetaV0(fileContent, &buffer)
		//	if err != nil {
		//		return err
		//	}
		//	filename += ".metadb"

		default:
			return E.New("unsupported output GeoIP database type: ", toType)
		}

	case "sing", "sing-geoip":
		switch strings.ToLower(toType) {
		//case "meta", "meta0", "meta-geoip0":
		//	err = convert.SingIPToMetaV0(fileContent, &buffer)
		//	if err != nil {
		//		return err
		//	}
		//	filename += ".metadb"

		default:
			return E.New("unsupported output GeoIP database type: ", toType)
		}

	case "v2ray":
		var geoipList []*v2raygeo.GeoIP
		geoipList, err = v2raygeo.LoadIP(fileContent)
		if err != nil {
			return err
		}
		switch strings.ToLower(toType) {
		case "sing", "sing-geoip":
			err = convert.V2RayIPToSing(geoipList, &buffer)
			if err != nil {
				return err
			}
			filename += ".db"

		case "meta", "meta0", "meta-geoip0":
			err = convert.V2RayIPToMetaV0(geoipList, &buffer)
			if err != nil {
				return err
			}
			filename += ".metadb"

		default:
			return E.New("unsupported output GeoIP database type: ", toType)
		}

	case "meta", "meta0", "meta-geoip0":
		switch strings.ToLower(toType) {
		case "sing", "sing-geoip":
			err = convert.MetaV0ToSing(fileContent, &buffer)
			if err != nil {
				return err
			}
			filename += ".db"

		default:
			return E.New("unsupported output GeoIP database type: ", toType)
		}
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
