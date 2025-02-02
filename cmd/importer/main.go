package main

import (
	"bytes"
	"log"
	"os"

	"github.com/cg219/common-game/internal/importer"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Specify a CSV to import")
    }

    csvbytes, err := os.ReadFile(os.Args[1])
    if err != nil {
        log.Fatalf("err loading file: %s", err.Error())
    }

    importer.Run(bytes.NewReader(csvbytes))
}

