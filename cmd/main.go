package main

import (
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

	w := webserver.NewWebserver(nil, 8080)

	// handle graceful shutdown
	sigs := make(chan os.Signal, 5)
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
		if err == http.ErrServerClosed {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "hgf %v\n", err)
		os.Exit(1)
	}
}
