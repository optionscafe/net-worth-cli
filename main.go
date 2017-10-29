//
// Date: 10/18/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
  "os"
  "fmt"
  "log"
  "time"
  "os/user"
  "github.com/joho/godotenv"
  "github.com/olekukonko/tablewriter"  
)

var (
  timeZone *time.Location
)

//
// Main...
//
func main() {
    
  // Get the current user.
  usr, err := user.Current()
    
  if err != nil {
    log.Fatal(err)
  }
    
  // Load .env file 
  err = godotenv.Load(usr.HomeDir + "/.net-worth-cli")

  if err != nil {
    log.Fatal("Error loading ~/.net-worth-cli file")
  } 

  // Set location.
  timeZone, _ = time.LoadLocation(os.Getenv("TIMEZONE"))
     
  // Make sure we have at least one arg
  if len(os.Args) <= 1 {
    PrintHelp()
    return
  }

  // Switch based on the first argument
  switch os.Args[1] {

    // List Accounts
    case "accounts-list":
      DoAccounts()

    // Create Account
    case "accounts-create":
      DoCreateAccount()

    // Mark Account
    case "accounts-mark":
      MarkAccountValue()

    // List Marks
    case "marks-list":
      ListMarks()

     // Create ledger entry
    case "ledger-create":
      DoCreateLedger()     

     // List ledger entry
    case "ledger-list":
      DoLedgerList()  

    // Print Help
    case "help":
      PrintHelp()

    // Print Help
    default:
      PrintHelp()

  }

}

//
// Print help
//
func PrintHelp() {

  fmt.Println("")
  fmt.Println("Actions:")
  fmt.Println("\n help")
  fmt.Println("\n marks-list")
  fmt.Println("\n accounts-list")
  fmt.Println("\n accounts-list {id}")
  fmt.Println("\n accounts-mark {id} {balance}")   
  fmt.Println("\n accounts-create \"{name}\" {balance}")
  fmt.Println("\n ledger-list")  
  fmt.Println("\n ledger-create {account_id} {date} {amount} \"{category_name}\" \"{note}\"")    
  fmt.Println("")
}

//
// Print table.
//
func PrintTable(rows [][]string, headers []string) {

  // Print table to screen.
  table := tablewriter.NewWriter(os.Stdout) 

  // Build table headers
  table.SetHeader(headers) 

  // Build table rows
  for _, v := range rows {
    table.Append(v)
  } 

  // Send output 
  table.Render()
}

/* End File */