package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgaunet/template-website/internal/webserver"
)

//go:generate templ generate

// for website, if you need to handle big file or stream, you must stay on stdlib or chi (just router)
// otherwise, you can use fiber or echo
// for api, you can use fiber or echo

func main() {
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
	err := w.Run()
	if err != nil {
		if err == http.ErrServerClosed {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "hgf %v\n", err)
		os.Exit(1)
	}
}
