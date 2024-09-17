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
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cg219/common-game/auth"
	"github.com/cg219/common-game/game"
	"github.com/cg219/common-game/internal/data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pressly/goose/v3"
	"github.com/tursodatabase/go-libsql"
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

type regsiterPost struct {
    Username string `json:"username"`
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
    dbName := os.Getenv("LOCAL_DB_NAME")
    dbUrl := os.Getenv("TURSO_DATABASE_URL")
    dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
    tmp, err := os.MkdirTemp("", "libdata-*")

    if  err != nil {
        return err
    }

    defer os.RemoveAll(tmp)
    dbPath := filepath.Join(tmp, dbName)
    conn, err := libsql.NewEmbeddedReplicaConnector(dbPath, dbUrl, libsql.WithAuthToken(dbAuthToken), libsql.WithSyncInterval(60))

    if err != nil {
        return err
    }

    defer conn.Close()
    globalContext = context.Background()
    db := sql.OpenDB(conn)
    defer db.Close()
    provider, err := goose.NewProvider(goose.DialectSQLite3, db, os.DirFS("./migrations"))

    if err != nil {
        return err
    }

    results, err := provider.Up(context.Background())

    if err != nil {
        return err
    }

    for _, r := range results {
        log.Println("goose: %s, %s", r.Source.Path, r.Duration)
    }

    globalQuery = data.New(db)

    if err != nil {
        return err
    }

    srv.mux.HandleFunc("GET /", getHome())
    srv.mux.Handle("GET /sf/", http.StripPrefix("/sf", http.FileServer(http.Dir("./web"))))
    srv.mux.Handle("POST /api/game", handle(createGame))
    srv.mux.Handle("PUT /api/game", playerOnly(handle(updateGame)))
    srv.mux.Handle("POST /auth/register", handle(createRegistration))
    srv.mux.Handle("POST /auth/verify", handle(verifyRegistration))

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
        tmpl := template.Must(template.ParseFiles("templates/pages/auth.html"))
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

func createRegistration(w http.ResponseWriter, r *http.Request) error {
    w.Header().Add("Content-Type", "application/json")
    var body regsiterPost

    defer r.Body.Close()

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return err
    }

    reg := auth.CreateRegistration(body.Username, fmt.Sprintf("@%s", body.Username))

    err =  json.NewEncoder(w).Encode(reg)
    if err != nil {
        return err
    }

    return nil
}

func verifyRegistration(w http.ResponseWriter, _ *http.Request) error {
    return nil
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

    id, err := gc.Q.SaveNewGame(gc.Ctx, data.SaveNewGameParams{
        Active: sql.NullBool{ Bool: true, Valid: true },
        PlayerID: sql.NullInt64{ Int64: 0, Valid: true },
        Start: sql.NullInt64{ Int64: time.Now().UTC().UnixMilli(), Valid: true },
    })

    if err != nil {
        return err
    }

    statusCh, moveCh := game.Run()
    store[int(id)] = &storeData{
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
        GameId: int(id),
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

    gc := &game.GameConfig{
        Q: globalQuery,
        Ctx: globalContext,
    }

    id := r.Context().Value(GameId).(int)

    d, ok := store[id]

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

    d.mch <- game.Move{
        Words: words,
    }

    status := <- d.sch

    switch status.Status.Status() {
    case game.Playing:
        gc.Q.UpdateGameTurns(gc.Ctx, data.UpdateGameTurnsParams{
            ID: int64(id),
            Wrong: sql.NullInt64{ Int64: int64(d.game.Metadata.WrongTurns), Valid: true },
            Turns: sql.NullInt64{ Int64: int64(d.game.Metadata.TotalTurns), Valid: true },
        }) 
    case game.Win:
        gc.Q.UpdateGameStatus(gc.Ctx, data.UpdateGameStatusParams{
            ID: int64(id),
            End: sql.NullInt64{ Int64: int64(time.Now().UTC().UnixMilli()), Valid: true },
            Active: sql.NullBool{ Bool: false, Valid: true },
            Win: sql.NullBool{ Bool: true, Valid: true },
        })
    case game.Lose:
        gc.Q.UpdateGameStatus(gc.Ctx, data.UpdateGameStatusParams{
            ID: int64(id),
            End: sql.NullInt64{ Int64: int64(time.Now().UTC().UnixMilli()), Valid: true },
            Active: sql.NullBool{ Bool: false, Valid: true },
            Win: sql.NullBool{ Bool: false, Valid: true },
        })
    default:
        gc.Q.UpdateGame(gc.Ctx, data.UpdateGameParams{
            ID: int64(id),
            Turns: sql.NullInt64{ Int64: int64(d.game.Metadata.TotalTurns), Valid: true },
            Wrong: sql.NullInt64{ Int64: int64(d.game.Metadata.WrongTurns), Valid: true },
        })
    }

    tmpl := template.Must(template.ParseFiles("templates/fragments/game-board.html"))

    gr := &GameResponse{
        GameId: id,
        Words: d.game.WordsWithData(),
        TurnsLeft: d.game.MaxTurns - d.game.Metadata.WrongTurns,
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
