package models

type Transaction struct {
  ID             int     `json:id`
  Amount         int     `json:amount`
  Business_name  string  `json:business_name`
  Type_trans     string  `json:type`
  Users_ID       int     `json:users_id`
}
