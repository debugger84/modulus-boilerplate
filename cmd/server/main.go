package main

import (
	graphInit "boilerplate/internal/graph"
	"boilerplate/internal/pgx"
	"boilerplate/internal/post"
	"boilerplate/internal/user"
	"fmt"
	application "github.com/debugger84/modulus-application"
	db "github.com/debugger84/modulus-db-pg-gorm"
	graphql "github.com/debugger84/modulus-graphql"
	logger "github.com/debugger84/modulus-logger-zap"
	router "github.com/debugger84/modulus-router-httprouter"
	"net/http"
	"runtime"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, this is a boilerplate for the Modulus framefork!\n")
}

func main() {
	runtime.GOMAXPROCS(24)
	loggerConfig := logger.NewModuleConfig(nil)

	userConfig := user.NewModuleConfig()
	postConfig := post.NewModuleConfig()
	routes := router.NewRoutes()
	routes.Get("home", "/", hello)

	routerConfig := router.NewModuleConfig()
	routerConfig.Routes.AddFromRoutes(routes)

	dbConfig := db.NewModuleConfig()

	graphQlConfig := graphql.NewModuleConfig()
	graphQlInitConfig := graphInit.NewModuleConfig()

	pgxConfig := pgx.NewModuleConfig()

	app := application.New(
		[]interface{}{
			loggerConfig,
			routerConfig,
			dbConfig,
			userConfig,
			pgxConfig,
			graphQlConfig,
			graphQlInitConfig,
			postConfig,
		},
	)
	err := app.Run()
	if err != nil {
		panic("Cannot run application: " + err.Error())
	}
}
