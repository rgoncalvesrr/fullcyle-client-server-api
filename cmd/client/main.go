package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type cotacao struct {
	Bid string `json:"bid"`
}

func BuscaCotacao(ctx context.Context) (*cotacao, error) {

	select {
	case <-ctx.Done():
		log.Println("Tempo limite atingido")
		return nil, errors.New("requisição cancelada")
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
		if err != nil {
			return nil, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Erro ao fazer a requisição", err)
			return nil, err
		}
		defer res.Body.Close()

		r, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var body cotacao

		if err = json.Unmarshal(r, &body); err != nil {
			return nil, err
		}

		return &body, nil
	}
}

func GravarCotacao(valor string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Dólar:%s", valor)
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cota, err := BuscaCotacao(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	GravarCotacao(cota.Bid)
}
