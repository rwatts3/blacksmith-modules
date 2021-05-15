package sqlikedestination

/*
Statement holds a prepared statement including its query and values to load into
a SQL-like database.

Example:

  sqlikedestination.Statement{
    Query: "INSERT INTO users ("first_name", "last_name", "email") VALUES ($1, $2, $3);",
    Values: [][]interface{}{
      {"John", "Doe", "johndoe@example.com"},
      {"Jane", "Doe", "janedoe@example.com"},
    },
  }
*/
type Statement struct {
	Query  string          `json:"query"`
	Values [][]interface{} `json:"values"`
}
