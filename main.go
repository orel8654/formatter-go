package main

import (
	"get_rating_card/formatter"
)

func main() {
	var path string
	fileData := formatter.ReaderCSV(path)
	newFileCSV := formatter.Getter(fileData)
	formatter.Writer(newFileCSV)
}
