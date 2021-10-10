package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/astronomy/ap-manager/handlers"
	"gitlab.com/astronomy/ap-manager/helpers"
	"gitlab.com/astronomy/ap-manager/storage"
)

func main() {
	db, err := sql.Open("mysql", "root:rootPassword$@tcp(127.0.0.1:3306)/ap-manager")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	argLength := len(os.Args[1:])
	if argLength == 1 && os.Args[1] == "logFile" {
		a := helpers.ProcessLogFile()
		storage.CreateSessions(db, a)
	}

	if argLength == 0 {
		env := &handlers.Env{DB: db}
		http.Handle("/health-check", handlers.Handler{Env: env, H: handlers.HealthCheckRouteHandler, Method: http.MethodGet})
		http.Handle("/sessions", handlers.Handler{Env: env, H: handlers.SessionsRouteHandler, Method: http.MethodGet})
		log.Println("Serving on host - port 8085.....")
		errHTTP := http.ListenAndServe(":8085", nil)
		if errHTTP != nil {
			log.Fatal("Web server for http is toast", errHTTP)
		}
	}
}
