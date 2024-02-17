package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRLQuotationResponse struct {
	USDBRlQuotation quotation `json:"USDBRL"`
}

type quotation struct {
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	High      string `json:"high"`
	Low       string `json:"low"`
}

type bidQuotationOutput struct {
	Bid string `json:"Dólar"`
}

var db *sql.DB

const awesomeApiURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {

	var err error
	db, err = sql.Open("sqlite3", "./bid.sqlite")
	if err != nil {
		log.Fatal("could open database connection", err)
	}
	defer db.Close()

	sqlQuery := `create table IF NOT EXISTS quotation (bid text);`
	_, err = db.Exec(sqlQuery)
	if err != nil {
		log.Fatal("could not execute query", err)
	}

	http.HandleFunc("/cotacao", handler)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", awesomeApiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not make new request with context", slog.String("url", awesomeApiURL), slog.Any("error", err))
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not get quotation at awesomeapi", slog.String("url", awesomeApiURL), slog.Any("error", err))
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not read response body", slog.Any("body", res.Body), slog.Any("error", err))
		return
	}

	var q USDBRLQuotationResponse
	if err := json.Unmarshal(body, &q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not unmarshal quotation response", slog.String("body", string(body)), slog.Any("error", err))
		return
	}

	bidOutput := bidQuotationOutput{
		Bid: q.USDBRlQuotation.Bid,
	}

	ctx, cancel = context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not begin transaction", slog.Any("transaction", tx), slog.Any("error", err))
		return
	}

	_, err = tx.ExecContext(ctx, "insert into quotation(bid) values(?)", bidOutput.Bid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not execute insert quotation", slog.String("quotation", bidOutput.Bid), slog.Any("error", err))
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("could not commit transaction", slog.Any("transaction", tx), slog.Any("error", err))
		tx.Rollback()
		return
	}

	o, err := json.Marshal(bidOutput)
	if err != nil {
		fmt.Println("Error marshaling bidOutput:", err)
	}

	w.Write(o)
}
