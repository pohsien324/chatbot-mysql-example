# Line Bot echo-reply example with MySQL connection
The simply echo-reply example for Line bot. In this example, the webhook server will connect to MySQL service and get the matched response for the incoming message. You can easily create events by inserting the keyword/response record into MySQL server.

## Preparation

Make sure you have imported the following packages:
1. [line-bot-sdk-go](https://github.com/line/line-bot-sdk-go)
2. [go-sql-driver](https://github.com/go-sql-driver/mysql)

## How to execute?
```{bash}
$ go get github.com/pohsienshih/chatbot-mysql-example/line
```
```{bash}
$ export CHANNEL_SECRET=<yoursecret>
$ export CHANNEL_TOKEN=<yourtoken>
$ export PORT=<the port you want to listen on>

$ cd $GOPATH/src/pohsienshih/chatbot-mysql-example/line
$ go build -o webhook .
$ ./webhook <mysql server ip> <mysql server port> <mysql server username> <mysql server password>
```
> Make sure you already have MySQL service.

## Notice
TLS connection for this example is not yet supported. You can expose your service by using [ngrok](https://ngrok.com/).
```{bash}
$ ngrok http <port>
```

