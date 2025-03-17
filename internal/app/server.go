package app

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cg219/common-game/internal/database"
	"github.com/cg219/common-game/internal/game"
	"github.com/cg219/common-game/pkg/argon2id"
	"github.com/cg219/common-game/pkg/webtoken"
	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
    mux *http.ServeMux
    appcfg *AppCfg
    log *slog.Logger
    hasher *argon2id.Argon2id
    games map[int]*LiveGameData
    storage *Storage
}

type LiveGameData struct {
    game *game.Game
    mch chan<- game.Move
    sch <-chan game.StatusGroup
}

type GameResponse struct {
    Words []game.WordData`json:"words"`
    GameId int `json:"id"`
    TurnsLeft int `json:"moveLeft"`
    Status int `json:"status"`
    HasMove bool  `json:"hasMove"`
    Move GameResponseMove `json:"move,omitempty"`
}

type GameResponseSubject struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

type GameResponseMove struct {
    Correct bool `json:"correct"`
    Words  []string `json:"words,omitempty"`
    Subjects []GameResponseSubject `json:"subjects,omitempty"`
}

type SuccessResp struct {
    Success bool `json:"success"`
}

type TokenPacket struct{
    AccessToken string
    RefreshToken string
}

type ResponseError struct {
    Code int `json:"code"`
    Success bool `json:"success"`
    Messaage string `json:"message"`
}

type CandlerFunc func(w http.ResponseWriter, r *http.Request) error

const (
    INTERNAL_ERROR = "Internal Server Error"
    AUTH_ERROR = "Authentication Error"
    USERNAME_EXISTS_ERROR = "Username Exists Error"
    MISSING_PARAMS_ERROR = "Missing Parameters Error"
    GOTO_NEXT_HANDLER_ERROR = "Redirect Error"
    REDIRECT_ERROR = "Intentional Redirect Error"
)
const (
    CODE_USER_EXISTS = iota
    AUTH_FAIL
    AUTH_NOT_ALLOWED
    INTERNAL_SERVER_ERROR
    BAD_REQUEST
)

func NewServer(cfg *AppCfg) *Server {
    logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

    return &Server{
        mux: http.NewServeMux(),
        appcfg: cfg,
        log: logger, 
        hasher: argon2id.NewArgon2id(16 * 1024, 2, 1, 16, 32),
        games: make(map[int]*LiveGameData),
        storage: NewStorage(cfg.database, logger),
    }
}

func loadActiveGames(srv *Server) {
    _, games := srv.storage.GetActiveGames()

    for _, v := range games {
        srv.games[v.Id] = v.Data
    }
}

func addRoutes(srv *Server) {
    static, err := fs.Sub(srv.appcfg.config.Frontend, "static-app/assets")

    if err != nil {
        log.Fatal("error creating file subsystem")
    }

    srv.mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
    })

    srv.mux.Handle("GET /", srv.handle(srv.RedirectAuthenticated("/game", true), srv.getLoginPage))
    srv.mux.Handle("GET /game", srv.handle(srv.RedirectAuthenticated("/", false), srv.getGamePage))
    srv.mux.Handle("GET /report", srv.handle(srv.RedirectAuthenticated("/", false), srv.getReportPage))
    srv.mux.Handle("GET /account", srv.handle(srv.RedirectAuthenticated("/", false), srv.getAccountPage))
    srv.mux.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServer(http.FS(static))))
    srv.mux.Handle("POST /api/generate-apikey/{name}", srv.handle(srv.UserOnly, srv.GenerateAPIKey))
    srv.mux.Handle("POST /api/forgot-password", srv.handle(srv.ForgotPassword))
    srv.mux.Handle("POST /api/reset-password", srv.handle(srv.ResetPassword))
    srv.mux.Handle("POST /api/report", srv.handle(srv.UserOnly, srv.ReportBug))
    srv.mux.Handle("POST /api/game", srv.handle(srv.UserOnly, srv.CreateGame))
    srv.mux.Handle("PUT /api/game", srv.handle(srv.UserOnly, srv.UpdateGame))
    srv.mux.Handle("POST /auth/register", srv.handle(srv.Register))
    srv.mux.Handle("POST /auth/login", srv.handle(srv.Login))
    srv.mux.Handle("POST /auth/logout", srv.handle(srv.UserOnly, srv.Logout))
    srv.mux.Handle("GET /validate/{validvalue}", srv.handle(srv.ValidateRegistration)) 
    srv.mux.Handle("GET /reset/{resetvalue}", srv.handle(srv.getResetPage))
    srv.mux.Handle("POST /reset/{resetvalue}", srv.handle(srv.GetResetPasswordData))
    srv.mux.Handle("GET /healthcheck", srv.handle(srv.HealthCheck))
}

