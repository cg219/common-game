package main

import (
    "context"
    "database/sql"
    "encoding/json"
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
    "github.com/spiretechnology/go-webauthn"
    "github.com/tursodatabase/go-libsql"
)

type Server struct {
    mux *http.ServeMux
    games map[int]*LiveGameData
    ctx context.Context
    db *data.Queries
    wa webauthn.WebAuthn
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

type authPost struct {
    Username string `json:"username"`
}

type LiveGameData struct {
    game *game.Game
    mch chan<- game.Move
    sch <-chan game.StatusGroup
}

type ForwardRequestError struct {
    Error error
    DerivedError error
    ResponseWriter http.ResponseWriter
    Request *http.Request
    NextHandler http.Handler
}

type RegistrationPost struct {
    Username string `json:"username"`
    Response webauthn.RegistrationResponse `json:"response"`
}

type AuthPost struct {
    Username string `json:"username"`
    Response webauthn.AuthenticationResponse `json:"response"`
}

type CookieValues struct {
    GID string
    UID string
}

type MHandlerFunc func(w http.ResponseWriter, r *http.Request) error
type ContextKey int
type DerivedError int

const (
    GameId ContextKey = iota
    Error
)

const (
    Unauthorized DerivedError = iota
)

func (e DerivedError) String() string {
    return []string{"Unauthorized"}[e]
}

func WordExists(words []string, w string) bool {
    return slices.Contains(words, w)
}

func NewServer(db *data.Queries, wa webauthn.WebAuthn) *Server {
    return &Server{
        mux: http.NewServeMux(),
        games: make(map[int]*LiveGameData),
        ctx: context.Background(),
        db: db,
        wa: wa,
    }
}

func handle(h MHandlerFunc) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mwe := r.Context().Value(Error)

        if mwe != nil {
            mwe := mwe.(error)

            switch(mwe.Error()) {
            case Unauthorized.String():
                getDefault(w)                
                return
            }
        }

        if err := h(w, r); err != nil {
            log.Printf("ERRR: %s", err)
            w.Write(getErrResponse(err))
        }
    })
}

func (h MHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := h(w, r); err != nil {
        log.Println("ERRR")
    }
}

func startServer() error {
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

    q := data.New(db)
    srv := NewServer(q, webauthn.New(webauthn.Options{
        RP: webauthn.RelyingParty{
            ID: "localhost",
            Name: "The Common Game",
        },
        Credentials: auth.NewCredentials(q),
    }))

    srv.mux.Handle("GET /", authOnly(handle(srv.getGamePage)))
    srv.mux.Handle("GET /favicon.ico", http.FileServer(http.Dir("./web")))
    srv.mux.Handle("GET /sf/", http.StripPrefix("/sf", http.FileServer(http.Dir("./web"))))
    srv.mux.Handle("POST /api/game", authOnly(handle(srv.createGame)))
    srv.mux.Handle("PUT /api/game", playerOnly(handle(srv.updateGame)))
    srv.mux.Handle("POST /auth", handle(srv.createAuth))
    srv.mux.Handle("POST /auth/auth-verify", handle(srv.verifyAuth))
    srv.mux.Handle("POST /auth/register", handle(srv.createRegistration))
    srv.mux.Handle("POST /auth/verify", handle(srv.verifyRegistration))

    return http.ListenAndServe(":3000", srv.mux)
}

func forwardError(f *ForwardRequestError) {
    log.Printf("Forwared Error: %s\n", f.Error)
    ctx := context.WithValue(f.Request.Context(), Error, f.DerivedError)
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

func getDefault(w http.ResponseWriter) {
    tmpl := template.Must(template.ParseFiles("templates/pages/auth.html"))
    w.Header().Add("Content-Type", "text/html")
    tmpl.Execute(w, nil)
}

func (s *Server) getAuthPage(w http.ResponseWriter, r *http.Request) error {
    tmpl := template.Must(template.ParseFiles("templates/pages/auth.html"))
    w.Header().Add("Content-Type", "text/html")

    tmpl.Execute(w, nil)
    return nil
}

func (s *Server) getGamePage(w http.ResponseWriter, r *http.Request) error {
    tmpl := template.Must(template.ParseFiles("templates/pages/game.html"))
    w.Header().Add("Content-Type", "text/html")

    tmpl.Execute(w, nil)
    return nil
}

func authOnly(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.URL)
        fre := &ForwardRequestError{
            NextHandler: h,
            ResponseWriter: w,
            Request: r,
        }

        uidCookie, err := r.Cookie("the-connect-game-uid")
        if err != nil {
            fre.Error = err
            fre.DerivedError = fmt.Errorf(Unauthorized.String())
            forwardError(fre)
            return
        }

        stoken := uidCookie.Value
        uidToken, err := jwt.ParseWithClaims(stoken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("notsecure"), nil })

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        uid, err := uidToken.Claims.GetSubject()

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        if uid == "" {
            fre.Error = fmt.Errorf("Unauthorized")
            forwardError(fre)
            return
        }

        values := CookieValues{
            UID: uid,
        }

        ctx := context.WithValue(r.Context(), "cookieValues", values)
        r = r.WithContext(ctx)
        h.ServeHTTP(w, r)
    })
}

