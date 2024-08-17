package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cg219/common-game/game"
	"github.com/golang-jwt/jwt/v5"
)

type server struct {
    mux *http.ServeMux
}

type gameResponse struct {
    Words []string `json:"words"`
    GameId int `json:"id"`
    Token string `json:"token"`
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

type ContextKey int

var store map[int]*storeData

const (
    GameId ContextKey = iota
    Error
)

func newServer() *server {
    return &server{
        mux: http.NewServeMux(),
    }
}

func startServer() error {
    srv := newServer()
    store = make(map[int]*storeData)

    srv.mux.HandleFunc("GET /", getHome())
    srv.mux.HandleFunc("POST /game", createGame())
    srv.mux.Handle("PUT /game", mwGetAuth(updateGame()))

    return http.ListenAndServe(":3000", srv.mux)
}

func getHome() http.HandlerFunc {
    return func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Add("Content-Type", "text/plain")
        w.Write([]byte("Yay we're here!!"))
    }
}

func createGame() http.HandlerFunc {
    return func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        game, err := game.Create()

        if err != nil {
            w.Write(getErrResponse(err))
            return
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
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute).UTC()),
            Subject: fmt.Sprintf("%d", id),
        })

        stoken, err := token.SignedString([]byte("notsecure"))

        if err != nil {
            w.Write(getErrResponse(err))
            return
        }

        gr := &gameResponse{
            GameId: id,
            Words: game.Words(),
            Token: stoken,
        }

        res, err := json.Marshal(gr)

        if err != nil {
            w.Write(getErrResponse(err))
            return
        }

        w.Write(res)
    }
}

func mwGetAuth(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        stoken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
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

func updateGame() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")

        cerr := r.Context().Value(Error)

        if cerr, ok:= cerr.(error); ok {
            w.Write(getErrResponse(error(cerr)))
            return
        }

        id := r.Context().Value(GameId).(int)

        data, ok := store[id]

        if !ok {
            w.Write(getErrResponse(fmt.Errorf("game with id %d not found", id)))
            return
        }

        body := &gamePost{}
        err := json.NewDecoder(r.Body).Decode(body)

        if err != nil {
            w.Write(getErrResponse(err))
            return
        }

        data.mch <- game.Move{
            Words: body.Words,
        }

        status := <- data.sch

        moveRes := &moveResponse{
            Correct: status.Status.Metadata.Correct,
            GameId:  id,
            TurnsLeft: data.game.MaxTurns - data.game.Metadata.WrongTurns,
            Status: status.Status.Status().Enum(),
        }


        if moveRes.Correct {
            moveRes.Subject = status.Status.Metadata.Subject.Name
            moveRes.Words = status.Status.Metadata.Move.Words[:]
        }

        res, err := json.Marshal(moveRes)

        if err != nil {
            w.Write(getErrResponse(err))
            return
        }

        log.Println(id)
        w.Write(res)
    }
}
