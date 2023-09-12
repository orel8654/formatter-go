package formatter

//Добавление
//error wrap
//multierr
//time.Time при чтении

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
)

const Buffer = 3
const (
	ColID = iota
	ColLastLogin
	ColDateJoined
	ColUsername
	ColFirstName
	ColPhone
	ColCountCard
	ColProductID
)

type FileCSV struct {
	ID,
	LastLogin,
	DateJoined,
	Username,
	FirstName,
	Phone,
	CountCard,
	ProductID,
	url string
	Rating any
}

type ResponseData struct {
	AvItemID int64   `json:"av_item_id"`
	AvComps  int64   `json:"av_comps"`
	AvRating float64 `json:"av_rating"`
}

func fetcByURL(url string) (data ResponseData, err error) {
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return ResponseData{}, err
	}
	defer response.Body.Close()
	return data, json.NewDecoder(response.Body).Decode(&data)
}

func Fetch(id int, work <-chan FileCSV, result chan<- FileCSV) {
	for c := range work {
		if c.ProductID == "" {
			c.Rating = 0
			result <- c
			continue
		}
		data, err := fetcByURL(c.url)
		if err != nil {
			fmt.Println(err)
		}
		c.Rating = data.AvRating
		result <- c
	}
}

func Getter(data []FileCSV) []FileCSV {
	var count int64
	worker := make(chan FileCSV, Buffer)
	result := make(chan FileCSV, Buffer)

	newFileCSV := make([]FileCSV, 0, len(data))

	for i := 1; i <= Buffer; i++ {
		go Fetch(i, worker, result)
	}

	for _, val := range data {
		val.url = "http://127.0.0.1:8000/api/v3/rating?marketplace=wb&item_id=" + val.ProductID
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

func ReaderCSV(path string) ([]FileCSV, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]FileCSV, 0, len(records)-1)
	for i, line := range records {
		if i == 0 || len(line) == 0 {
			continue
		}
		rec := FileCSV{
			ID:         line[ColID],
			LastLogin:  line[ColLastLogin],
			DateJoined: line[ColDateJoined],
			Username:   line[ColUsername],
			FirstName:  line[ColFirstName],
			Phone:      line[ColPhone],
			CountCard:  line[ColCountCard],
			ProductID:  line[ColProductID],
		}
		result = append(result, rec)
	}
	return result, nil
}

func Writer(data []FileCSV) error {
	file, err := os.Create("newResult.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&data, file); err != nil {
		return err
	}
	return nil
}
