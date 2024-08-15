package main

import (
	"encoding/json"
	"net/http"

	"github.com/cg219/common-game/game"
)

type server struct {
    mux *http.ServeMux
}

type gameResponse struct {
    Words []string `json:"words"`
    GameId int `json:"id"`
}

var store map[int]*game.Game

func newServer() *server {
    return &server{
        mux: http.NewServeMux(),
    }
}

func startServer() error {
    srv := newServer()
    store = make(map[int]*game.Game)

    srv.mux.HandleFunc("GET /", serveHome())
    srv.mux.HandleFunc("GET /game", serveGame())

    return http.ListenAndServe(":3000", srv.mux)
}

func serveHome() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/plain")
        w.Write([]byte("Yay we're here!!"))
    }
}

func serveGame() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        game, err := game.Create()

        if err != nil {
            panic(err)
        }

        id := len(store)

        store[id] = game

        gr := &gameResponse {
            GameId: id,
            Words: make([]string, 3),
        }

        res, err := json.Marshal(gr)

        if err != nil {
            panic(err)
        }

        w.Header().Add("Content-Type", "application/json")
        w.Write(res)
    }
}
