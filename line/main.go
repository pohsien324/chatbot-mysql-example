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
	"github.com/line/line-bot-sdk-go/linebot"
)

func createSchema(db *sql.DB) {
	db.Exec("create database chatbot;")
	db.Exec("use chatbot;")
	db.Exec("create table line( id int not null primary key auto_increment,keyword text, response text)CHARSET=utf8 COLLATE=utf8_general_ci; ")
	db.Exec("insert into line(keyword,response) values ('Hi','Hi there');")
	fmt.Println("Create database chatbot.")
}

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	var (
		mysqlServer     = os.Args[1]
		mysqlServerPort = os.Args[2]
		mysqlUsername   = os.Args[3]
		mysqlPassword   = os.Args[4]
	)

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

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:

					//Get the response from Database
					var response string
					row := db.QueryRow("SELECT response FROM chatbot.line  WHERE keyword = ?", message.Text)
					erro := row.Scan(&response)
					if erro != nil {
						if erro == sql.ErrNoRows {
							// No response match the keyword.
							response = "Sorry, I don't know what you say?"
						} else if strings.Contains(erro.Error(), "Unknown database") {
							// database chatbot is not exist.
							response = "There is an internal problem, please contact to the administrator."
							fmt.Println(erro)
							createSchema(db)
						} else {
							// Can not connect to mysql server
							response = "There is an internal problem, please contact to the administrator."
							fmt.Println(erro)
						}
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
