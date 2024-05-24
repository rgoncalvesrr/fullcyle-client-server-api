package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rgoncalvesrr/fullcyle-client-server-api/internal/core/entity"
)

type DataBody struct {
	Data entity.Cotacao `json:"USDBRL"`
}

type CotacaoOutputDTO struct {
	Bid string `json:"bid"`
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {

	defer fmt.Println("Requisição finalizada")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cota, err := BuscaCotacao(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(CotacaoOutputDTO{Bid: cota.Bid})

}

func BuscaCotacao(ctx context.Context) (*entity.Cotacao, error) {
	select {
	case <-ctx.Done():
		log.Println("Tempo limite atingido")
		return nil, errors.New("requisição cancelada")
	default:
		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

		var body DataBody

		if err = json.Unmarshal(r, &body); err != nil {
			return nil, err
		}

		return &body.Data, nil
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", mux)
}