func playerOnly(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gidCookie, err := r.Cookie("the-connect-game-gid")
        if err != nil {
            log.Print(err)
        }

        stoken := gidCookie.Value
        gidToken, err := jwt.ParseWithClaims(stoken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("notsecure"), nil })

        uidCookie, err := r.Cookie("the-connect-game-uid")
        if err != nil {
            log.Print(err)
        }

        stoken = uidCookie.Value
        uidToken, err := jwt.ParseWithClaims(stoken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte("notsecure"), nil })

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

        gid, err := gidToken.Claims.GetSubject()

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        uid, err := uidToken.Claims.GetSubject()

        if err != nil {
            fre.Error = err
            forwardError(fre)
            return
        }

        values := CookieValues{
            GID: gid,
            UID: uid,
        }

        ctx := context.WithValue(r.Context(), "cookieValues", values)
        r = r.WithContext(ctx)
        h.ServeHTTP(w, r)
    })
}

func (s *Server) createAuth(w http.ResponseWriter, r *http.Request) error {
    w.Header().Add("Content-Type", "application/json")
    var body authPost

    defer r.Body.Close()

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return err
    }

    log.Printf("Body: %s\n", body)

    reg, err := s.wa.CreateAuthentication(s.ctx, webauthn.User{
        ID: body.Username,
        Name: body.Username,
        DisplayName: fmt.Sprintf("@%s", body.Username),
    })

    if err != nil {
        return err
    }

    log.Println(reg)

    data, _ := json.Marshal(reg)

    w.Header().Add("Content-Length", strconv.Itoa(len(data)))
    w.Write(data)

    return nil
}

func (s *Server) verifyAuth(w http.ResponseWriter, r *http.Request) error {
    w.Header().Add("Content-Type", "application/json")
    var body AuthPost 

    defer r.Body.Close()

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return err
    }

    log.Printf("Body: %s\n", body)

    u := webauthn.User{
        ID: strings.ToLower(body.Username),
        Name: strings.ToLower(body.Username),
        DisplayName: fmt.Sprintf("@%s", body.Username),
    }

    reg, err := s.wa.VerifyAuthentication(s.ctx, u, &body.Response)

    if err != nil {
        log.Printf("Auth Verification Error: %s, %s\n", u.ID, err)
        return fmt.Errorf("Verification Failed")
    }

    if reg == nil {
        log.Printf("Auth Verification Failed: %s\n", u.ID)
        return fmt.Errorf("Verification Failed")
    } 

    token := auth.NewToken("user", u.ID, time.Now().Add(60 * time.Minute).UTC())
    if err = token.Create(); err != nil {
        return err
    }

    cookie := auth.NewCookie("uid", token.Value(), int(time.Hour.Seconds()))
    http.SetCookie(w, &cookie)
    w.WriteHeader(http.StatusOK)
    return nil
}


func (s *Server) createRegistration(w http.ResponseWriter, r *http.Request) error {
    w.Header().Add("Content-Type", "application/json")
    var body authPost

    defer r.Body.Close()

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return err
    }

    reg, err := s.wa.CreateRegistration(s.ctx, webauthn.User{
        ID: body.Username,
        Name: body.Username,
        DisplayName: fmt.Sprintf("@%s", body.Username),
    })

    if err != nil {
        return err
    }

    log.Println(reg)

    data, _ := json.Marshal(reg)

    w.Header().Add("Content-Length", strconv.Itoa(len(data)))
    w.Write(data)

    return nil
}

func (s *Server) verifyRegistration(w http.ResponseWriter, r *http.Request) error {
    w.Header().Add("Content-Type", "application/json")
    var body RegistrationPost

    defer r.Body.Close()

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        return err
    }

    u := webauthn.User{
        ID: strings.ToLower(body.Username),
        Name: strings.ToLower(body.Username),
        DisplayName: fmt.Sprintf("@%s", body.Username),
    }

    reg, err := s.wa.VerifyRegistration(s.ctx, u, &body.Response)

    if err != nil {
        log.Println("Registration Verification Error: %s, %s", u.ID, err)
        return fmt.Errorf("Verification Failed")
    }

    if reg == nil {
        log.Println("Registration Verification Failed: %s", u.ID)
        return fmt.Errorf("Verification Failed")
    } 

    token := auth.NewToken("user", u.ID, time.Now().Add(60 * time.Minute).UTC())
    if err = token.Create(); err != nil {
        return err
    }

    cookie := auth.NewCookie("uid", token.Value(), int(time.Hour.Seconds()))
    http.SetCookie(w, &cookie)
    w.WriteHeader(http.StatusOK)
    return nil
}


