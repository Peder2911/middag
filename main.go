package main

import (
   "database/sql"
   _ "embed"
   "fmt"
   _ "github.com/lib/pq"
   "github.com/peder2911/middag/crud"
   "log"
   "net/http"
)

//go:embed dist/index.html
var index_html string

func main() {
   cfg := ReadConfig()
   connstr := fmt.Sprintf(
      "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
      cfg.database.username,
      cfg.database.password,
      cfg.database.host,
      cfg.database.port,
      cfg.database.name,
   )
   log.Println(connstr)
   database, err := sql.Open(
      "postgres",
      connstr,
   )
   if err != nil {
      panic(fmt.Sprintf("Failed to connect to database: %s", err))
   }

   mux := http.NewServeMux()

   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte(index_html))
   })

   mux.Handle("/api/recipes/", crud.RecipeHandler{Database: database})
   mux.Handle("/api/ingredients/", crud.IngredientHandler{Database: database})
   log.Fatal(http.ListenAndServe(":8000", mux))
}
