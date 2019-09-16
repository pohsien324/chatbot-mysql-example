package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paked/messenger"
)

var (
	verifyToken = os.Getenv("VERIFY_TOKEN")
	pageToken   = os.Getenv("PAGE_TOKEN")
	port        = os.Getenv("PORT")
)

func createSchema(db *sql.DB) {
	db.Exec("create database chatbot;")
	db.Exec("use chatbot;")
	db.Exec("create table facebook( id int not null primary key auto_increment,keyword text, response text)CHARSET=utf8 COLLATE=utf8_general_ci; ")
	db.Exec("insert into facebook(keyword,response) values ('Hi','Hi there');")
	fmt.Println("Create database chatbot.")
}

func main() {
	var (
		mysqlServer     = os.Args[1]
		mysqlServerPort = os.Args[2]
		mysqlUsername   = os.Args[3]
		mysqlPassword   = os.Args[4]
	)

	if verifyToken == "" || pageToken == "" {
		log.Fatal("missing arguments")
	}

	if port == "" {
		port = "80"
	}

	if mysqlServer == "" || mysqlServerPort == "" || mysqlUsername == "" || mysqlPassword == "" {
		log.Fatal("missing database arguments")
	}

	//  Connect and initial the mysql server
	connectionInformation := fmt.Sprintf("%s:%s@tcp(%s:%s)/", mysqlUsername, mysqlPassword, mysqlServer, mysqlServerPort)
	db, err := sql.Open("mysql", connectionInformation)
	defer db.Close()

	err = db.Ping()
	for err != nil {
		fmt.Println("Can not connect to the mysql server.")
		fmt.Println("Retry after 5 seconds.")
		time.Sleep(5 * time.Second)
		err = db.Ping()
		fmt.Println(err)
	}
	createSchema(db)
	fmt.Println("Initial database successfully.")

	// Setup Messenger Bot Webhook Server
	client := messenger.New(messenger.Options{
		VerifyToken: verifyToken,
		Token:       pageToken,
	})

	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		//Get the response from Database
		var response string
		row := db.QueryRow("SELECT response FROM chatbot.facebook  WHERE keyword = ?", m.Text)
		erro := row.Scan(&response)
		if erro != nil {
			if erro == sql.ErrNoRows {
				// No response match the keyword.
				response = "Sorry, I don't know what you say?"
			} else if strings.Contains(erro.Error(), "Unknown database") {
				// Database chatbot is not exist.
				response = "There is an internal problem, please contact to the administrator."
				fmt.Println(erro)
				createSchema(db)
			} else {
				// Can not connect to mysql server
				response = "There is an internal problem, please contact to the administrator."
				fmt.Println(erro)
			}
		}
		r.Text(response, messenger.ResponseType)
	})

	if err := http.ListenAndServe(":"+port, client.Handler()); err != nil {
		log.Fatal(err)
	}
}
