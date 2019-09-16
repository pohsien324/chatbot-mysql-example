

# Facebook Messenger Bot echo-reply example with MySQL connection
The simply echo-reply example for messenger bot. In this instance, the webhook server will connect to MySQL service and get the matched response for the incoming message. You can easily create events by inserting the keyword/response record into MySQL server.

## Preparation

Make sure you have imported the following packages:
1. [paked/messenger](https://github.com/paked/messenger)
2. [go-sql-driver](https://github.com/go-sql-driver/mysql)

## How to execute?
```{bash}
$ go get github.com/pohsienshih/chatbot-mysql-example/messenger
```
```{bash}
$ export VERIFY_TOKEN=<your verify token>
$ export PAGE_TOKEN=<your page token>
$ export PORT=<the port you want to listen on>

$ cd $GOPATH/src/pohsienshih/chatbot-mysql-example/messenger
$ go build -o webhook .
$ ./webhook <mysql server ip> <mysql server port> <mysql server username> <mysql server password>
```
> Make sure you already have MySQL service.

## Notice
TLS connection for this example is not yet supported. You can expose your service by using [ngrok](https://ngrok.com/).
```{bash}
$ ngrok http <port>
```

