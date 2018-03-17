//
// Date: 10/28/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leekchan/accounting"
	"github.com/tidwall/gjson"
)

//
// Do Ledger list.
//
func DoLedgerList() {

	// Set output data.
	var rows [][]string

	// Set money format
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	// Make API request
	body, err := MakeGetRequest("/api/v1/ledgers")

	if err != nil {
		log.Fatal(err)
	}

	// Loop through the accounts and print them
	result := gjson.Parse(body)

	// Loop through and build rows of output table.
	result.ForEach(func(key, value gjson.Result) bool {

		// Get values from json
		id := gjson.Get(value.String(), "id").String()
		account_name := gjson.Get(value.String(), "account_name").String()
		category_name := gjson.Get(value.String(), "category_name").String()
		amount := gjson.Get(value.String(), "amount").Float()
		date := gjson.Get(value.String(), "date").String()
		note := gjson.Get(value.String(), "note").String()

		// Parse dates.
		layout := "2006-01-02T15:04:05Z"
		d, _ := time.Parse(layout, date)

		rows = append(rows, []string{id, d.In(timeZone).Format("01/02/2006"), account_name, category_name, ac.FormatMoney(amount), note})

		// keep iterating
		return true
	})

	fmt.Println("")

	// Print table and return.
	PrintTable(rows, []string{"Id", "Date", "Account", "Category", "Amount", "Note"})

	fmt.Println("")
}

//
// Create a new ledger.
//
func DoCreateLedger() {

	// Make sure we have the args we need.
	if len(os.Args) < 8 {
		PrintHelp()
		return
	}

	// Post data
	var postStr = []byte(`{"date":"` + os.Args[3] + `","amount":` + os.Args[4] + `,"account_id":` + os.Args[2] + `,"category_name":"` + os.Args[5] + `","symbol":"` + os.Args[6] + `", "note":"` + os.Args[7] + `"}`)

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("POST", os.Getenv("SERVER_URL")+"/api/v1/ledgers", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	// Close Body
	defer res.Body.Close()

	// If 400 it is a duplicate entry.
	if res.StatusCode == 400 {
		log.Fatal("Something went wrong creating this ledger entry.")
	}

	// Make sure the api responded with a 201
	if res.StatusCode != 201 {
		log.Fatal(fmt.Sprint("/api/v1/ledgers (POST) did not return a status code of 201 -", res.StatusCode))
	}

	// Print record.
	PrintOneLedgerRow(res.Body)
}

//
// Print one Ledger row.
//
func PrintOneLedgerRow(resBody io.ReadCloser) {

	// Set output data.
	var rows [][]string

	// Set money format
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	// Read the data we got.
	body, _ := ioutil.ReadAll(resBody)

	// Get the values we need.
	id := gjson.Get(string(body), "id").String()
	account_name := gjson.Get(string(body), "account_name").String()
	category_name := gjson.Get(string(body), "category_name").String()
	amount := gjson.Get(string(body), "amount").Float()
	date := gjson.Get(string(body), "date").String()
	note := gjson.Get(string(body), "note").String()

	// Parse dates.
	layout := "2006-01-02T15:04:05Z"
	d, _ := time.Parse(layout, date)

	rows = append(rows, []string{id, d.In(timeZone).Format("01/02/2006"), account_name, category_name, ac.FormatMoney(amount), note})

	fmt.Println("")

	// Print table and return.
	PrintTable(rows, []string{"Id", "Date", "Account", "Category", "Amount", "Note"})

	fmt.Println("")
}

/* End File */
