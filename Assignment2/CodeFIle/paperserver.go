package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PapersServer struct {
	mu           sync.Mutex
	papers       map[int]Paper
	nextPaperNum int
	rabbitCh     *amqp.Channel
	rabbitQueue  amqp.Queue
}

// rpc methods

func (s *PapersServer) AddPaper(args *AddPaperArgs, reply *AddPaperReply) error {
	if args.Author == "" || args.Title == "" {
		return fmt.Errorf("Invalid paper data")
	}

	s.mu.Lock()
	id := s.nextPaperNum
	s.nextPaperNum++

	paper := Paper{
		PaperNumber: id,
		Author:      args.Author,
		Title:       args.Title,
		Format:      args.Format,
		Content:     args.Content,
	}
	s.papers[id] = paper
	s.mu.Unlock()

	reply.PaperNumber = id
	reply.Success = true
	reply.Message = "Paper added successfully"

	msg := fmt.Sprintf("New Paper Added | ID=%d | Title=%s", id, args.Title)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = s.rabbitCh.PublishWithContext(
		ctx,
		"",
		s.rabbitQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)

	return nil
}

func (s *PapersServer) ListPapers(_ *ListPapersArgs, reply *ListPapersReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.papers {
		reply.Papers = append(reply.Papers, PaperInfo{
			PaperNumber: p.PaperNumber,
			Author:      p.Author,
			Title:       p.Title,
		})
	}
	return nil
}

func (s *PapersServer) GetPaperDetails(args *GetPaperArgs, reply *GetPaperDetailsReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.papers[args.PaperNumber]
	if !ok {
		return fmt.Errorf("Paper not found")
	}

	reply.Author = p.Author
	reply.Title = p.Title
	reply.Format = p.Format
	return nil
}

func (s *PapersServer) FetchPaperContent(args *FetchPaperArgs, reply *FetchPaperReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.papers[args.PaperNumber]
	if !ok {
		return fmt.Errorf("Paper not found")
	}

	reply.Content = p.Content
	reply.Format = p.Format
	return nil
}


func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Channel error:", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"paper_notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	server := &PapersServer{
		papers:      make(map[int]Paper),
		nextPaperNum: 1,
		rabbitCh:    ch,
		rabbitQueue: queue,
	}

	rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Server listen error:", err)
	}

	log.Println("Paper Server running on port 1234")
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}
