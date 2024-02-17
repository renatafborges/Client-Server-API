package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type bidQuotationOutput struct {
	Bid string `json:"Dólar"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusOK {
		return
	}

	fmt.Println(res.StatusCode)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("could not read response body", slog.Any("body", res.Body), slog.Any("error", err))
		return
	}

	var b bidQuotationOutput
	if err := json.Unmarshal(body, &b); err != nil {
		slog.Error("could not unmarshal quotation response", slog.String("body", string(body)), slog.Any("error", err))
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
	}

	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", b.Bid))
	fmt.Println("File created with success!")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing on file: %v\n", err)
	}

	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
