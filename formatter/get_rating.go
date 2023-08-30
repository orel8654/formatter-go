package formatter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
)

const Maxi int = 3

type FileCSV struct {
	Id         string
	LastLogin  string
	DateJoined string
	Username   string
	FirstName  string
	Phone      string
	CountCard  string
	ProductId  string
	Rating     any
	url        string
}

type ResponseData struct {
	AvItemId int64   `json:"av_item_id"`
	AvComps  int64   `json:"av_comps"`
	AvRating float64 `json:"av_rating"`
}

func NewFetch(id int, work <-chan FileCSV, result chan<- FileCSV) {
	for c := range work {
		if c.ProductId == "" {
			c.Rating = 0
			result <- c
		} else {
			fmt.Println(c.url)
			response, err := http.Get(c.url)
			c.Rating = 0
			if err != nil {
				fmt.Println(err)
				result <- c
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				result <- c
			}
			var data ResponseData
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
				result <- c
			}
			c.Rating = data.AvRating
			result <- c
		}
	}
}

func Getter(data []FileCSV) []FileCSV {
	var count int64
	worker := make(chan FileCSV, len(data))
	result := make(chan FileCSV, len(data))

	newFileCSV := []FileCSV{}

	for i := 1; i <= Maxi; i++ {
		go NewFetch(i, worker, result)
	}

	for _, val := range data {
		val.url = "http://127.0.0.1:8000/api/v3/rating?marketplace=wb&item_id=" + string(val.ProductId)
		worker <- val
	}
	close(worker)

	for i := 0; i < len(data); i++ {
		newFileCSV = append(newFileCSV, <-result)
		count++
		fmt.Printf("%d / %d\n", count, len(data))
	}
	return newFileCSV
}

func ReaderCSV(path string) []FileCSV {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	result := []FileCSV{}
	for i, line := range records {
		if i > 0 && len(line) != 0 {
			rec := FileCSV{Id: line[0], LastLogin: line[1], DateJoined: line[2], Username: line[3], FirstName: line[4], Phone: line[5], CountCard: line[6], ProductId: line[7]}
			result = append(result, rec)
		}
	}
	return result
}

func Writer(data []FileCSV) {
	file, err := os.Create("newResult.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&data, file); err != nil {
		panic(err)
	}
}
