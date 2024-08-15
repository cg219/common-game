package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cg219/common-game/game"
)

type server struct {
    mux *http.ServeMux
}

type gameResponse struct {
    Words []string `json:"words"`
    GameId int `json:"id"`
}

type moveResponse struct {
    Status string `json:"status"`
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

var store map[int]*storeData

func newServer() *server {
    return &server{
        mux: http.NewServeMux(),
    }
}

func startServer() error {
    srv := newServer()
    store = make(map[int]*storeData)

    srv.mux.HandleFunc("GET /", getHome())
    srv.mux.HandleFunc("GET /game", getGame())
    srv.mux.HandleFunc("POST /game/{id}", postGame())

    return http.ListenAndServe(":3000", srv.mux)
}

func getHome() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/plain")
        w.Write([]byte("Yay we're here!!"))
    }
}

func getGame() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        game, err := game.Create()

        if err != nil {
            panic(err)
        }

        id := len(store)

        statusCh, moveCh := game.Run()
        store[id] = &storeData{
            game: game,
            mch: moveCh,
            sch: statusCh,
        }

        gr := &gameResponse{
            GameId: id,
            Words: game.Words(),
        }

        res, err := json.Marshal(gr)

        if err != nil {
            panic(err)
        }

        w.Write(res)
    }
}

func postGame() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        pv := r.PathValue("id")
        id, err := strconv.Atoi(pv)

        if err != nil {
            w.Write([]byte("{\"status\": false}"))
            return
        }

        data, ok := store[id]

        if !ok {
            w.Write([]byte("{\"status\": false}"))
            return
        }

        body := &gamePost{}
        err = json.NewDecoder(r.Body).Decode(body)

        if err != nil {
            w.Write([]byte("{\"status\": false}"))
            return
        } 

        data.mch <- game.Move{
            Words: body.Words,
        }

        status := <- data.sch

        moveRes := &moveResponse{
            Status: status.Status.String(),
            GameId:  id,
            TurnsLeft: data.game.MaxTurns - data.game.Metadata.WrongTurns,
        }

        res, err := json.Marshal(moveRes)

        if err != nil {
            w.Write([]byte("{\"status\": false}"))
            return
        }

        log.Println(id)
        w.Write(res)
    }
}
