package app

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cg219/common-game/internal/database"
	"github.com/pressly/goose/v3"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

type Config struct {
    Data struct {
        Path string `yaml:"data"`
    } `yaml:"app"`
    Frontend embed.FS
    Migrations embed.FS
}

type AppCfg struct {
    config Config
    database *database.Queries
}

func NewConfig(frontend embed.FS, migrations embed.FS) *Config {
    cfg := &Config{}

    cfg.Data.Path = os.Getenv("APP_DATA")
    cfg.Frontend = frontend
    cfg.Migrations = migrations

    return cfg
}

func NewConfigFromSecrets(data []byte, frontend embed.FS, migrations embed.FS) *Config {
    cfg := &Config{}

    if err := yaml.Unmarshal(data, cfg); err != nil {
        log.Fatal("Error unmarshalling secrets file")
    }

    cfg.Frontend = frontend
    cfg.Migrations = migrations

    return cfg
}

func Run(config Config) error {
    cfg := &AppCfg{
        config: config,
    }

    cwd, _ := os.Getwd();
    db, err := sql.Open("sqlite", filepath.Join(cwd, config.Data.Path))
    if err != nil {
        return err
    }

    defer db.Close()

    goose.SetBaseFS(config.Migrations)
    goose.SetDialect("sqlite3")

    if err := goose.Up(db, "sql/migrations"); err != nil {
        return err
    }

    cfg.database = database.New(db)
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // App Logic
    go func() {
        StartServer(cfg)
    }()

    for {
        select {
        case <- ctx.Done():
            log.Println("terminating Run()")
            return nil
        }
    }
}

