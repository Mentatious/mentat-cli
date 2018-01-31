package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Mentatious/mentat-cli/commands"
	"github.com/Mentatious/mentat-cli/io/format"
	"github.com/wiedzmin/goodies"
	"github.com/ybbus/jsonrpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	goodies.HandleKeyboardInterrupt()

	app := kingpin.New("mentat-cli", "Mentat command line client")
	APIHost := app.Flag("apihost", "API host address").Short('a').Required().String()
	User := app.Flag("user", "Free-form user ID to associate imported entries with").Short('u').Required().String()
	Quiet := app.Flag("quiet", "Be quiet").Short('q').Bool()
	importCommand := app.Command("import", "Import data into Mentat DB")
	importDeliciousCommand := importCommand.Command("delicious", "Import Delicious bookmarks dump data")
	importPocketCommand := importCommand.Command("pocket", "Import Pocket bookmarks dump data")
	importFile := importCommand.Flag("file", "File to import from").Short('f').String()
	searchCommand := app.Command("search", "Search the saved data corpus")
	searchQuery := app.Flag("query", "Raw JSON query to API").Short('r').String()
	searchResultsFormat := app.Flag("format", "Results format, e.g. raw, json, etc").Short('x').String()
	searchOutput := app.Flag("output", "File to save results to").Short('o').Default("stdout").String()
	searchTypes := app.Flag("types", "Comma-separated list of entry types to search for").Short('t').String()
	searchContent := app.Flag("content", "Comma-separated list of entry types to search for").Short('c').String()
	searchTags := app.Flag("tags", "Comma-separated list of tags to search for").Short('T').String()
	searchPrio := app.Flag("prio", "Comma-separated list of priority values to search for").Short('p').String()
	deleteCommand := app.Command("delete", "delete entries from Mentat DB")
	deleteUUIDs := deleteCommand.Flag("uuids", "Comma-separated list of UUIDs of entries to delete").Short('i').String()
	cleanupCommand := app.Command("cleanup", "cleanup Mentat DB (type-wise)")
	cleanupTypes := cleanupCommand.Flag("type", "types of entries to cleanup").Short('s').String()

	parsedParams := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *APIHost == "" {
		fmt.Printf("No API host provided, exiting...")
		os.Exit(1)
	} else {
		apiHostCheck, err := http.Get(fmt.Sprintf("http://%s", *APIHost))
		if err != nil {
			fmt.Printf("Cannot dial %s, exiting... (%s)", *APIHost, err.Error())
			os.Exit(1)
		}
		defer apiHostCheck.Body.Close()
	}
	apiserverURL := fmt.Sprintf("http://%s/mentat/v1/", *APIHost)
	rpcClient := jsonrpc.NewRPCClient(apiserverURL)
	switch parsedParams {
	case importDeliciousCommand.FullCommand():
		commands.ImportDelicious(*importFile, rpcClient, *User, *Quiet)
	case importPocketCommand.FullCommand():
		commands.ImportPocket(*importFile, rpcClient, *User, *Quiet)
	case searchCommand.FullCommand():
		results := commands.Search(rpcClient, *User, *searchTypes, *searchContent, *searchTags, *searchPrio, *searchQuery)
		switch *searchResultsFormat {
		case "json":
			format.DumpJSON(results, *searchOutput)
		case "csv":
			format.DumpCSV(results, *searchOutput)
		default:
			fmt.Printf("No export format provided, dumping raw contents:\n%v", results)
		}
	case deleteCommand.FullCommand():
		var deleted int
		commands.Delete(rpcClient, *User, *deleteUUIDs, &deleted)
		fmt.Printf("deleted %d entries: OK", deleted)
	case cleanupCommand.FullCommand():
		var deleted int
		commands.Cleanup(rpcClient, *User, *cleanupTypes, &deleted)
		fmt.Printf("cleaned up %d entries: OK", deleted)
	default:
		fmt.Printf("Unknown dump format, exiting...")
		os.Exit(1)
	}
}
