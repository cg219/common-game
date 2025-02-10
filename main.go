package main

import (
	"context"
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
    cwd, _ := os.Getwd();

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

    if cfg.R2.Key != "" {
        s3cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2.Key, cfg.R2.Secret, "")), config.WithRegion("auto"))

        if err != nil {
            log.Fatal(err)
        }

        client := s3.NewFromConfig(s3cfg, func(o *s3.Options) {
            o.BaseEndpoint = aws.String(cfg.R2.Url)
        })

        res, err := client.GetObject(context.Background(), &s3.GetObjectInput{
            Bucket: aws.String("commongame"),
            Key: aws.String("database.db"),
        })

        if err != nil {
            log.Fatal(err)
        }

        dbfile, err := os.Create(filepath.Join(cwd, cfg.App.Data))
        if err != nil {
            log.Fatal(err)
        }

        data, err := io.ReadAll(res.Body)
        if err != nil {
            res.Body.Close()
            log.Fatal(err)
        }

        res.Body.Close()

        _, err = dbfile.Write(data)
        if err != nil {
            dbfile.Close()
            log.Fatal(err)
        }

        dbfile.Close()
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

