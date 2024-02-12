package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/metacubex/geo/geoip"
	"github.com/metacubex/geo/geosite"

	F "github.com/sagernet/sing/common/format"
	"github.com/spf13/cobra"
)

func init() {
	commandLook.PersistentFlags().StringVarP(&dbType, "type", "t", "", "specify database type")
	commandLook.PersistentFlags().StringVarP(&dbPath, "file", "f", "", "specify database file path")
	commandLook.PersistentFlags().BoolVarP(&immediate, "immediate", "i", false, "return immediately as soon as a result is found")
	commandLook.PersistentFlags().BoolVarP(&noResolve, "no-resolve", "", false, "set no resolve for domains")
	mainCommand.AddCommand(commandLook)
}

var commandLook = &cobra.Command{
	Use:   "look",
	Short: "Query geo information from databases",
	RunE:  look,
	Args:  cobra.ExactArgs(1),
}

var (
	immediate bool
	noResolve bool
)

func look(cmd *cobra.Command, args []string) error {
	var (
		ipPaths   []string
		sitePaths []string
		err       error
	)
	if dbPath == "" {
		ipPaths, err = findIP()
		if err != nil {
			fmt.Println("⚠", err)
		}
		sitePaths, err = findSite()
		if err != nil {
			fmt.Println("⚠", err)
		}
	} else {
		ipPaths = []string{dbPath}
		sitePaths = []string{dbPath}
	}

	ip := net.ParseIP(args[0])
	if ip == nil { // domain
		fmt.Println("🔎Querying from", sitePaths)
		result := make(map[string]struct{})
		startTime := time.Now()
		domainName := args[0]
		for _, filePath := range sitePaths {
			var db *geosite.Database
			db, err = geosite.FromFile(filePath)
			if err != nil {
				fmt.Println("❌Error when loading", filePath, "as a GeoSite database, skipped.")
				continue
			}

			if immediate {
				code := db.LookupCode(domainName)
				if code != "" {
					result[strings.ToLower(code)] = struct{}{}
					break
				}
			} else {
				codes := db.LookupCodes(domainName)
				for _, code := range codes {
					result[strings.ToLower(code)] = struct{}{}
				}
			}
		}

		os.Stdout.WriteString(F.ToString("🎉Query finished in ", time.Now().Sub(startTime), "!\n"))
		fmt.Print("Total ", len(result), " results (GeoSite codes):\n  ")
		for code := range result {
			os.Stdout.WriteString(code)
			os.Stdout.WriteString(" ")
		}

		if noResolve {
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		go func() {
			<-ctx.Done()
			fmt.Print("\n\n🌎Resolved ", domainName, " timeout")
		}()

		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", domainName)
		if err != nil {
			fmt.Println("\n\n❌Fail to resolve", domainName, ", skipped.")
		}
		if len(ips) != 0 {
			ip = ips[0]
		}

		fmt.Print("\n\n🌎Resolved ", domainName, " as ", ip, "\n\n")
	}

	fmt.Println("🔎Querying from", ipPaths)
	result := make(map[string]struct{})
	startTime := time.Now()
	for _, filePath := range ipPaths {
		var db *geoip.Database
		db, err = geoip.FromFile(filePath)
		if err != nil {
			fmt.Println("❌Error when loading", filePath, "as a GeoIP database, skipped.")
			continue
		}

		codes := db.LookupCode(ip)
		for _, code := range codes {
			result[strings.ToUpper(code)] = struct{}{}
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
	os.Stdout.WriteString("\n")
	return nil
}
