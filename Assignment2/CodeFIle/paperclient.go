package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	command := os.Args[1]
	server := os.Args[2]

	switch command {
	case "add":
		if len(os.Args) != 6 {
			fmt.Println("Usage: paperclient add <server> <author> <title> <file>")
			return
		}
		addPaper(server, os.Args[3], os.Args[4], os.Args[5])

	case "list":
		listPapers(server)

	case "detail":
		if len(os.Args) != 4 {
			fmt.Println("Usage: paperclient detail <server> <paper-id>")
			return
		}
		showDetails(server, os.Args[3])

	case "fetch":
		if len(os.Args) != 4 {
			fmt.Println("Usage: paperclient fetch <server> <paper-id>")
			return
		}
		fetchPaper(server, os.Args[3])

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  paperclient add <server> <author> <title> <file>")
	fmt.Println("  paperclient list <server>")
	fmt.Println("  paperclient detail <server> <paper-id>")
	fmt.Println("  paperclient fetch <server> <paper-id>")
}

func connect(server string) *rpc.Client {
	client, err := rpc.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	return client
}

func addPaper(server, author, title, filePath string) {
	client := connect(server)
	defer client.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	format := filepath.Ext(filePath)
	if len(format) > 0 {
		format = format[1:]
	}

	args := AddPaperArgs{
		Author:  author,
		Title:   title,
		Format:  format,
		Content: content,
	}

	var reply AddPaperReply
	if err := client.Call("PapersServer.AddPaper", &args, &reply); err != nil {
		log.Fatalf("AddPaper failed: %v", err)
	}

	fmt.Println("Paper uploaded successfully!")
	fmt.Println("Paper ID:", reply.PaperNumber)
}

func listPapers(server string) {
	client := connect(server)
	defer client.Close()

	var reply ListPapersReply
	if err := client.Call("PapersServer.ListPapers", &ListPapersArgs{}, &reply); err != nil {
		log.Fatalf("ListPapers failed: %v", err)
	}

	fmt.Printf("%-5s %-20s %-30s\n", "ID", "Author", "Title")
	for _, p := range reply.Papers {
		fmt.Printf("%-5d %-20s %-30s\n", p.PaperNumber, p.Author, p.Title)
	}
}

func showDetails(server, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Invalid paper ID")
	}

	client := connect(server)
	defer client.Close()

	var reply GetPaperDetailsReply
	if err := client.Call("PapersServer.GetPaperDetails",
		&GetPaperArgs{PaperNumber: id}, &reply); err != nil {
		log.Fatalf("GetPaperDetails failed: %v", err)
	}

	fmt.Println("Author:", reply.Author)
	fmt.Println("Title :", reply.Title)
	fmt.Println("Format:", reply.Format)
}

func fetchPaper(server, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Invalid paper ID")
	}

	client := connect(server)
	defer client.Close()

	var reply FetchPaperReply
	if err := client.Call("PapersServer.FetchPaperContent",
		&FetchPaperArgs{PaperNumber: id}, &reply); err != nil {
		log.Fatalf("FetchPaperContent failed: %v", err)
	}

	fileName := fmt.Sprintf("paper_%d.%s", id, reply.Format)
	err = os.WriteFile(fileName, reply.Content, 0644)
	if err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}

	fmt.Println("Paper saved as:", fileName)
}
