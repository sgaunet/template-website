package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgaunet/template-website/internal/webserver"
	"github.com/sgaunet/template-website/pkg/config"
)

//go:generate templ generate

// for website, if you need to handle big file or stream, you must stay on stdlib or chi (just router)
// otherwise, you can use fiber or echo
// for api, you can use fiber or echo

var version string
var defaultPort = 8080
var defaultSizeSignalChannel = 5

func printVersion() {
	fmt.Printf("%s\n", version)
}

func main() {
	var (
		err         error
		cfgFile     string
		versionFlag bool
	)
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()

	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromFileOrEnvVar(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML file: %s\n", err)
		os.Exit(1)
	}
	if !cfg.IsValid() {
		fmt.Fprintf(os.Stderr, "Invalid configuration\n")
		os.Exit(1)
	}

	w := webserver.NewWebserver(nil, defaultPort)

	// handle graceful shutdown
	sigs := make(chan os.Signal, defaultSizeSignalChannel)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err := w.Shutdown()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}()
	err = w.Run()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "hgf %v\n", err)
		os.Exit(1)
	}
}
