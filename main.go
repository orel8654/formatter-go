package main

import (
	"fmt"
	"get_rating_card/formatter"
	"os"
)

func main() {
	var path string = "/Users/egororlov/Desktop/get_rating_card/users.csv"
	fileData, err := formatter.ReaderCSV(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	newFileCSV := formatter.Getter(fileData)
	formatter.Writer(newFileCSV)
}
