package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Email struct {
    From string
    To string
    Subject string
    Body string
}

func (e *Email) Send(apikey string, endpoint string) {
    var buf bytes.Buffer

    writer := multipart.NewWriter(&buf)
    writer.WriteField("from", e.From)
    writer.WriteField("to", e.To)
    writer.WriteField("subject", e.Subject)
    writer.WriteField("text", e.Body)
    writer.Close()

    reader := bytes.NewReader(buf.Bytes())

    r, _ := http.NewRequestWithContext(context.Background(), "POST", endpoint, reader)
    r.SetBasicAuth("api", apikey)
    r.Header.Set("Content-Type", writer.FormDataContentType())

    c := http.Client{
        Timeout: 30 * time.Second,
    }

    res, err := c.Do(r)
    if err != nil {
        log.Println(err)
    }

    defer res.Body.Close()

    var val map[string]interface{}

    err = json.NewDecoder(res.Body).Decode(&val)
    if err != nil && err != io.EOF {
        log.Println(err)
    }

    log.Println("email response")
    log.Println(val)
}

func StartEmailServer(cfg *AppCfg) {
    open := false

    go func() {
        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
        defer stop()

        <- ctx.Done()

        log.Println("Shutting Down Email Server")
    }()

    for {
        select {
        case m, ok := <-cfg.emails:
            if !ok {
                return
            }

            log.Printf("new email: %v\n", m)

            if !open {
                open = true
                m.Send(cfg.config.Email.Key, cfg.config.Email.Host)
            }

            log.Println("email sent")
        case <-time.After(30 * time.Second):
            if open {
                log.Println("closing conn")
                open = false
                log.Println("closed conn")
            }
        }
    }
}
