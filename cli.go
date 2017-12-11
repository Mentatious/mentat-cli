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
	importCommand := app.Command("import", "Import data into Mentat DB")
	importFormat := importCommand.Flag("type", "Type of data to import").Short('t').String()
	importFile := importCommand.Flag("file", "File to import from").Short('f').String()
	importAPIHost := importCommand.Flag("apihost", "API host address").Short('a').String()
	importQuiet := importCommand.Flag("quiet", "Be quiet").Short('q').Bool()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case importCommand.FullCommand():
		if *importAPIHost == "" {
			fmt.Printf("No API host provided, exiting...")
			os.Exit(1)
		} else {
			apiHostCheck, err := http.Get(fmt.Sprintf("http://%s", *importAPIHost))
			if err != nil {
				fmt.Printf("Cannot dial %s, exiting... (%s)", *importAPIHost, err.Error())
				os.Exit(1)
			}
			defer apiHostCheck.Body.Close()
		}
		apiserverURL := fmt.Sprintf("http://%s/mentat/v1/", *importAPIHost)
		rpcClient := jsonrpc.NewRPCClient(apiserverURL)
		if *importFormat == "delicious" {
			importers.ImportDelicious(*importFile, rpcClient, log, *importQuiet)
		} else if *importFormat == "pocket" {
			importers.ImportPocket(*importFile, rpcClient, log, *importQuiet)
		} else {
			fmt.Printf("Unknown dump format: '%s', exiting...", *importFormat)
			os.Exit(1)
		}
	}
}
