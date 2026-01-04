package main

import (
	"fmt"
	"paxosLab/paxos"
)

func main() {
	acceptors := []*paxos.Acceptor{
		{Alive: true},
		{Alive: true},
		{Alive: false},
		{Alive: false}, // failed
		{Alive: false},
	}

	proposer := paxos.Proposer{ProposalNumber: 1}
	value := proposer.Propose("Distributed Systems", acceptors)

	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}
}
