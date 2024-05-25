package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/glebarez/go-sqlite"

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
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	cota, err := BuscaCotacao(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	if err = SalvaCotacao(ctx, cota); err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(CotacaoOutputDTO{Bid: cota.Bid})

}

func SalvaCotacao(ctx context.Context, c *entity.Cotacao) error {
	select {
	case <-ctx.Done():
		return errors.New("tempo excedido para fazer a gravação no banco de dados")
	default:
		db, err := sql.Open("sqlite", "./cotacoes.db")
		if err != nil {
			return errors.Join(errors.New("erro ao abrir banco de dados"), err)
		}
		defer db.Close()

		_, err = db.ExecContext(ctx,
			`create table if not exists cotacoes(
				code  	  	varchar, 
				code_in    	varchar, 
				name 	   	varchar,
				high       	varchar,
				low        	varchar,
				var_bid     varchar,
				pct_change  varchar,
				bid        	varchar,
				ask        	varchar,
				time_stamp  varchar,
				create_date varchar);`)
		if err != nil {
			return errors.Join(errors.New("erro criar estrutura do banco de dados"), err)
		}

		stmt, err := db.PrepareContext(ctx,
			`insert into cotacoes 
				(code, code_in, name, high, low, var_bid, pct_change, bid, ask, time_stamp, create_date)
			values
				(?,?,?,?,?,?,?,?,?,?,?)`)

		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, c.Code, c.CodeIn, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask, c.TimeStamp, c.CreateDate)

		if err != nil {
			return err
		}
	}
	return nil
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
