//
// Date: 10/20/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "log"
  "time"
  "github.com/tidwall/gjson"  
  "github.com/leekchan/accounting"  
)

//
// List all marks
//
func ListMarks() {

  // Set output data.
  var rows [][]string

  // Set money format
  ac := accounting.Accounting{Symbol: "$", Precision: 2}

  // Make API request
  body, err := MakeGetRequest("/api/v1/marks")

  if err != nil {
    log.Fatal(err)
  } 

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
