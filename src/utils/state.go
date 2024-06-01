package utils

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ClientState struct {
	IsAuthed bool
	Username string
	UserId   string
}

var DefaultClientState = ClientState{
	IsAuthed: false,
        Username: "anon",
        UserId: "",
}

type ContextClaimType struct{}

var ClaimsContextKey struct{} = ContextClaimType{}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId   string `json:"sub"`
	Username string `json:"username"`
	JwtType  string `json:"jwt_type"`
	Family   string `json:"family"`
}

func (c *CustomClaims) String() string {
	return fmt.Sprintf("%#v", c)
}

func GetClientState(r *http.Request) ClientState {
	claims, ok := r.Context().Value(ClaimsContextKey).(*CustomClaims)

	if claims == nil {
		state := DefaultClientState
		// leaving room in the future to add things to state even if claims is nil
		return state
	}

	return ClientState{
		IsAuthed: ok,
		Username: claims.Username,
		UserId:   claims.ID,
	}
}
