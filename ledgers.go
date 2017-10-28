//
// Date: 10/28/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "os"
  "io"  
  "log" 
  "fmt"
  "time" 
  "bytes"  
  "net/http"
  "io/ioutil"
  "github.com/tidwall/gjson"
  "github.com/leekchan/accounting"   
)

//
// Create a new ledger.
//
func DoCreateLedger() {

  // Make sure we have the args we need.
  if len(os.Args) < 6 {
    PrintHelp()
    return
  }

  // Post data
  var postStr = []byte(`{"date":"` + os.Args[3] + `","amount":` + os.Args[4] + `,"account_id":` + os.Args[2] + `,"note":"` + os.Args[5] + `"}`);

  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("POST", os.Getenv("SERVER_URL") + "/api/v1/ledgers", bytes.NewBuffer(postStr))
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + os.Getenv("ACCESS_TOKEN"))   
 
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
  amount := gjson.Get(string(body), "amount").Float()
  date := gjson.Get(string(body), "date").String() 
  note := gjson.Get(string(body), "note").String() 

  // Parse dates.
  layout := "2006-01-02T15:04:05Z"
  d, _ := time.Parse(layout, date)

  rows = append(rows, []string{ id, d.In(timeZone).Format("01/02/2006"), account_name, ac.FormatMoney(amount), note })

  fmt.Println("")

  // Print table and return.
  PrintTable(rows, []string{ "Id", "Date", "Account", "Amount", "Note" })

  fmt.Println("")
}

/* End File */
