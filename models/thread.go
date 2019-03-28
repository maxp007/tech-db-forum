package models

type Thread struct {
	description string
	Author      string `json:"author, string"`
	Created     string `json:"created,string"`
	Forum       string
	Id          int32
	Message     string
	Slug        string
	Title       string
	Votes       int32
}
