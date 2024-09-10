package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cg219/common-game/game"
	"github.com/cg219/common-game/internal/data"
	"github.com/golang-jwt/jwt/v5"
    _ "github.com/tursodatabase/go-libsql"
)

type server struct {
    mux *http.ServeMux
}

type GameResponse struct {
    Words []game.WordData`json:"words"`
    GameId int `json:"id"`
    TurnsLeft int `json:"moveLeft"`
    Status int `json:"status"`
    HasMove bool
    WordExists func([]string, string) bool
    Move struct {
        Correct bool `json:"correct"`
        Words  []string `json:"words,omitempty"`
    } `json:"move,omitempty"`
}

type errorResponse struct {
    Error string `json:"error"`
}

type moveResponse struct {
    Status int `json:"status"`
    Correct bool `json:"correct"`
    Words  []string `json:"words,omitempty"`
    Subject string `json:"subject,omitempty"`
    GameId int `json:"id"`
    TurnsLeft int `json:"moveLeft"`
}

type gamePost struct {
    Words [4]string `json:"words"`
}

type storeData struct {
    game *game.Game
    mch chan<- game.Move
    sch <-chan game.StatusGroup
}

type ForwardRequestError struct {
    Error error
    ResponseWriter http.ResponseWriter
    Request *http.Request
    NextHandler http.Handler
}

type MHandlerFunc func(w http.ResponseWriter, r *http.Request) error
type ContextKey int

var store map[int]*storeData
var globalQuery *data.Queries
var globalContext context.Context

const (
    GameId ContextKey = iota
    Error
)

func WordExists(words []string, w string) bool {
    return slices.Contains(words, w)
}

func newServer() *server {
    return &server{
        mux: http.NewServeMux(),
    }
}

func handle(h MHandlerFunc) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := h(w, r); err != nil {
            w.Write(getErrResponse(err))
            log.Println("ERRR")
        }
    })
}

func (h MHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := h(w, r); err != nil {
        log.Println("ERRR")
    }
}

func startServer() error {
    srv := newServer()
    store = make(map[int]*storeData)

    globalContext = context.Background()
    ddl, err := os.ReadFile("./configs/schema.sql")
    if err != nil {
        return err
    }

    db, err := sql.Open("libsql", "file:./database.db")
    if err != nil {
        return err
    }

    defer db.Close()

    if _, err := db.ExecContext(globalContext, string(ddl)); err != nil {
        return err
    }

    globalQuery = data.New(db)

    if err != nil {
        return err
    }

    srv.mux.HandleFunc("GET /", getHome())
    srv.mux.Handle("GET /sf/", http.StripPrefix("/sf", http.FileServer(http.Dir("./web"))))
    srv.mux.Handle("POST /api/game", handle(createGame))
    srv.mux.Handle("PUT /api/game", playerOnly(handle(updateGame)))

    return http.ListenAndServe(":3000", srv.mux)
}

func forwardError(f *ForwardRequestError) {
    log.Println(f.Error)
    ctx := context.WithValue(f.Request.Context(), Error, errors.New("JWT Error"))
    r := f.Request.WithContext(ctx)
    f.NextHandler.ServeHTTP(f.ResponseWriter, r)
}

func getErrResponse(e error) []byte {
    message := "Server Error: Something went wrong :c"
    res := &errorResponse{ Error: message }
    data, err := json.Marshal(res)
    log.Println(e.Error())

    if err != nil {
        panic(nil)
    }

    return data
}

func getHome() http.HandlerFunc {
    return func(w http.ResponseWriter, _ *http.Request) {
        tmpl := template.Must(template.ParseFiles("templates/pages/game.html"))
        w.Header().Add("Content-Type", "text/html")

        tmpl.Execute(w, nil)
    }
}

func playerOnly(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("the-connect-game")

        if err != nil {
            log.Print(err)
        }

        stoken := cookie.Value
        token, err := jwt.ParseWithClaims(stoken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("notsecure"), nil })

        fre := &ForwardRequestError{
            NextHandler: h,
            ResponseWriter: w,
            Request: r,
        }

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        sid, err := token.Claims.GetSubject()

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        id, err := strconv.Atoi(sid)

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        ctx := context.WithValue(r.Context(), GameId, id)
        r = r.WithContext(ctx)
        h.ServeHTTP(w, r)
    })
}

func createGame(w http.ResponseWriter, _ *http.Request) error {
    gc := &game.GameConfig{
        Q: globalQuery,
        Ctx: globalContext,
    }

    game, err := game.Create(*gc)

    if err != nil {
        return err
    }

    id := len(store)

    statusCh, moveCh := game.Run()
    store[id] = &storeData{
        game: game,
        mch: moveCh,
        sch: statusCh,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
        Issuer: "common-game",
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute).UTC()),
        Subject: fmt.Sprintf("%d", id),
    })

    stoken, err := token.SignedString([]byte("notsecure"))

    if err != nil {
        return err
    }

    tmpl := template.Must(template.ParseFiles("templates/fragments/game-board.html"))

    gr := &GameResponse{
        GameId: id,
        Words: game.WordsWithData(),
        TurnsLeft: game.MaxTurns - game.Metadata.WrongTurns,
        Status: int(game.CheckStatus()),
    }

    cookie := http.Cookie{
        Name: "the-connect-game",
        Value: stoken,
        Path: "/",
        MaxAge: 3600,
        HttpOnly: true,
        Secure: true,
        SameSite: http.SameSiteLaxMode,
    }

    http.SetCookie(w, &cookie)
    w.Header().Add("Content-Type", "text/html")
    tmpl.Execute(w, gr)

    return nil
}

func updateGame(w http.ResponseWriter, r *http.Request) error {
    cerr := r.Context().Value(Error)

    if cerr, ok:= cerr.(error); ok {
        w.Header().Add("Content-Type", "application/json")
        w.Write(getErrResponse(error(cerr)))
        return nil
    }

    id := r.Context().Value(GameId).(int)

    data, ok := store[id]

    if !ok {
        w.Header().Add("Content-Type", "application/json")
        w.Write(getErrResponse(fmt.Errorf("game with id %d not found", id)))
        return nil
    }

    err := r.ParseForm()

    if err != nil {
        return err
    }

    var words [4]string

    for k, v := range r.Form {
        if strings.EqualFold(k,"words") {
            copy(words[:], v[:4])
        }
    }

    data.mch <- game.Move{
        Words: words,
    }

    status := <- data.sch

    tmpl := template.Must(template.ParseFiles("templates/fragments/game-board.html"))

    gr := &GameResponse{
        GameId: id,
        Words: data.game.WordsWithData(),
        TurnsLeft: data.game.MaxTurns - data.game.Metadata.WrongTurns,
        Status: status.Status.Status().Enum(),
        HasMove: true,
        WordExists: WordExists,
        Move: struct{Correct bool "json:\"correct\""; Words []string "json:\"words,omitempty\""}{
            Correct: status.Status.Metadata.Correct,
            Words: status.Status.Metadata.Move.Words[:],
        },
    }

    w.Header().Add("Content-Type", "text/html")
    tmpl.Execute(w, gr)

    return nil
}
