package models

import (
	"encoding/xml"
)

// Struct for the API response (adjust according to the actual API response)
type Rates struct {
	XMLName xml.Name `xml:"rates"`
	Date    string   `xml:"date"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	Fullname    string `xml:"fullname"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Date        string `xml:"date"`
}

type ResponseItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
	Value string `json:"value"`
	Adate string `json:"adate"`
}

type Response struct {
	Success bool `json:"success"`
}

type UserRepository interface {
	GetByDate(Date string) ([]*ResponseItem, error)
	GetByDateCode(Date string, Code string) ([]*ResponseItem, error)
	Exists(user *Item) (int, error)
	Update(user *Item) error
	Insert(user *Item) error
}
