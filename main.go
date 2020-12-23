package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"

	// we do a blank import here because this packages init() function simply registers the package as the SQL driver for mysql(mariaDB), init() runs before main()
	_ "github.com/go-sql-driver/mysql"
)

func Dbconn() (db *sql.DB) {
	cm := Configuration_manager{}
	cm.v = viper.New()
	if cm.Load("./app.conf") == false {
		log.Fatalf("Can not load config file \n")
		os.Exit(1)
	}
	app := cm.GetAppConfig()
	connect := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", app.User_db, app.Password_db, app.IP_db, app.Dbname)
	db, err := sql.Open("mysql", connect)
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Println(err)
	if err != nil {
		log.Fatalf("Failed to open SQL database: %s", err)
	}
	return db

}

func main() {

	// Setup the mux router and register our routes
	r := mux.NewRouter()
	r.HandleFunc("/home", HomeHandler)
	// r.HandleFunc("/second-route", CoolNewHandler)

	// if we want to do something to each request before entering the route handler, for example checking an auth token
	// we can implement a middleware

	r.Use(loggingMiddleware, authMiddleware)

	// Now lets setup the server with a graceful shutddown
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
