package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rgoncalvesrr/fullcyle-client-server-api/internal/core/entity"
)

type resp struct {
	Data entity.Cotacao `json:"USDBRL"`
}

func main() {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição %v", err)
	}
	defer req.Body.Close()
	_body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler stream %v", err)
	}
	var r resp

	if err = json.Unmarshal(_body, &r); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer decodificar json %v", err)
	}
	fmt.Printf("%+v", r.Data)
}
