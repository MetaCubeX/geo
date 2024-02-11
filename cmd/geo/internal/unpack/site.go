package unpack

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/metacubex/geo/encoding/v2raygeo"

	"github.com/sagernet/sing/common"
	"github.com/spf13/cobra"
)

func init() {
	CommandSite.PersistentFlags().StringVarP(&code, "code", "c", "", "specify output code")
	CommandSite.PersistentFlags().StringVarP(&outDir, "out-dir", "d", "", "specify output directory")
}

var CommandSite = &cobra.Command{
	Use:   "site",
	Short: "Unpack GeoSite resources",
	RunE:  unpack,
	Args:  cobra.ExactArgs(1),
}

func unpack(cmd *cobra.Command, args []string) error {
	fmt.Println("➕Loading file:", args[0])
	fileContent, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	fmt.Println("📦Unpacking GeoSite database")
	startTime := time.Now()

	var geositeList []*v2raygeo.GeoSite
	geositeList, err = v2raygeo.LoadSite(fileContent)
	if err != nil {
		return err
	}

	if outDir == "" {
		outDir = "output"
	}
	err = os.Mkdir(outDir, 0o666)
	if err != nil && !os.IsExist(err) {
		return err
	}

	var count uint64
	for _, site := range geositeList {
		if code == "" || strings.EqualFold(code, site.CountryCode) {
			domains := processGeositeEntry(site)
			err = os.WriteFile(path.Join(outDir, strings.ToLower(site.CountryCode)),
				[]byte(strings.Join(domains, "\n")), 0o666)
			if err != nil {
				fmt.Println("❌Error when saving", site.CountryCode, "to text file, skipped:", err)
			}
			count += 1
			if code != "" {
				break // fast path
			}
		}
	}

	fmt.Println("🎉Successfully unpacked", count, "codes in", time.Now().Sub(startTime))
	return nil
}

func processGeositeEntry(vGeositeEntry *v2raygeo.GeoSite) []string {
	var domains []string
	var entry strings.Builder

	for _, domain := range vGeositeEntry.Domain {
		entry.Reset()
		entry.WriteString(strings.ToLower(domain.Type.String()))
		entry.WriteString(":")
		entry.WriteString(domain.Value)

		for _, attribute := range domain.Attribute {
			entry.WriteString(" @" + attribute.Key)
		}

		domains = append(domains, entry.String())
	}

	domains = common.Uniq(domains)
	sort.Strings(domains)
	return domains
}
