package crud

import (
   "database/sql"
   "encoding/json"
   "fmt"
   "log"
   "path"
   "net/http"
   "strconv"
)

type RecipesListing struct {
   Recipes []Recipe `json:"recipes"`
}

type RecipeHandler struct {
   Database *sql.DB
}

func (rh RecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
   if r.Method == "GET" {
      if r.URL.Path == "/api/recipes/" {
         rh.list_recipes(w, r)
      } else {
         rh.show_recipe(w, r)
      }
   } else if r.Method == "POST" {
      rh.create_new_recipe(w, r)
   } else if r.Method == "PUT" {
      rh.update_recipe(w, r)
   } else {
      w.WriteHeader(http.StatusMethodNotAllowed)
   }
}

func (rh RecipeHandler) list_recipes(w http.ResponseWriter, r *http.Request) {
   var page int
   var err error
   var response_data []byte 
   var response_model RecipesListing = RecipesListing{}

   page, err = strconv.Atoi(r.URL.Query().Get("page"))
   if err != nil {
      page = 0
   }

   response_model.Recipes = rh.fetch_recipes_list(page)
   response_data, err = json.Marshal(response_model)
   if err != nil {
      w.WriteHeader(500)
   } else {
      w.Write(response_data)
   }
}

func (rh RecipeHandler) show_recipe(w http.ResponseWriter, r *http.Request) {
   matches, err := path.Match("/api/recipes/*",r.URL.Path)
   if err != nil {
      panic(err)
   }
   if ! matches {
      log.Printf("Unknown path: %s\n", r.URL.Path)
      w.WriteHeader(http.StatusNotFound)
      return
   }

   _, recipe_id := path.Split(r.URL.Path)
   result, err := rh.Database.Query("select * from recipe where id = $1", recipe_id)
   if err != nil {
      log.Printf("Error when fetching a recipe: %s\n", err)
      w.WriteHeader(http.StatusNotFound)
      return
   }
   exists := result.Next()
   if ! exists {
      log.Printf("Recipe %s was not found.\n", recipe_id)
      w.WriteHeader(http.StatusNotFound)
      return
   }

   var recipe Recipe
   result.Scan(&recipe)

   data,err := json.Marshal(recipe)
   if err != nil {
      log.Printf("Something went wrong when serializing a recipe: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
   }
   w.Write(data)
}

func (rh RecipeHandler) create_new_recipe(w http.ResponseWriter, r *http.Request) {
}

func (rh RecipeHandler) update_recipe(w http.ResponseWriter, r *http.Request) {
}

func (rh RecipeHandler) fetch_recipes_list(page int) []Recipe {
   var err error
   var recipes []Recipe = make([]Recipe, 0)

   result, err := rh.Database.Query(`
      select * from recipe where id > $1 order by id asc limit $2;
   `, page, PAGESIZE)

   if err != nil {
      log.Println(fmt.Sprintf("Error when fetching recipes from database: %s", err))
      return recipes
   }

   for result.Next(){
      var recipe Recipe
      result.Scan(&recipe)
      recipes = append(recipes, recipe)
   }
   return recipes
}
