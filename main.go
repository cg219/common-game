package main

import (
	"embed"
	"log"
	"os"

	"github.com/cg219/common-game/internal/app"
)

//go:embed static-app
var Frontend embed.FS

//go:embed sql/migrations/*.sql
var Migrations embed.FS

func main() {
    var cfg *app.Config
    done := make(chan struct{})
    secretsPath := os.Getenv("APP_CREDTENTIALS")
    _, err := os.Stat(secretsPath)

    if err != nil {
        if os.IsNotExist(err) {
            log.Printf("secrets file not found: %s\nFalling back to env variables\n", secretsPath)
            cfg = app.NewConfig(Frontend, Migrations)
        } else if os.IsPermission(err) {
            log.Printf("incorrect permissions on secret file: %s\nFalling back to env variables\n", secretsPath)
            cfg = app.NewConfig(Frontend, Migrations)
        } else {
            log.Fatal(err)
        }
    } else {
        data, err := os.ReadFile(secretsPath)
        if err != nil {
            log.Printf("error loading secrets file: %s; err: %s\nFalling back to env variables\n", secretsPath, err.Error())
        }

        cfg = app.NewConfigFromSecrets(data, Frontend, Migrations)
    }

    go func() {
        if err := app.Run(*cfg); err != nil {
            log.Fatal(err)
            close(done)
            return
        }
        log.Println("Exiting app func")

        close(done)
    }()

    <- done

    log.Println("Exiting main safely")
}

