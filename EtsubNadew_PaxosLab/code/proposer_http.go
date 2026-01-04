package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"paxosLab/paxos"
)

var acceptorURLs = []string{
	"http://localhost:8001",
	"http://localhost:8002",
	"http://localhost:8003",
}

func main() {
	proposal := paxos.Prepare{ProposalNumber: 1}
	promises := 0

	for _, url := range acceptorURLs {
		body, _ := json.Marshal(proposal)
		resp, err := http.Post(url+"/prepare", "application/json", bytes.NewBuffer(body))
		if err != nil {
			continue
		}

		var promise paxos.Promise
		json.NewDecoder(resp.Body).Decode(&promise)

		if promise.ProposalNumber == proposal.ProposalNumber {
			promises++
		}
	}

	if promises <= len(acceptorURLs)/2 {
		fmt.Println("Consensus not reached (prepare phase)")
		return
	}

	accepted := 0
	accept := paxos.Accept{
		ProposalNumber: 1,
		Value:          "Distributed Systems",
	}

	for _, url := range acceptorURLs {
		body, _ := json.Marshal(accept)
		resp, err := http.Post(url+"/accept", "application/json", bytes.NewBuffer(body))
		if err != nil {
			continue
		}

		var ack paxos.Accepted
		json.NewDecoder(resp.Body).Decode(&ack)

		if ack.ProposalNumber == accept.ProposalNumber {
			accepted++
		}
	}

	if accepted > len(acceptorURLs)/2 {
		fmt.Println("Consensus reached on value: Distributed Systems")
	} else {
		fmt.Println("Consensus not reached (accept phase)")
	}
}
