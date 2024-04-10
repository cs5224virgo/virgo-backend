package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	perrors "github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// ErrInvalidToken when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
)

// Claims are the structure of the jwt
type Claims struct {
	UserID uint `json:"virgo_id"`
	// other necessary jwt claim info
	jwtgo.RegisteredClaims
}

// NewToken creates a new jwt token for the given parameters
func NewToken(userID uint, username string, expiration time.Time) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			Subject:   username,
			ExpiresAt: jwtgo.NewNumericDate(expiration),
			Issuer:    "virgo",
		},
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt_signing_key")))
}

// ParseToken parses and check a jwt token if its valid
func ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwtgo.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt_signing_key")), nil
	})
	if err != nil {
		return nil, perrors.Wrap(err, "error while parsing jwt token")
	}

	// validate
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// probably wont be, but still, check if userid is zero
	if token.Claims.(*Claims).UserID == 0 {
		return nil, fmt.Errorf("%w: id in jwt token is zero wtf", ErrInvalidToken)
	}

	return token.Claims.(*Claims), err
}
