package main

// data model

type Paper struct {
	PaperNumber int
	Author      string
	Title       string
	Format      string
	Content     []byte
}

// request and response structs

type AddPaperArgs struct {
	Author  string
	Title   string
	Format  string
	Content []byte
}

type AddPaperReply struct {
	PaperNumber int
	Success     bool
	Message     string
}

type ListPapersArgs struct{}

type PaperInfo struct {
	PaperNumber int
	Author      string
	Title       string
}

type ListPapersReply struct {
	Papers []PaperInfo
}

type GetPaperArgs struct {
	PaperNumber int
}

type GetPaperDetailsReply struct {
	Author string
	Title  string
	Format string
}

type FetchPaperArgs struct {
	PaperNumber int
}

type FetchPaperReply struct {
	Content []byte
	Format  string
}
