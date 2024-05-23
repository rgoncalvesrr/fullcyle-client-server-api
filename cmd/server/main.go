package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rgoncalvesrr/fullcyle-client-server-api/internal/core/entity"
)

type DataBody struct {
	Data entity.Cotacao `json:"USDBRL"`
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(w, "Erro ao fazer a requisição %v", err)
	}
	defer req.Body.Close()
	_body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "Erro ao ler stream %v", err)
	}
	var _dataBody DataBody

	if err = json.Unmarshal(_body, &_dataBody); err != nil {
		fmt.Fprintf(w, "Erro ao fazer decodificar json %v", err)
	}

	_json, err := json.Marshal(_dataBody.Data)
	if err != nil {
		fmt.Fprintf(w, "Erro codificar json %v", err)
	}

	w.Write(_json)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BuscaCotacaoHandler)
	http.ListenAndServe(":3000", mux)
}
