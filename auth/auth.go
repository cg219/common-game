package auth

import (
	"cmp"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
    subject string
    expiresAt time.Time
    Name string
    value string
}

func NewToken(n, s string, e time.Time) Token {
    return Token{
        Name: n,
        subject: s,
        expiresAt: e,
    }
}

func NewCookie(n, v string, ma int) http.Cookie {
    return http.Cookie{
        Name: fmt.Sprintf("the-connect-game-%s", n),
        Value: v,
        Path: "/",
        MaxAge: ma,
        HttpOnly: true,
        Secure: true,
        SameSite: http.SameSiteLaxMode,
    }
}

func (t *Token) Value() string {
    return t.value
}

func (t *Token) Create() error {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
        Issuer: "common-game",
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(t.expiresAt),
        Subject: fmt.Sprintf("%s", t.subject),
    })

    stoken, err := token.SignedString([]byte("notsecure"))

    if err != nil {
        return fmt.Errorf("Error creating %s token. val: %t, err %s", cmp.Or(t.Name, "new"), t.subject, err)
    }

    t.value = stoken

    return  nil
}
