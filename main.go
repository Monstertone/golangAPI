package main

import (
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "strconv"

  _ "database/sql"
  _ "github.com/lib/pq"
  _ "github.com/subosito/gotenv"
)

type Transaction struct {
  ID             int     `json:id`
  Users_ID       int     `json:users_ID`
  Amount         string  `json:amount`
  Type           string  `json:type`
  Business_name  string  `json:business_name`
}

var transactions []Transaction

func main() {
  router := mux.NewRouter()

  transactions = append(transactions,
Transaction{ID: 1, Users_ID: 1, Amount: "44.97", Type:"gasoline", Business_name:"Costco"},
Transaction{ID: 2, Users_ID: 2, Amount: "16.99", Type:"paper supplies", Business_name:"Office Depot"},
Transaction{ID: 3, Users_ID: 3, Amount: "108.43", Type:"computer parts", Business_name:"Frys Electronics"},
Transaction{ID: 4, Users_ID: 1, Amount: "79.99", Type:"printer", Business_name:"Frys Electronics"},
Transaction{ID: 5, Users_ID: 3, Amount: "34.27", Type:"learning courses", Business_name:"Udemy"},
Transaction{ID: 6, Users_ID: 3, Amount: "299.99", Type:"galaga machine", Business_name:"Walmart"})


  router.HandleFunc("/transactions", getTransactions).Methods("GET")
  router.HandleFunc("/transactions/{id}", getTransaction).Methods("GET")
  router.HandleFunc("/transactions", addTransaction).Methods("POST")
  router.HandleFunc("/transactions", updateTransaction).Methods("PUT")
  router.HandleFunc("/transactions/{id}", deleteTransaction).Methods("DELETE")

  log.Fatal(http.ListenAndServe(":8000", router))
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(transactions)
}
func getTransaction(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  i, _ := strconv.Atoi(params["id"])

  for _, transaction := range transactions {
    if transaction.ID == i {
      json.NewEncoder(w).Encode(&transaction)
    }
  }

}

func addTransaction(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction
  json.NewDecoder(r.Body).Decode(&transaction)

  transactions = append(transactions, transaction)

  json.NewEncoder(w).Encode(transactions)
  log.Println(transaction)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
  var transaction Transaction

  json.NewDecoder(r.Body).Decode(&transaction)

        for i, item := range transactions {
        if item.ID == transaction.ID {
          transactions[i] = transaction

        }
    }

    json.NewEncoder(w).Encode(transactions)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  id, _ := strconv.Atoi(params["id"])

  for i, item := range transactions {
    if item.ID == id {
      transactions = append(transactions[:i], transactions[i+1:]...)
    }
  }

  json.NewEncoder(w).Encode(transactions)

}