func (h CandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := h(w, r); err != nil {
        fmt.Println("OOPS")
    }
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) error {
    resp := SuccessResp{ Success: true }

    encode(w, http.StatusOK, resp)
    return nil
}

func (s *Server) GetGame(ctx context.Context, uid int64) (*game.Game, int, error) {
    gid, err := s.appcfg.database.GetActiveGidForUid(ctx, uid)
    if err != nil {
        if err != sql.ErrNoRows {
            s.log.Error("Error performing active game check", "err",err)
            return nil, 0, fmt.Errorf(INTERNAL_ERROR)
        }
    }

    if gid == 0 {
        err, id, g := s.storage.GetNewGameWithContext(ctx, s.appcfg.connection, int(uid))
        if err != nil {
            return nil, 0, fmt.Errorf(INTERNAL_ERROR)
        }

        statusCh, moveCh := g.Run()
        s.games[int(id)] = &LiveGameData{
            game: g,
            mch: moveCh,
            sch: statusCh,
        }

        return g, int(id), nil
    } else {
        d, ok := s.games[int(gid)]

        if !ok {
            s.log.Error("Getting game from server games", "gid", gid, "err", err)
            return nil, int(gid), fmt.Errorf(INTERNAL_ERROR)
        }

        return d.game, int(gid), nil
    }
}

