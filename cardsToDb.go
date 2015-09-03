package main

import (
  "encoding/csv"
  "os"
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/caneroj1/hush"
  "strconv"
)

func getCSVFile(file string) *os.File {
  f, err := os.Open(file)
  if err != nil {
    panic(err)
  }

  return f
}

func getDBConnection(secrets hush.Hush) *sql.DB {
  fmt.Println("Initiating db connection.")

  // get dbname
  dbname, ok := secrets.GetString("dbname")
  if !ok {
    panic("Could not get database name.")
  }

  // get user
  user, ok := secrets.GetString("user")
  if !ok {
    panic("Could not get user.")
  }

  // get password
  password, ok := secrets.GetString("password")
  if !ok {
    panic("Could not get password.")
  }

  connectionString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, password)
  db, err := sql.Open("postgres", connectionString)
  if err != nil {
    panic(err)
  }

  fmt.Println("Successful")
  return db
}

// WriteToDB writes the records in the CSV to the db
func WriteToDB() {
  secrets := hush.Hushfile()
  db := getDBConnection(secrets)
  csvFile := csv.NewReader(getCSVFile(os.Args[2]))
  records, err := csvFile.ReadAll()

  if err != nil {
    panic(err)
  }

  table, ok := secrets.GetString("table")
  if !ok {
    panic("Could not get the table name.")
  }

  sql := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3, $4);", table)
  stmt, err := db.Prepare(sql)
  if err != nil {
    panic(err)
  }

  for _, record := range records {
    cardType, err := strconv.ParseInt(record[2], 0, 0)
    if err != nil {
      panic(err)
    }

    cardBlanks, err := strconv.ParseInt(record[3], 0, 0)
    if err != nil {
      panic(err)
    }

    res, err := stmt.Exec(record[1], cardType, cardBlanks, true)
    if err != nil {
      panic(err)
    }
    fmt.Println(res)
  }
}
