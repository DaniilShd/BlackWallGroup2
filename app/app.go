package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DaniilShd/BlackWallGroup/domain"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ClientMapRequest struct {
	domain domain.ClientMapCache
}

type ClientMapRequestAndHandler struct {
	MapRequest ClientMapRequest
	Handler    ClientHandler
}

func Start() {

	dbClient := getDbClient()

	cdb := ClientHandler{domain.NewClientRepositoryDb(dbClient)}
	// created map for cache request
	mapRequest := ClientMapRequest{domain.NewMapRequest(dbClient)}
	ch := ClientMapRequestAndHandler{
		MapRequest: mapRequest,
		Handler:    cdb,
	}

	router := mux.NewRouter()

	router.HandleFunc("/clients/{id_client:[0-9]+}", ch.makeTransaction).Methods(http.MethodPost)
	router.HandleFunc("/history/{id_client:[0-9]+}", cdb.getHistoryTransaction).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}

// Connecting to databases
func getDbClient() *sqlx.DB {
	HOST := "localhost"
	USER := "testuser"
	PASSWORD := "1234"
	DATABASE := "TestTransaction"

	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	client, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(); err != nil {
		fmt.Println("bad connection DB")
		return nil
	}

	fmt.Println("connection to DB")
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
