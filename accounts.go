//
// Date: 10/18/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "io"
  "os"
  "fmt"
  "log"
  "time"
  "bytes"
  "net/http"
  "io/ioutil"
  "github.com/tidwall/gjson"
  "github.com/leekchan/accounting"
  "github.com/olekukonko/tablewriter"  
)

//
// Process Accounts
//
func DoAccounts() {

  // List all accounts
  if len(os.Args) == 2 {
    AccountsList()
    return
  }

  // List just one account.
  if len(os.Args) == 3 {
    AccountList()
    return
  }

  PrintHelp()

}

//
// Create a new account.
//
func DoCreateAccount() {

  // Make sure we have the args we need.
  if len(os.Args) < 4 {
    PrintHelp()
    return
  }

  // Post data
  var postStr = []byte(`{"name":"` + os.Args[2] +  `","balance":` + os.Args[3] + `}`)

  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("POST", os.Getenv("SERVER_URL") + "/api/v1/accounts", bytes.NewBuffer(postStr))
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
    log.Fatal("We already have an account with the name " + os.Args[2])
  }  
  
  // Make sure the api responded with a 201
  if res.StatusCode != 201 {
    log.Fatal(fmt.Sprint("/api/v1/accounts (POST) did not return a status code of 201 -", res.StatusCode))
  }   

  // Print record.
  PrintOneAccountRow(res.Body)
}

//
// List an account by id.
//
func AccountList() {
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", os.Getenv("SERVER_URL") + "/api/v1/accounts/" + os.Args[2], nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + os.Getenv("ACCESS_TOKEN"))   
 
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

  // Print record.
  PrintOneAccountRow(res.Body)
}

//
// List accounts
//
func AccountsList() {
  
  // Set output data.
  var rows [][]string

  // Keep track of the total.
  var total = 0.00;

  // Set money format
  ac := accounting.Accounting{Symbol: "$", Precision: 2}  

  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", os.Getenv("SERVER_URL") + "/api/v1/accounts", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + os.Getenv("ACCESS_TOKEN"))   
 
  res, err := client.Do(req)
      
  if err != nil {
    log.Fatal(err)  
  }        
  
  // Close Body
  defer res.Body.Close()    
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    log.Fatal(fmt.Sprint("/api/v1/accounts did not return a status code of 200 -", res.StatusCode))
  }    
     
  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body) 

  // Loop through the accounts and print them
  result := gjson.Parse(string(body))

  // Loop through and build rows of output table.
  result.ForEach(func(key, value gjson.Result) bool {
        
    id := gjson.Get(value.String(), "id").String()
    name := gjson.Get(value.String(), "name").String()
    balance := gjson.Get(value.String(), "balance").Float()
    units := gjson.Get(value.String(), "units").String()    
    createdAt := gjson.Get(value.String(), "created_at").String()
    updatedAt := gjson.Get(value.String(), "updated_at").String()    

    // Parse dates.
    layout := "2006-01-02T15:04:05Z"
    c, _ := time.Parse(layout, createdAt)
    u, _ := time.Parse(layout, updatedAt)

    // Keep track of the total balances.
    total = total + balance

    rows = append(rows, []string{ id, name, ac.FormatMoney(balance), units, c.In(timeZone).Format("01/02/2006"), u.In(timeZone).Format("01/02/2006") })

    // keep iterating
    return true 
  })
 
  // Print table to screen.
  table := tablewriter.NewWriter(os.Stdout) 

  // Build table headers
  table.SetHeader([]string{ "Id", "Name", "Balance", "Units", "Created At", "Updated At" }) 

  // Build table rows
  for _, v := range rows {
    table.Append(v)
  } 

  // Set footer
  table.SetFooter([]string{"", "Total", ac.FormatMoney(total), "", "", ""}) 

  // Send output 
  table.Render()
}

//
// Mark an account's value
//
func MarkAccountValue() {

  // Make sure we have the args we need.
  if len(os.Args) < 4 {
    PrintHelp()
    return
  }

  // Post data
  var postStr = []byte(`{"balance":` + os.Args[3] +  `,"date":"` + time.Now().UTC().Format("2006-01-02") + `"}`)

  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("POST", os.Getenv("SERVER_URL") + "/api/v1/accounts/" + os.Args[2] + "/marks", bytes.NewBuffer(postStr))
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + os.Getenv("ACCESS_TOKEN"))   
 
  res, err := client.Do(req)
      
  if err != nil {
    log.Fatal(err)  
  }        
  
  // Close Body
  defer res.Body.Close()    

  // See if there was an error
  if res.StatusCode != 201 {
    log.Fatal("There was an error marking this asset.")
  }

  // Print result
  fmt.Println("") 
  fmt.Println("Account as been marked at $" + os.Args[3] + ".")

  // Print record.
  PrintOneAccountRow(res.Body)
}

//
// Print one account row.
//
func PrintOneAccountRow(resBody io.ReadCloser) {

  // Set output data.
  var rows [][]string

  // Set money format
  ac := accounting.Accounting{Symbol: "$", Precision: 2}

  // Read the data we got.
  body, _ := ioutil.ReadAll(resBody)

  // Get the values we need.
  id := gjson.Get(string(body), "id").String()
  name := gjson.Get(string(body), "name").String()
  balance := gjson.Get(string(body), "balance").Float()
  units := gjson.Get(string(body), "units").String()    
  createdAt := gjson.Get(string(body), "created_at").String()
  updatedAt := gjson.Get(string(body), "updated_at").String()

  // Parse dates.
  layout := "2006-01-02T15:04:05Z"
  c, _ := time.Parse(layout, createdAt)
  u, _ := time.Parse(layout, updatedAt)

  rows = append(rows, []string{ id, name, ac.FormatMoney(balance), units, c.In(timeZone).Format("01/02/2006"), u.In(timeZone).Format("01/02/2006") })

  fmt.Println("")

  // Print table and return.
  PrintTable(rows, []string{ "Id", "Name", "Balance", "Units", "Created At", "Updated At" })

  fmt.Println("")
}

/* End File */