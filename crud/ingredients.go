package crud

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

type IngredientHandler struct {
   Database *sql.DB
}

func (ih IngredientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
   if r.Method == "GET" {
      if r.URL.Path == "/api/ingredients/" {
         ih.list_ingredients(w, r)
      } else {
         if m,_ := path.Match("/api/ingredients/*", r.URL.Path); m {
            ih.show_ingredient(w, r)
         } else {
            w.WriteHeader(http.StatusNotFound)
            return
         }
      }
   } else if r.Method == "POST" {
      if r.URL.Path == "/api/ingredients/" {
         ih.create_ingredient(w, r)
      } else {
         w.WriteHeader(http.StatusNotFound)
      }
//   } else if r.Method == "PUT" {
//      if m,_ := path.Match("/api/ingredients/*", r.URL.Path); m {
//         ih.update_ingredient(w, r)
//      } else {
//         w.WriteHeader(http.StatusNotFound)
//         return
//      }
   } else if r.Method == "DELETE" {
      if m,_ := path.Match("/api/ingredients/*", r.URL.Path); m {
         ih.delete_ingredient(w,r)
      } else {
         w.WriteHeader(http.StatusNotFound)
      }
   } else {
      w.WriteHeader(http.StatusMethodNotAllowed)
   }
}

func (ih IngredientHandler) list_ingredients(w http.ResponseWriter, r *http.Request){
   type response_model struct {
      Ingredients []Ingredient `json:"ingredients"`
   }

   page := get_page_from_request(r)

   var ingredients []Ingredient = make([]Ingredient, 0)
   result, err := ih.Database.Query(`
      select name,id from ingredient where id > $1 order by id asc limit $2
   `, page*PAGESIZE, PAGESIZE)
   if err != nil {
      log.Printf("Error when listing ingredients: %s\n", err)
   } else {
      var ingredient Ingredient 
      for result.Next() {
         ingredient = Ingredient{}
         err := result.Scan(&ingredient.Name, &ingredient.Id)
         if err != nil {
            log.Printf("Failed to scan ingredient row: %s\n", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
         }
         ingredients = append(ingredients, ingredient)
      }
   }

   response := response_model{Ingredients:ingredients}
   response_data,err := json.Marshal(response)
   if err != nil{
      log.Printf("Error when serializing ingredient response: %s", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
   w.Write(response_data)
}

func (ih IngredientHandler) show_ingredient(w http.ResponseWriter, r *http.Request){
   _,ingredient_id_string := path.Split(r.URL.Path)
   ingredient_id, err := strconv.Atoi(ingredient_id_string)
   if err != nil {
      w.WriteHeader(http.StatusNotFound)
      return 
   }
   var response_model Ingredient
   row := ih.Database.QueryRow(`
      select id,name from ingredient where id = $1
   `, ingredient_id)
   err = row.Scan(&response_model.Id, &response_model.Name)
   if err != nil {
      if err == sql.ErrNoRows {
         w.WriteHeader(http.StatusNotFound)
         return
      }
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Error when fetching ingredient: %s\n", err)
      return
   }
   response_data, err := json.Marshal(response_model)
   if err != nil {
      log.Printf("Error when serializing ingredient: %s", err)
   }
   w.Write(response_data)
}

func (ih IngredientHandler) create_ingredient(w http.ResponseWriter, r *http.Request){
   type NewIngredient struct {
      Name string `json:"name"`
   }

   var err error
   request_data, err := io.ReadAll(r.Body)
   if err != nil {
      log.Printf("Failed to read from request: %s\n", err)
      w.WriteHeader(http.StatusBadRequest)
      return
   }
   var new_ingredient NewIngredient
   err = json.Unmarshal(request_data, &new_ingredient)
   if err != nil {
      log.Printf("Failed to deserialize JSON from request: %s\n", err)
      w.WriteHeader(http.StatusBadRequest)
      return
   }
   result, err := ih.Database.Query(`
      insert into ingredient(name) values ($1) returning id
   `, new_ingredient.Name)
   if err != nil {
      log.Printf("Failed to insert ingredient into database: %s\n", err)
      w.WriteHeader(http.StatusBadRequest)
      return
   }
   response_model := Ingredient{Name: new_ingredient.Name}
   result.Next()
   err = result.Scan(&response_model.Id)
   if err != nil {
      log.Printf("Failed to scan newly inserted row from database: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
   }
   response_data, err := json.Marshal(response_model)
   if err != nil{
      log.Printf("Failed to marshal JSON of inserted ingredient: %s", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   w.Write(response_data)
}

func (ih IngredientHandler) delete_ingredient(w http.ResponseWriter, r *http.Request) {
   _,ingredient_id_string := path.Split(r.URL.Path)
   ingredient_id, err := strconv.Atoi(ingredient_id_string)
   if err != nil {
      w.WriteHeader(http.StatusNotFound)
      return
   }
   result := ih.Database.QueryRow(`
      delete from ingredient where id = $1 returning id
   `, ingredient_id)
   err = result.Scan(&ingredient_id)
   if err != nil {
      if err == sql.ErrNoRows {
         w.WriteHeader(http.StatusNotFound)
         return
      }
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Error when trying to delete ingredient %v: %s", ingredient_id, err)
      return
   }
   w.WriteHeader(http.StatusNoContent)
}
