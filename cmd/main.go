package main

import (
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/bwmarrin/snowflake"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	val := time.Now().Unix()
	val = (val % 999) + 1
	s, err := snowflake.NewNode(val)
	if err != nil {
		panic(err)
	}
	opts := handler.NewServerOptions{
		Repository: repo,
		Snowflake:  s,
		JWTCreate:  handler.JWTIssue,
		JWTVal:     handler.JWTValidate,
	}
	return handler.NewServer(opts)
}
