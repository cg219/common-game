package auth

import (
	"crypto/rand"
	"encoding/base64"
)

const ChallengeSize = 32
const UserIdSize = 32

type Challenge = [ChallengeSize]byte
type UserId = [UserIdSize]byte
type RegistrationResponse struct {
    Challenge string `json:"challenge"`
    UserId string `json:"userid"`
    UserName string `json:"username"`
    UserDisplay string `json:"displayName"`
}

func CreateChallenge() (Challenge, error){
    var c Challenge

    if _, err := rand.Read(c[:]); err != nil {
        return c, err
    }

    return c, nil
}

func CreateUserId() (UserId, error) {
    var u UserId

    if _, err := rand.Read(u[:]); err != nil {
        return u, err
    }

    return u, nil
}

func CreateRegistration(name string, display string) RegistrationResponse {
    challenge, err := CreateChallenge()
    if err != nil {
        panic(err)
    }

    userId, err := CreateUserId()
    if err != nil {
        panic(err)
    }

    challengeStr := base64.RawStdEncoding.EncodeToString(challenge[:])
    userStr := base64.RawStdEncoding.EncodeToString(userId[:])
    reg := RegistrationResponse{
        Challenge: challengeStr,
        UserId: userStr,
        UserName: name,
        UserDisplay: display,
    }

    return reg
}
