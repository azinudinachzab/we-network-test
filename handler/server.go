package handler

import (
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/constant"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
)

type Server struct {
	Repository repository.RepositoryInterface
	Snowflake  *snowflake.Node
	JWTCreate  func(ttl time.Duration, content interface{}) (string, error)
	JWTVal     func(token string) (interface{}, error)
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Snowflake   *snowflake.Node
	JWTCreate  func(ttl time.Duration, content interface{}) (string, error)
	JWTVal     func(token string) (interface{}, error)
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Snowflake: opts.Snowflake,
		JWTCreate: opts.JWTCreate,
		JWTVal:    opts.JWTVal,
	}
}

func JWTIssue(ttl time.Duration, content interface{}) (string, error) {
	privateKey := constant.PrivateKey
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func JWTValidate(token string) (interface{}, error) {
	publicKey := constant.PublicKey
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return "", err
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, err
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}
