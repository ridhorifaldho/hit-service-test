package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/subosito/gotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ReqOrder struct {
	RequestID int64    `json:"request_id"`
	Data      []*Order `json:"data"`
}

type Order struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	RequestID int64     `json:"request_id"`
	Customer  string    `json:"customer"`
	Quantity  uint      `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func init() {
	gotenv.Load()
}

func main() {

	// URL endpoint API yang akan dihit
	url := os.Getenv("URL")
	fmt.Println("URL : ", url)
	// Jumlah data yang ingin diambil
	totalData, _ := strconv.Atoi(os.Getenv("TOTAL_LOOP"))
	fmt.Println("total looping", totalData)

	for i := 1; i <= totalData; i++ {

		// Membuat request http
		order := Order{
			Customer:  "Joko " + strconv.Itoa(i),
			Quantity:  10 + uint(i),
			Price:     100.00,
			Timestamp: time.Now(),
		}

		// Membuat ReqOrder struct
		reqOrder := ReqOrder{
			RequestID: int64(i),
			Data:      []*Order{&order},
		}

		// Mengencode ReqOrder struct ke dalam format JSON
		jsonData, err := json.Marshal(reqOrder)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Membuat request http
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Menambahkan header content-type JSON
		req.Header.Set("Content-Type", "application/json")

		// Membuat http client dan mengirimkan request
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyString := string(body)

		fmt.Println("Response Status :", resp.Status, "| Response data : ", bodyString)
	}
}
