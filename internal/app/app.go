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
    App struct {
        Url string `yaml:"url"`
        Admin string `yaml:"email"`
        Data string `yaml:"data"`
    } `yaml:"app"`
    Email struct {
        User string `yaml:"user"`
        Host string `yaml:"host"`
        From string `yaml:"from"`
        Key string `yaml:"key"`
    } `yaml:"email"`
    R2 struct {
        Key string `yaml:"key"`
        Secret string `yaml:"secret"`
        Token string `yaml:"token"`
        Url string `yaml:"url"`
    } `yaml:"r2"`
    Frontend embed.FS
    Migrations embed.FS
}

type AppCfg struct {
    config Config
    database *database.Queries
    connection *sql.DB
    emails chan Email
}

func NewConfig(frontend embed.FS, migrations embed.FS) *Config {
    cfg := &Config{}
    cfg.App.Data = os.Getenv("APP_DATA")
    cfg.App.Url = os.Getenv("APP_URL")
    cfg.App.Admin = os.Getenv("APP_ADMIN_EMAIL")
    cfg.Email.User = os.Getenv("SMTP_USER")
    cfg.Email.From = os.Getenv("SMTP_FROM")
    cfg.Email.Host = os.Getenv("SMTP_HOST")
    cfg.Email.Key = os.Getenv("SMTP_KEY")
    cfg.R2.Key = os.Getenv("R2_KEY")
    cfg.R2.Secret = os.Getenv("R2_SECRET")
    cfg.R2.Token = os.Getenv("R2_TOKEN")
    cfg.R2.Url = os.Getenv("R2_URL")
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
        emails: make(chan Email),
    }

    cwd, _ := os.Getwd();
    db, err := sql.Open("sqlite", filepath.Join(cwd, config.App.Data))
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
    cfg.connection = db
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    defer close(cfg.emails)

    // App Logic
    go func() {
        StartEmailServer(cfg)
    }()

    go func() {
        StartServer(cfg)
    }()

    <-ctx.Done()
    log.Println("terminating Run()")
    return nil
}

