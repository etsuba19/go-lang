package main

import (
	"encoding/json"
	"net/http"
	"os"
	"paxosLab/paxos"
)

var acceptor = &paxos.Acceptor{
	Alive: true,
}

func prepareHandler(w http.ResponseWriter, r *http.Request) {
	var p paxos.Prepare
	json.NewDecoder(r.Body).Decode(&p)

	promise := acceptor.HandlePrepare(p)
	json.NewEncoder(w).Encode(promise)
}

func acceptHandler(w http.ResponseWriter, r *http.Request) {
	var a paxos.Accept
	json.NewDecoder(r.Body).Decode(&a)

	accepted := acceptor.HandleAccept(a)
	json.NewEncoder(w).Encode(accepted)
}

func main() {
	port := os.Args[1]
	http.HandleFunc("/prepare", prepareHandler)
	http.HandleFunc("/accept", acceptHandler)
	http.ListenAndServe(":"+port, nil)
}