func (s *Server) createGame(w http.ResponseWriter, r *http.Request) error {
    gameWords, err := s.db.GetSubjectsForGame(s.ctx)
    if err != nil {
        log.Printf("Error retreiving subjects: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    var values CookieValues
    rawCookieValue := r.Context().Value("cookieValues")

    if rawCookieValue != nil {
        values = rawCookieValue.(CookieValues) 
        log.Println(values)
    }

    log.Println(values.UID)

    uid := values.UID

    if uid == "" {
        return fmt.Errorf("Unauthorized")
    }

    game := game.Create(gameWords)
    id, err := s.db.SaveNewGame(s.ctx, data.SaveNewGameParams{
        Active: sql.NullBool{ Bool: true, Valid: true },
        Start: sql.NullInt64{ Int64: time.Now().UTC().UnixMilli(), Valid: true },
    })

    if err != nil {
        log.Printf("Error creating game: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    err = s.db.SaveUserToGame(s.ctx, data.SaveUserToGameParams{ Uid: uid, Gid: id})
    if err != nil {
        log.Printf("Error saving user to game: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    statusCh, moveCh := game.Run()
    s.games[int(id)] = &LiveGameData{
        game: game,
        mch: moveCh,
        sch: statusCh,
    }

    tmpl := template.Must(template.ParseFiles("templates/fragments/game-board.html"))
    token := auth.NewToken("game", fmt.Sprintf("%d", id), time.Now().Add(6 * time.Hour).UTC())
    if err = token.Create(); err != nil {
        log.Printf("Error creating token: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    gr := &GameResponse{
        GameId: int(id),
        Words: game.WordsWithData(),
        TurnsLeft: game.MaxTurns - game.Metadata.WrongTurns,
        Status: int(game.CheckStatus()),
    }

    cookie := auth.NewCookie("gid", token.Value(), int(6 * time.Hour.Seconds()))
    http.SetCookie(w, &cookie)
    w.Header().Add("Content-Type", "text/html")
    tmpl.Execute(w, gr)

    return nil
}

func (s *Server) updateGame(w http.ResponseWriter, r *http.Request) error {
    cerr := r.Context().Value(Error)

    if cerr, ok:= cerr.(error); ok {
        w.Header().Add("Content-Type", "application/json")
        w.Write(getErrResponse(error(cerr)))
        return nil
    }

    var values CookieValues
    rawCookieValue := r.Context().Value("cookieValues")

    if rawCookieValue != nil {
        values = rawCookieValue.(CookieValues) 
        log.Println(values)
    }

    gid, err := strconv.Atoi(values.GID)
    if err != nil {
        log.Printf("Error Getting Game ID Value: %s", err)
        return fmt.Errorf("Unathorized")
    }

    uid := values.UID
    if uid == "" {
        log.Printf("User ID empty")
        return fmt.Errorf("Unathorized")
    }

    guid, err := s.db.GetGameUidById(s.ctx, int64(gid))
    if err != nil {
        log.Printf("Error Getting Game: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    if guid != uid {
        log.Printf("Game UID doesn't match User ID guid: %s, uid: %s", guid, uid)
        return fmt.Errorf("Unauthorized")
    }

    d, ok := s.games[gid]

    if !ok {
        w.Header().Add("Content-Type", "application/json")
        w.Write(getErrResponse(fmt.Errorf("game with id %d not found", gid)))
        return nil
    }

    err = r.ParseForm()
    if err != nil {
        log.Printf("Error parsing form: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    var words [4]string

    for k, v := range r.Form {
        if strings.EqualFold(k,"words") {
            copy(words[:], v[:4])
        }
    }

    d.mch <- game.Move{ Words: words }
    status := <- d.sch

    switch status.Status.Status() {
    case game.Playing:
        s.db.UpdateGameTurns(s.ctx, data.UpdateGameTurnsParams{
            ID: int64(gid),
            Wrong: sql.NullInt64{ Int64: int64(d.game.Metadata.WrongTurns), Valid: true },
            Turns: sql.NullInt64{ Int64: int64(d.game.Metadata.TotalTurns), Valid: true },
        }) 
    case game.Win:
        s.db.UpdateGameStatus(s.ctx, data.UpdateGameStatusParams{
            ID: int64(gid),
            End: sql.NullInt64{ Int64: int64(time.Now().UTC().UnixMilli()), Valid: true },
            Active: sql.NullBool{ Bool: false, Valid: true },
            Win: sql.NullBool{ Bool: true, Valid: true },
        })
    case game.Lose:
        s.db.UpdateGameStatus(s.ctx, data.UpdateGameStatusParams{
            ID: int64(gid),
            End: sql.NullInt64{ Int64: int64(time.Now().UTC().UnixMilli()), Valid: true },
            Active: sql.NullBool{ Bool: false, Valid: true },
            Win: sql.NullBool{ Bool: false, Valid: true },
        })
    default:
        s.db.UpdateGame(s.ctx, data.UpdateGameParams{
            ID: int64(gid),
            Turns: sql.NullInt64{ Int64: int64(d.game.Metadata.TotalTurns), Valid: true },
            Wrong: sql.NullInt64{ Int64: int64(d.game.Metadata.WrongTurns), Valid: true },
        })
    }

    tmpl := template.Must(template.ParseFiles("templates/fragments/game-board.html"))
    gr := &GameResponse{
        GameId: gid,
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
