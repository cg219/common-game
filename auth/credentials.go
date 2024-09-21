package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cg219/common-game/internal/data"
	"github.com/spiretechnology/go-webauthn"
)

type Credentials struct {
    Db *data.Queries
}

type CredentialsJSON struct {
    ID string `json:"id"`
    Type string `json:"type"`
    PublicKey string `json:"publickey"`
    PublickKeyAlg int `json:"publickeyalg"`
}

func NewCredentials(db *data.Queries) Credentials {
    return Credentials{ Db: db }
}

func (c Credentials) GetCredentials(ctx context.Context, user webauthn.User) ([]webauthn.Credential, error) {
    log.Printf("User: %s\n",  user)
    u, err := c.Db.GetUserById(ctx, sql.NullString{ String: user.ID, Valid: true })
    if err != nil {
        return nil, err
    }

    if u.KeyString != "" {
        keys := strings.Split(u.KeyString, "|:|")
        creds := make([]webauthn.Credential, len(keys))

        for _, v := range keys {
            var cj CredentialsJSON
            if err = json.Unmarshal([]byte(v), &cj); err != nil {
                return nil, fmt.Errorf("Error decoding credentials: %s", err)
            }

            cid, err := base64.StdEncoding.DecodeString(cj.ID)
            if err != nil {
                return nil, fmt.Errorf("Error decoding credentials: %s", err)
            }

            pub, err := base64.StdEncoding.DecodeString(cj.PublicKey)
            if err != nil {
                return nil, fmt.Errorf("Error decoding credentials: %s", err)
            }

            c := &webauthn.Credential{
                ID: cid,
                Type: cj.Type,
                PublicKey: pub,
                PublicKeyAlg: cj.PublickKeyAlg,
            }

            creds = append(creds, *c)
        }

        return creds, nil
    }
    return nil, nil
}

func (c Credentials) GetCredential(ctx context.Context, user webauthn.User, credentialId []byte) (*webauthn.Credential, error) {
    log.Printf("Getting USer By Key: %s", base64.RawURLEncoding.EncodeToString(credentialId))
    u, err := c.Db.GetUserByKey(ctx, sql.NullString{ String: base64.RawURLEncoding.EncodeToString(credentialId), Valid: true })
    if err != nil {
        return nil, err
    } 

    if u.Keys.Valid {
        var cj CredentialsJSON
        if err = json.Unmarshal([]byte(u.Keys.String), &cj); err != nil {
            return nil, fmt.Errorf("Error decoding credentials: %s", err)
        }

        cid, err := base64.StdEncoding.DecodeString(cj.ID)
        if err != nil {
            return nil, fmt.Errorf("Error decoding credentials: %s", err)
        }

        pub, err := base64.StdEncoding.DecodeString(cj.PublicKey)
        if err != nil {
            return nil, fmt.Errorf("Error decoding credentials: %s", err)
        }

        c := &webauthn.Credential{
            ID: cid,
            Type: cj.Type,
            PublicKey: pub,
            PublicKeyAlg: cj.PublickKeyAlg,
        }

        return c, nil
    }

    return nil, nil
}

func (c Credentials) StoreCredential(ctx context.Context, user webauthn.User, credential webauthn.Credential, meta webauthn.CredentialMeta) error {
    cj := &CredentialsJSON{
        ID: base64.StdEncoding.EncodeToString(credential.ID),
        Type: credential.Type,
        PublicKey: base64.StdEncoding.EncodeToString(credential.PublicKey),
        PublickKeyAlg: credential.PublicKeyAlg,
    }

    credData, err := json.Marshal(cj)
    if err != nil {
        return err
    }

    row := data.SaveUserParams{
        Uid: sql.NullString{ String: user.ID, Valid: true },
        KeyID: sql.NullString{ String: base64.RawURLEncoding.EncodeToString(credential.ID), Valid: true },
        Keys: sql.NullString{ String: string(credData), Valid: true },
    }

    log.Printf("Save: %s", row)

    err = c.Db.SaveUser(ctx, row)
    if err != nil {
        return err
    }

    return nil
}
