package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/machinebox/sdk-go/boxutil"
)

func main() {
	fbep := os.Getenv("FACEBOX_URL")
	if fbep == "" {
		fmt.Println(`FACEBOX_URL must be specified.`)
		os.Exit(-1)
	}

	facebox := facebox.New(fbep)
	fmt.Println(`Face ID by Machine Box - https://machinebox.io/`)

	fmt.Println("Waiting for Facebox to be ready...")
	boxutil.WaitForReady(context.Background(), facebox)
	fmt.Println("Done!")

	fids := ":9000"
	srv := NewServer(facebox)
	if err := http.ListenAndServe(fids, srv); err != nil {
		log.Fatalln(err)
	}
}

