package main

import (
	"get_rating_card/formatter"
)

func main() {
	fileData := formatter.ReaderCSV("/Users/egororlov/Desktop/get_rating_card/users.csv")
	newFileCSV := formatter.Getter(fileData)
	formatter.Writer(newFileCSV)
}
