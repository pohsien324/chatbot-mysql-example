package main

import (
	"database/sql"
	"fmt"
	"time"
	"strings"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/go-sql-driver/mysql"

)

func createSchema(db *sql.DB) {
	db.Exec("create database chatbot;")
	db.Exec("use chatbot;")
	db.Exec("create table telegram( id int not null primary key auto_increment,keyword text, response text)CHARSET=utf8 COLLATE=utf8_general_ci; ")
	db.Exec("insert into telegram(keyword,response) values ('Hi','Hi there');")
	fmt.Println("Create database chatbot.")
}


func main() {

	var (
		mysqlServer     = os.Args[1]
		mysqlServerPort = os.Args[2]
		mysqlUsername   = os.Args[3]
		mysqlPassword   = os.Args[4]
	)
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		log.Fatal("Missing tolen")
	}

	if mysqlServer == "" || mysqlServerPort == "" || mysqlUsername == "" || mysqlPassword == "" {
		log.Fatal("Missing database arguments")
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


	// Setup telegram bot
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := client.GetUpdatesChan(u)

	// Get the message from client
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//Get the response from database
		var response string
		row := db.QueryRow("SELECT response FROM chatbot.telegram  WHERE keyword = ?", update.Message.Text)
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,response)
		client.Send(msg)
	}

}
