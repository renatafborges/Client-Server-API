<h1 align="center">
   <a href="#"> Client Server API </a>
</h1>

<h3 align="center">
    In this challenge I applied what I learned about http webserver, contexts, database and file manipulation with Golang.
</h3>

<h4 align="center"> 
	 Status: Finished
</h4>

<p align="center"> 
 <a href="#how-it-works">How it works</a> • 
 <a href="#author">Author</a> • 

</p>

## How it works

- client/main.go make an HTTP request to server/main.go requesting the Dollar quote.
- server/main.go consume the API containing the Dollar and Real exchange rate at the address: https://economia.awesomeapi.com.br/json/last/USD-BRL and return the result to client in JSON format.
- using context package, server/main.go register each quote received in the SQLite database, with the maximum timeout to call the Dollar quote API being 200ms and the maximum timeout to be able to persist the data in the database being 10ms.
- client/main.go will only receive the current exchange rate value from server/main.go ("bid" field in the JSON).
- using context package, client/main.go have a maximum timeout of 300ms to receive the result from server/main.go.
- the 3 contexts return an error in logs if the execution time is insufficient.
- client/main.go will save the current quote in "cotacao.txt" file in format: Dollar: {value}
- the endpoint generated by server/main.go is: /cotacao and the port used by HTTP server is 8080.
  
---

This project is divided into one part:
1. Backend

### Pre-requisites

Before you begin, you will need to have the following tools installed on your machine:
[Git] (https://git-scm.com), 
[Golang] (https://go.dev/doc/install)
In addition, it is good to have an editor to work with the code like [VSCode] (https://code.visualstudio.com/)

#### Running Project

```bash

# Clone this repository
$ git clone https://github.com/renatafborges/Client-Server-API.git

# Access the project folder cmd/terminal
$ cd Client-Server-API

# go to the server folder
$ cd server

# in server folder run main.go
$ go run main.go

# The server will start at port: 8080 with the following message
Server is running on :8080

# open another terminal tab and access client folder
(to verify the current folder)
$ ls
(to move up a folder level)
$ cd ..
(access client folder)
$ cd client 

# in client folder run main.go
$ go run main.go

# The following message will appear at terminal in case of success
200
File created with success!

# in client folder the file will be created
cotacao.txt

# you may delete this file in case of another test
# you can use the extension SQLite Viewer and SQLite to access bid.sqlite database in server/bid.sqlite
```
## Author
Made with love by Renata Borges 👋🏽 [Get in Touch!](Https://www.linkedin.com/in/renataborgestech)
---
