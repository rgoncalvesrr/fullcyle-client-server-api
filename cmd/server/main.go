package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rgoncalvesrr/fullcyle-client-server-api/internal/core/entity"
)

type DataBody struct {
	Data entity.Cotacao `json:"USDBRL"`
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {

	cota, err := BuscaCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cota)
}

func BuscaCotacao() (*entity.Cotacao, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	r, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var body DataBody

	if err = json.Unmarshal(r, &body); err != nil {
		return nil, err
	}

	return &body.Data, nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BuscaCotacaoHandler)
	http.ListenAndServe(":3000", mux)
}
