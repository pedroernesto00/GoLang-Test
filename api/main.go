package main

import (
	"encoding/json"
	"log"
	"net/http"
	//"fmt"
	"io/ioutil"
	"github.com/gorilla/mux"
	"math/big"
	"math"
)

type api struct {
	Page int `json:page`
	TotalPages int `json:totalPages`
	ItemsOnPage int `json:itemsOnPage`
	Address string  `json:address`
	Balance string `json:balance`
	TotalReceived string  `json:totalReceived`
	TotalSent string `json:totalSent`
	UnconfirmedBalance string `json:unconfirmedBalance`
	UnconfirmedTxs int `json:unconfirmedTxs`
	Txs int `json:txs`
	TxsIds []string `json:txids`

}

type Balance struct {
	Balance string
}

func GetBalance(w http.ResponseWriter, r *http.Request) {

	//Getting the json from api
	resp, err := http.Get("https://blockbook-bitcoin.tronwallet.me/api/v2/address/" + r.URL.Path[1:])

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	content := api{}

	err = json.Unmarshal(body, &content)
	if err != nil {
		log.Fatal(err)
	}
	
	//Dividing the balance by 10**8
	real_balance := new(big.Float)
	divisor := new(big.Float).SetFloat64(math.Pow(10,8))
	real_balance.SetString(content.Balance)
	real_balance.Quo(real_balance, divisor)

	//Writing the resulting json
	bal := Balance{}
	bal.Balance = real_balance.String()

	balanceJson, err := json.Marshal(bal)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(balanceJson)
	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/balance/{address}", GetBalance).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}