//
// Date: 10/20/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "github.com/leekchan/accounting"
  "github.com/tidwall/gjson"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "time"
)

//
// List all marks
//
func ListMarks() {

  // Set output data.
  var rows [][]string

  // Set money format
  ac := accounting.Accounting{Symbol: "$", Precision: 2}

  // Setup http client
  client := &http.Client{}

  // Setup api request
  req, _ := http.NewRequest("GET", os.Getenv("SERVER_URL")+"/api/v1/marks", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

  res, err := client.Do(req)

  if err != nil {
    log.Fatal(err)
  }

  // Close Body
  defer res.Body.Close()

  // Make sure the api responded with a 200
  if res.StatusCode == 404 {
    log.Fatal("No results found.")
  }

  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body)

  // Loop through the accounts and print them
  result := gjson.Parse(string(body))

  // Loop through and build rows of output table.
  result.ForEach(func(key, value gjson.Result) bool {

    date := gjson.Get(value.String(), "date").String()
    units := gjson.Get(value.String(), "units").String()
    balance := gjson.Get(value.String(), "balance").Float()
    unitPrice := gjson.Get(value.String(), "price_per").Float()

    // Parse dates.
    layout := "2006-01-02T15:04:05Z"
    d, _ := time.Parse(layout, date)

    rows = append(rows, []string{d.Format("01/02/2006"), units, ac.FormatMoney(unitPrice), ac.FormatMoney(balance)})

    // keep iterating
    return true
  })

  // Print record.
  PrintTable(rows, []string{"Date", "Units", "Price", "Balance"})

}

/* End File */