func (s *Server) CreateGame(w http.ResponseWriter, r *http.Request) error {
    username := r.Context().Value("username").(string)
    user, err := s.appcfg.database.GetUser(r.Context(), username)
    if err != nil {
        s.log.Error("Error retreiving user", "err",err)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    game, id, err := s.GetGame(r.Context(), user.ID)
    subjects := make([]GameResponseSubject, 0)

    for _, v := range game.CompletedSubjects {
        subjects = append(subjects, GameResponseSubject{
            Id: v,
            Name: game.Subjects[v].Name,
        }) 
    }

    gr := &GameResponse{
        GameId: id,
        Words: game.WordsWithData(),
        TurnsLeft: game.MaxTurns - game.Metadata.WrongTurns,
        Status: int(game.CheckStatus()),
        HasMove: true,
        Move: GameResponseMove{
            Subjects: subjects,
        },
    }

    encode(w, http.StatusOK, gr)
    return nil
}

func (s *Server) UpdateGame(w http.ResponseWriter, r *http.Request) error {
    type Body struct {
        Words [4]string `json:"words"`
        Gid int `json:"gid"`
    }

    username := r.Context().Value("username").(string)
    user, err := s.appcfg.database.GetUser(r.Context(), username)
    if err != nil {
        s.log.Error("Error retreiving user", "err",err)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    body, err := decode[Body](r)
    if err != nil {
        s.log.Error("Error decoding body", "err",err)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    r.Body.Close()

    guid, err := s.appcfg.database.GetGameUidByGameId(r.Context(), int64(body.Gid))
    if err != nil {
        log.Printf("Error Getting Game: %s", err)
        return fmt.Errorf("Internal Server Error")
    }

    if guid != user.ID {
        s.log.Error("Game UID doesn't match User ID", "guid", guid, "uid", user.ID, "err", err)
        return fmt.Errorf(AUTH_ERROR)
    }

    d, ok := s.games[body.Gid]

    if !ok {
        s.log.Error("Getting game from server games", "gid", body.Gid, "err", err)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    s.log.Info("making move", "game", body.Gid)

    d.mch <- game.Move{ Words: body.Words }
    status := <- d.sch

    subjects := s.storage.UpdateGameWithContext(r.Context(), int64(body.Gid), status, d.game)

    gr := &GameResponse{
        GameId: body.Gid,
        Words: d.game.WordsWithData(),
        TurnsLeft: d.game.MaxTurns - d.game.Metadata.WrongTurns,
        Status: status.Status.Status().Enum(),
        HasMove: true,
        Move: GameResponseMove{
            Correct: status.Status.Metadata.Correct,
            Words: status.Status.Metadata.Move.Words[:],
            Subjects: subjects,
        },
    }

    encode(w, http.StatusOK, gr)
    return nil
}

func (s *Server) GenerateAPIKey(w http.ResponseWriter, r *http.Request) error {
    username := r.Context().Value("username")
    user, _ := s.appcfg.database.GetUser(r.Context(), username.(string))

    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    key := make([]byte, 24)

    for i := range key {
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            s.log.Error("generating api key", "err", err)
            return fmt.Errorf(INTERNAL_ERROR)
        }

        key[i] = charset[n.Int64()]
    }

    s.appcfg.database.SaveApiKey(r.Context(), database.SaveApiKeyParams{
        Key: string(key),
        Uid: sql.NullInt64{ Valid: true, Int64: user.ID },
        Name: r.PathValue("name"),
    })

    type KeyResp struct {
        Key string `json:"apikey"`
    }

    resp := KeyResp{ Key: string(key) }
    encode(w, http.StatusOK, resp)
    return nil
}

func (s *Server) ResetPassword(w http.ResponseWriter, r *http.Request) error {
    resettimer := time.Now().Unix()
    type Body struct {
        Username string `json:"username"`
        Reset string `json:"reset"`
        Password string `json:"password"`
        PasswordConfirm string `json:"passwordConfirm"`
    }

    body, err := decode[Body](r)
    if err != nil {
        return err
    }

    if !strings.EqualFold(body.Password, body.PasswordConfirm) {
        return fmt.Errorf(AUTH_ERROR)
    }

    hashPass, _ := s.hasher.EncodeFromString(body.Password)

    s.appcfg.database.ResetPassword(r.Context(), database.ResetPasswordParams{
        Reset: sql.NullString{ String: body.Reset, Valid: true },
        ResetTime: sql.NullInt64{ Int64: resettimer, Valid: true },
        Password: hashPass,
    })

    data := SuccessResp{ Success: true }

    encode(w, 200, data)
    return nil
}

func (s *Server) ReportBug(w http.ResponseWriter, r *http.Request) error {
    type Body struct {
        Problem string `json:"problem"`
        Result string `json:"result"`
        Steps string `json:"steps"`
    }

    user, _ := s.appcfg.database.GetUser(r.Context(), r.Context().Value("username").(string))

    defer r.Body.Close()

    data, err := decode[Body](r)
    err = s.appcfg.database.ReportBug(r.Context(), database.ReportBugParams{
        Problem: data.Problem,
        Result: data.Result,
        Steps: data.Steps,
        Uid: user.ID,
    })

    if err != nil {
        s.log.Error("reporting bug", "err", err)
    }

    e := Email{
        From: s.appcfg.config.Email.From,
        To: s.appcfg.config.App.Admin,
        Subject: "The Common Game - Bug Report",
        Body: fmt.Sprintf("Bug Report from: %s - %s\n\nProblem:\n%s\n\nExpected Results:\n%s\n\nReproduction Steps:\n%s\n\n", user.Username, user.Email, data.Problem, data.Result, data.Steps),
    }

    s.appcfg.emails <- e

    encode(w, http.StatusOK, SuccessResp{ Success: true })
    return nil
}

func (s *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
    resettimer := time.Now().Add(time.Minute * 15).Unix()
    resetbytes := make([]byte, 32)
    rand.Read(resetbytes)
    reset := base64.URLEncoding.EncodeToString(resetbytes)[:16]
    err := r.ParseForm()

    if err != nil {
        s.log.Error("parsing form", "err", err)
    }

    username := r.FormValue("username")

    err = s.appcfg.database.SetPasswordReset(r.Context(), database.SetPasswordResetParams{
        Reset: sql.NullString{ String: reset, Valid: true },
        ResetTime: sql.NullInt64{ Int64: resettimer, Valid: true },
        Username: username,
    })

    if err != nil {
        s.log.Error("resetting pass", "err", err)
    }

    user, err := s.appcfg.database.GetUser(r.Context(), username)
    if err != nil {
        s.log.Error("retreiving user for reset pass", "username", username, "err", err)
    }

    e := Email{
        From: s.appcfg.config.Email.From,
        To: user.Email,
        Subject: "The Common Game - Forgot Password",
        Body: fmt.Sprintf("Reset your password link:\n%s/reset/%s", s.appcfg.config.App.Url, reset),
    }

    s.appcfg.emails <- e

    s.log.Info("Reset Link:", "url", fmt.Sprintf("Reset your password link:\n%s/reset/%s", s.appcfg.config.App.Url, reset))
    return nil
}

func (s *Server) getFile(w http.ResponseWriter, filepath string) {
    data, err := s.appcfg.config.Frontend.ReadFile(filepath)
    if err != nil {
        s.log.Error("getting file", "path", filepath)
        return
    }

    w.Header().Add("Content-Type", "text/html")
    w.Write(data)
}

func (s *Server) getAccountPage(w http.ResponseWriter, r *http.Request) error {
    s.getFile(w, "static-app/entrypoints/settings.html")
    return nil
}

func (s *Server) getResetPage(w http.ResponseWriter, r *http.Request) error {
    s.getFile(w, "static-app/entrypoints/reset.html")
    return nil
}

func (s *Server) getLoginPage(w http.ResponseWriter, r *http.Request) error {
    s.getFile(w, "static-app/entrypoints/auth.html")
    return nil
}

func (s *Server) getGamePage(w http.ResponseWriter, r *http.Request) error {
    s.getFile(w, "static-app/entrypoints/game.html")
    return nil
}

func (s *Server) getReportPage(w http.ResponseWriter, r *http.Request) error {
    s.getFile(w, "static-app/entrypoints/report.html")
    return nil
}

func (s *Server) setTokens(w http.ResponseWriter, r *http.Request, username string) {
    accessToken := webtoken.NewToken("accessToken", username, "notsecure", time.Now().Add(time.Hour * 1))
    refreshToken := webtoken.NewToken("refreshToken", webtoken.GenerateRefreshString(), "notsecure", time.Now().Add(time.Hour * 24 * 30))
    accessToken.Create("thecommongame")
    refreshToken.Create("thecommongame")
    cookieValue := webtoken.CookieAuthValue{ AccessToken: accessToken.Value(), RefreshToken: refreshToken.Value() }
    cookie := webtoken.NewAuthCookie("thecommongame", "/", cookieValue, int(time.Hour * 24 * 30))

    s.appcfg.database.SaveUserSession(r.Context(), database.SaveUserSessionParams{
        Accesstoken: accessToken.Value(),
        Refreshtoken: refreshToken.Subject(),
    })

    http.SetCookie(w, &cookie)
}

func (s *Server) unsetTokens(w http.ResponseWriter, r *http.Request) {
    accesstoken := r.Context().Value("accesstoken").(string)
    refreshtoken := r.Context().Value("refreshtoken").(string)
    s.log.Info("unset tokens", "refresh", refreshtoken, "access", accesstoken)

    s.appcfg.database.InvalidateUserSession(r.Context(), database.InvalidateUserSessionParams{ Accesstoken: accesstoken, Refreshtoken: refreshtoken, })
    cookie := webtoken.NewAuthCookie("thecommongame", "/", webtoken.CookieAuthValue{}, int(0))

    http.SetCookie(w, &cookie)
    *r = *r.WithContext(context.Background())
}

func (s *Server) authenticateRequest(r *http.Request, username string) {
    ctx := context.WithValue(r.Context(), "username", username)
    updatedRequest := r.WithContext(ctx)

    *r = *updatedRequest
}

func (s *Server) getAuthGookie(r *http.Request) (string, string) {
    cookie, err := r.Cookie("thecommongame")
    if err != nil {
        s.log.Error("Cookie Retrieval", "cookie", "thecommongame", "method", "UserOnly", "request", r, "error", err.Error())
        return "", ""
    }

    value, err := base64.StdEncoding.DecodeString(cookie.Value)
    if err != nil {
        s.log.Error("Base64 Decoding", "cookie", cookie.Value, "method", "UserOnly", "request", r, "error", err.Error())
        return "", ""
    }

    var cookieValue webtoken.CookieAuthValue
    err = json.Unmarshal(value, &cookieValue)
    if err != nil {
        s.log.Error("Invalid Cookie Value", "cookie", cookie.Value, "method", "UserOnly", "request", r, "error", err.Error())
        return "", ""
    }

    return cookieValue.AccessToken, cookieValue.RefreshToken
}

func (s *Server) login(ctx context.Context, username string, password string) bool {
    existingUser, err := s.appcfg.database.GetUserWithPassword(ctx, username)
    if err != nil {
        if err == sql.ErrNoRows {
            return false
        }

        s.log.Error("sql err", "err", err)
        return false
    }

    s.log.Info("existing user", "user", existingUser)

    if existingUser.Username == "" {
        return false 
    }

    if !existingUser.Valid.Valid {
        return false
    }

    correct, _ := s.hasher.Compare(password, existingUser.Password)
    if !correct {
        s.log.Info("Password Mismatch", "password", password)
        return false
    }

    return true
}

func (s* Server) refreshAccessToken(ctx context.Context, refreshExpire int64, refreshTokenString, refreshValue, username string, w http.ResponseWriter) {
    accessToken := webtoken.NewToken("accessToken", username, "notsecure", time.Now().Add(time.Hour * 1))
    accessToken.Create("thecommongame")
    cookieValue := webtoken.CookieAuthValue{ AccessToken: accessToken.Value(), RefreshToken: refreshTokenString }
    cookie := webtoken.NewAuthCookie("thecommongame", "/", cookieValue, int(refreshExpire))

    s.appcfg.database.SaveUserSession(ctx, database.SaveUserSessionParams{
        Accesstoken: accessToken.Value(),
        Refreshtoken: refreshValue,
    })

    http.SetCookie(w, &cookie)
    s.log.Info("Refresh User Tokens", "username", username)
}

func (s *Server) isAuthenticated(ctx context.Context, ats, rts string) (bool, string, func(http.ResponseWriter) context.Context, context.Context) {
    accessTokenExpired := true
    refreshTokenExpired := true
    accessToken, err := webtoken.GetParsedJWT(ats, "notsecure")
    if err != nil {
        fmt.Println()

        if !strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
            s.log.Error("Invalid AccessToken", "accessToken", ats, "method", "IsAuthenticated", "error", err.Error())
            return false, "", nil, nil
        }
    } else {
        accessTokenExpired = false
    }

    refreshToken, err := webtoken.GetParsedJWT(rts, "notsecure")
    if err != nil {
        if !strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
            s.log.Error("Invalid RefreshToken", "refreshToken", rts, "method", "isAuthenticated", "error", err.Error())
            return false, "", nil, nil
        }
    } else {
        refreshTokenExpired = false
    }

    rfs, err := refreshToken.Claims.GetSubject()
    if err != nil {
        s.log.Error("Invalid RefreshToken", "method", "isAuthenticated", "error", err.Error())
        return false, "", nil, nil
    }

    var rf webtoken.Subject
    err = json.Unmarshal([]byte(rfs), &rf)
    if err != nil {
        s.log.Error("Invalid RefreshToken", "refreshToken", rfs, "method", "isAuthenticated", "error", err.Error())
        return false, "", nil, nil
    }

    if refreshTokenExpired {
        s.log.Error("Expired RefreshToken", "refreshToken", rts, "method", "isAuthenticated")
        s.appcfg.database.InvalidateUserSession(ctx, database.InvalidateUserSessionParams{
            Accesstoken: ats,
            Refreshtoken: rf.Value,
        })
        return false, "", nil, nil
    }

    _, err = s.appcfg.database.GetUserSession(ctx, database.GetUserSessionParams{
        Accesstoken: ats,
        Refreshtoken: rf.Value,
    })
    if err != nil {
        s.log.Error("Retreiving User Session", "method", "isAuthenticated", "error", err.Error())
        return false, "", nil, nil
    }

    us, err := accessToken.Claims.GetSubject()
    if err != nil {
        s.log.Error("Invalid AccessToken", "method", "isAuthenticated", "error", err.Error())
        return false, "", nil, nil
    }

    var username webtoken.Subject
    err = json.Unmarshal([]byte(us), &username)
    if err != nil {
        s.log.Error("Invalid AccessToken", "accessToken", us, "method", "isAuthenticated", "error", err.Error())
        return false, "", nil, nil
    }

    if accessTokenExpired {
        s.log.Error("Expired AccessToken", "accessToken", ats, "method", "isAuthenticated")
        s.appcfg.database.InvalidateUserSession(ctx, database.InvalidateUserSessionParams{
            Accesstoken: ats,
            Refreshtoken: rf.Value,
        })

        expiresAt, _ := refreshToken.Claims.GetExpirationTime()

        return false, username.Value, func(w http.ResponseWriter) context.Context {
            s.refreshAccessToken(ctx, expiresAt.Unix(), rts, rf.Value, username.Value, w)
            ctx = context.WithValue(ctx, "accesstoken", ats)
            ctx = context.WithValue(ctx, "refreshtoken", rf.Value)

            return ctx
        }, nil
    }

    ctx = context.WithValue(ctx, "accesstoken", ats)
    ctx = context.WithValue(ctx, "refreshtoken", rf.Value)

    return true, username.Value, nil, ctx 
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) error {
    type RegisterBody struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Email string `json:"email"`
    }

    body, err := decode[RegisterBody](r)
    if err != nil {
        return err
    }

    if body.Username == "" || body.Email == "" || body.Password == "" {
       return fmt.Errorf(MISSING_PARAMS_ERROR)
    }

    hashPass, err := s.hasher.EncodeFromString(body.Password)
    if err != nil {
        s.log.Error("Encoding Password", "password", body.Password)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    err, validToken := s.storage.NewUserWithContext(r.Context(), body.Email, body.Username, hashPass)
    if err != nil {
        return err
    }

    e := Email{
        From: s.appcfg.config.Email.From,
        To: body.Email,
        Subject: "The Common Game - Validate Email",
        Body: fmt.Sprintf("Validate your email link:\n%s/validate/%s", s.appcfg.config.App.Url, validToken),
    }

    s.appcfg.emails <- e

    if err != nil {
        s.log.Error("Saving New User", "username", body.Username, "err", err)
        return fmt.Errorf(INTERNAL_ERROR)
    }

    encode(w, http.StatusOK, SuccessResp{ Success: true })
    s.log.Info("Register Body", "body", body)
    return nil
}

func (s *Server) ValidateRegistration(w http.ResponseWriter, r *http.Request) error {
    err, valid, username := s.storage.ValidateNewUserWithContext(r.Context(), r.PathValue("validvalue"))
    if err != nil {
        return err;
    }

    if valid {
        s.setTokens(w, r, username)
        http.Redirect(w, r, "/game", http.StatusSeeOther)
        return nil
    }
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) error {
    type LoginBody struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    body, err := decode[LoginBody](r)
    if err != nil {
        s.log.Info("decoding loginbody", "body", body, "err", err)
        return err
    }

    if !s.login(r.Context(), body.Username, body.Password) {
        return fmt.Errorf(AUTH_ERROR)
    }

    s.setTokens(w, r, body.Username)
    encode(w, http.StatusOK, SuccessResp{ Success: true })
    s.log.Info("Login", "username", body.Username)
    return nil
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) error {
    s.unsetTokens(w, r)
    encode(w, http.StatusOK, SuccessResp{ Success: true })
    s.log.Info("Logout", "success", true)
    return nil
}

func (s *Server) GetResetPasswordData(w http.ResponseWriter, r *http.Request) error {
    type Data struct {
        Valid bool `json:"valid"`
        Username string `json:"username"`
        Reset string `json:"reset"`
    }

    reset := r.PathValue("resetvalue")

    dbValue, _ := s.appcfg.database.CanResetPassword(r.Context(), database.CanResetPasswordParams{
        ResetTime: sql.NullInt64{ Int64: time.Now().Unix(), Valid: true },
        Reset: sql.NullString{ String: reset, Valid: true },
    })

    data := Data{ Valid: dbValue.Valid, Username: dbValue.Username, Reset: reset }

    encode(w, 200, data)
    return nil
}

func StartServer(cfg *AppCfg) error {
    srv := NewServer(cfg)

    loadActiveGames(srv)
    addRoutes(srv)

    server := &http.Server{
        Addr: fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")),
        Handler: srv.mux,
    }

    go func(s *http.Server) {
        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
        defer stop()

        <- ctx.Done()

        log.Println("Shutting Down Server")

        if err := s.Shutdown(ctx); err != nil {
            log.Println("Shutdown error")
        }
    }(server)

    return server.ListenAndServe()
}

