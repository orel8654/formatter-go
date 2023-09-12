package main

import (
	"fmt"
	"get_rating_card/formatter"
	"os"
)

func main() {
	var path string
	fileData, err := formatter.ReaderCSV(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	newFileCSV := formatter.Getter(fileData)
	formatter.Writer(newFileCSV)
}
