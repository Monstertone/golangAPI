package main

import (

  "log"
  "os"
  "database/sql"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/lib/pq"
  "github.com/subosito/gotenv"



)

type Transaction struct {
  ID             int     `json:id`
  Amount         int     `json:amount`
  Business_name  string  `json:business_name`
  Type_trans     string  `json:type`
  Users_ID       int     `json:users_id`
}

var transactions []Transaction
var db *sql.DB

func init(){
  gotenv.Load()
}

func logFatal(err error){
  if err != nil {
    log.Fatal(err)
  }
}

func main() {

  pgUrl, err := pq.ParseURL(os.Getenv("POSTGRESQL_URL"))
  logFatal(err)

  db, err = sql.Open("postgres", pgUrl)
  logFatal(err)

  err = db.Ping()
  logFatal(err)

  log.Println(pgUrl)

  router := mux.NewRouter()

  router.HandleFunc("/transactions", getTransactions).Methods("GET")
  router.HandleFunc("/transactions/{id}", getTransaction).Methods("GET")
  router.HandleFunc("/transactions", addTransaction).Methods("POST")
  router.HandleFunc("/transactions", updateTransaction).Methods("PUT")
  router.HandleFunc("/transactions/{id}", deleteTransaction).Methods("DELETE")

  log.Fatal(http.ListenAndServe(":8000", router))
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction
  transactions = []Transaction{}

  rows, err := db.Query("select * from transactions")
  logFatal(err)

  defer rows.Close()

  for rows.Next(){
    err :=rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Business_name,&transaction.Type_trans, &transaction.Users_ID)
    logFatal(err)

    transactions = append(transactions, transaction)
  }
   json.NewEncoder(w).Encode(transactions)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction
  params := mux.Vars(r)

  rows := db.QueryRow("select * from transactions where id=$1", params["id"])

  err := rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Business_name, &transaction.Type_trans, &transaction.Users_ID)
  logFatal(err)

  json.NewEncoder(w).Encode(transaction)


}

func addTransaction(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction
  var transactionID int

  json.NewDecoder(r.Body).Decode(&transaction)

  err := db.QueryRow("insert into transactions (id, amount, business_name, type_trans, users_id) values ($1, $2, $3, $4, $5) RETURNING id;",
              transaction.ID, transaction.Amount, transaction.Business_name, transaction.Type_trans, transaction.Users_ID).Scan(&transactionID)

  logFatal(err)

  json.NewEncoder(w).Encode(transactionID)

}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction
  json.NewDecoder(r.Body).Decode(&transaction)

  result, err := db.Exec("update transactions set amount=$1, business_name=$2, type_trans=$3, users_id=$4 where id=$5 RETURNING id",
  &transaction.Amount, &transaction.Business_name, &transaction.Type_trans, &transaction.Users_ID, &transaction.ID)

rowsUpdated, err := result.RowsAffected()
logFatal(err)

json.NewEncoder(w).Encode(rowsUpdated)

}



func deleteTransaction(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  log.Println(params["id"])

  result, err := db.Exec("delete from transactions where id=$1", params["id"])
  logFatal(err)

  rowsDeleted, err := result.RowsAffected()
  logFatal(err)

  json.NewEncoder(w).Encode(rowsDeleted)

}
