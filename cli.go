package main

import (
	"fmt"
	"github.com/Mentatious/mentat-cli/importers"
	"github.com/ybbus/jsonrpc"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// InitLogging ... Initialize loggers
func InitLogging(debug bool, showLoc bool) (*zap.Logger, *zap.SugaredLogger) {
	var rawlog *zap.Logger
	var log *zap.SugaredLogger
	var cfg zap.Config
	var err error
	if debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.DisableCaller = !showLoc
	rawlog, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	log = rawlog.Sugar()
	return rawlog, log
}

func main() {
	_, log := InitLogging(false, false)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Caught ^C, exiting...")
		os.Exit(1)
	}()

	app := kingpin.New("mentat-cli", "Mentat command line client")
	APIHost := app.Flag("apihost", "API host address").Short('a').Required().String()
	User := app.Flag("user", "Free-form user ID to associate imported entries with").Short('u').Required().String()
	Quiet := app.Flag("quiet", "Be quiet").Short('q').Bool()
	importCommand := app.Command("import", "Import data into Mentat DB")
	importDeliciousCommand := importCommand.Command("delicious", "Import Delicious bookmarks dump data")
	importPocketCommand := importCommand.Command("pocket", "Import Pocket bookmarks dump data")
	importFile := importCommand.Flag("file", "File to import from").Short('f').String()

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
		importers.ImportDelicious(*importFile, rpcClient, *User, log, *Quiet)
	case importPocketCommand.FullCommand():
		importers.ImportPocket(*importFile, rpcClient, *User, log, *Quiet)
	default:
		fmt.Printf("Unknown dump format, exiting...")
		os.Exit(1)
	}
}
