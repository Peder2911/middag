package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
)

type RecipesListing struct {
   Recipes []Recipe `json:"recipes"`
}

type RecipeDetail struct {
   Id          int `json:"id"`
   Name        string `json:"name"`
   Ingredients []RecipeIngredientListing `json:"ingredients"`
}

type RecipeIngredientListing struct {
   Name          string `json:"name"`
   Amount        int `json:"amount"`
   MeasuringUnit string `json:"measuring_unit"`
}

type RecipePost struct {
   Name        string `json:"name"`
   Ingredients []RecipeIngredientListing `json:"ingredients"`
}

type RecipeHandler struct {
   Database *sql.DB
}

func (rh RecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
   if r.Method == "GET" {
      if r.URL.Path == "/api/recipes/" {
         rh.list_recipes(w, r)
      } else {
         matches, err := path.Match("/api/recipes/*",r.URL.Path)
         if err != nil {
            panic(err)
         }
         if ! matches {
            log.Printf("Unknown path: %s\n", r.URL.Path)
            w.WriteHeader(http.StatusNotFound)
            return
         }
         rh.show_recipe(w, r)
      }
   } else if r.Method == "POST" {
      rh.create_new_recipe(w, r)
   //} else if r.Method == "PUT" {
   //   rh.update_recipe(w, r)
   } else {
      w.WriteHeader(http.StatusMethodNotAllowed)
   }
}

type RecipeIngredientHandler struct {
   Database *sql.DB
   PathRegexp *regexp.Regexp
}

func NewRecipeIngredientHandler(database *sql.DB) RecipeIngredientHandler {
   return RecipeIngredientHandler {
      Database: database,
      PathRegexp: regexp.MustCompile("/api/recipes/(?P<recipe_id>[0-9]+)/ingredients/"),
   }
}

func (rih RecipeIngredientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

   matches, err := path.Match("/api/recipes/*/ingredients/", r.URL.Path)
   recipe_id := rih.PathRegexp.FindAllStringSubmatch()[
   if err != nil {
      panic(err)
   }
   if ! matches {
      log.Printf("Unknown path: %s\n", r.URL.Path)
      w.WriteHeader(http.StatusNotFound)
      return
   }

   if r.Method == "GET" {
      rih.get_recipe_ingredients(w, r)
   } else if r.Method == "PUT" {
      rih.set_recipe_ingredients(w, r)
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
   _, recipe_id := path.Split(r.URL.Path)
   result, err := rh.Database.Query("select recipe.id, recipe.name from recipe where id = $1", recipe_id)
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

   var recipe RecipeDetail
   result.Scan(&recipe.Id, &recipe.Name)

   items, err := rh.Database.Query(`
      select ingredient.name, recipe_ingredient.amount, recipe_ingredient.measuring_unit 
      from ingredient      
      join recipe_ingredient on ingredient.id = recipe_ingredient.ingredient
      where recipe_ingredient.recipe=$1 
      `, recipe.Id)

   for items.Next() {
      item := RecipeIngredientListing{}
      items.Scan(&item.Name, &item.Amount, &item.MeasuringUnit)
      recipe.Ingredients = append(recipe.Ingredients, item)
   }

   data,err := json.Marshal(recipe)
   if err != nil {
      log.Printf("Something went wrong when serializing a recipe: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
   w.Write(data)
}

func (rh RecipeHandler) create_new_recipe(w http.ResponseWriter, r *http.Request) {
   var err error
   var recipe_post RecipePost

   request_data, err := io.ReadAll(r.Body)
   if err != nil {
      log.Printf("Failed to read data: %s\n", err)
      w.WriteHeader(http.StatusBadRequest)
      return
   }

   err = json.Unmarshal(request_data, &recipe_post)
   if err != nil {
      log.Printf("Received bad data: %s\n", err)
      w.WriteHeader(http.StatusBadRequest)
      return
   }

   tx, err := rh.Database.Begin()
   if err == nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   result,err := tx.Query("insert into recipe (name) values ($1) returning id", recipe_post.Name)
   if err != nil {
      log.Printf("Error when inserting recipe: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      tx.Rollback()
      return
   }
   response := Recipe{Name: recipe_post.Name}
   result.Next()
   err = result.Scan(&response.Id)
   if err != nil {
      log.Printf("Error when reading ID of inserted recipe: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   for _, ingredient := range(recipe_post.Ingredients) {
      inserted := 0
      result, err := tx.Query(`
         with rows as (
            insert into recipe_ingredient (recipe, ingredient, amount, measuring_unit) 
            select recipe.id, ingredient.id, $1, $2 from ingredient 
               cross join (select id from recipe where recipe.name=$3) as recipe
               where ingredient.name=$4
               returning recipe, ingredient)
         select count(*) from rows
         `,
         ingredient.Amount, ingredient.MeasuringUnit, recipe_post.Name, ingredient.Name)
      if err != nil {
         log.Printf("Error when running query to insert recipe ingredient: %s\n", err)
         w.WriteHeader(http.StatusInternalServerError)
         tx.Rollback()
         return
      }
      if ! result.Next() {
         log.Printf("Insert recipe ingredient query returned no rows!")
         w.WriteHeader(http.StatusInternalServerError)
         tx.Rollback()
         return
      }

      err = result.Scan(&inserted)
      if err != nil {
         log.Printf("Failed to scan number of inserted rows: %s", err)
         w.WriteHeader(http.StatusInternalServerError)
         tx.Rollback()
         return
      }

      if inserted == 0 {
         log.Printf("Ingredient not found: %s", ingredient.Name)
         w.WriteHeader(http.StatusBadRequest)
         tx.Rollback()
         return
      } else if inserted > 1 {
         log.Printf("Something went wrong when inserting recipe ingredient relation: Tried to insert %v assoc. rows!", inserted) 
         w.WriteHeader(http.StatusInternalServerError)
         tx.Rollback()
         return
      }
   }

   err = tx.Commit()
   if err != nil {
      log.Printf("Error when committing recipe data: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   data,err := json.Marshal(response)
   if err != nil {
      log.Printf("Error when marshalling recipe data: %s\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type","application/json")
   w.Write(data)
   w.WriteHeader(http.StatusOK)
   return
}

//func (rh RecipeHandler) update_recipe(w http.ResponseWriter, r *http.Request) {
//}

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
