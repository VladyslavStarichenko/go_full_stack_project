package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var database *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "goldsilver12"
	dbname   = "practice3_go"
)

func main() {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/", actorsHandler)
	router.HandleFunc("/create_actor", createActorHandler)
	router.HandleFunc("/edit_actor/{id:[0-9]+}", EditActorPage).Methods("GET")
	router.HandleFunc("/edit_actor/{id:[0-9]+}", EditActorHandler).Methods("POST")
	router.HandleFunc("/remove_actor/{id:[0-9]+}", DeleteActorHandler)
	router.HandleFunc("/performance_statistics", performanceStatisticsHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
