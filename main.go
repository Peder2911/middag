package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Ingredient struct {
   gorm.Model
   ID   uint    `json:"id"`
   Name string `gorm:"unique" json:"name"`
}

type Recipe struct {
   gorm.Model
   ID          uint          `json:"id"`
   Name        string       `gorm:"unique" json:"name"`
   Ingredients []Ingredient `gorm:"many2many:recipe_ingredients" json:"ingredients"`
}

type MeasuringUnit struct {
   gorm.Model
   ID uint
   Name string `gorm:"unique"`
   RecipeIngredient []RecipeIngredient
}

type RecipeIngredient struct {
   gorm.Model
   RecipeId      uint `gorm:"primary_key"`
   IngredientId  uint `gorm:"primary_key"`
   Amount        uint
   MeasuringUnitID uint
}

//go:embed dist/index.html
var index_html string

func get_page_from_request(r *http.Request) int {
   page,err := strconv.Atoi(r.URL.Query().Get("page"))
   if err != nil {page = 0}
   return page
}

func main() {
   config := ReadConfig()
    
   db, err := gorm.Open(postgres.Open(fmt.Sprintf(
         "user=%s password=%s host=%s port=%s dbname=%s",
         config.Database.Username, 
         config.Database.Password, 
         config.Database.Host, 
         config.Database.Port, 
         config.Database.Name)), &gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info})})
   if err != nil {
      panic(err)
   }

   db.AutoMigrate(&Recipe{}, &Ingredient{}, &MeasuringUnit{}, &RecipeIngredient{})
   err = db.SetupJoinTable(&Recipe{}, "Ingredients", &RecipeIngredient{})
   if err != nil {
      panic(err)
   }

   log.Println("Initialization complete.")

   r := mux.NewRouter() 
   r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
      w.Write([]byte(index_html))
   })

   s := r.PathPrefix("/api").Subrouter()

   //recipe_handlers := RecipeHandlers{db}
   //s.HandleFunc("/recipe", recipe_handlers.List).Methods("GET")
   //s.HandleFunc("/recipe", recipe_handlers.Create).Methods("POST")
   //s.HandleFunc("/recipe/{recipe_id:[0-9]+}", recipe_handlers.Detail).Methods("GET")
   //s.HandleFunc("/recipe/{recipe_id:[0-9]+}", recipe_handlers.Delete).Methods("DELETE")
   //s.HandleFunc("/recipe/{recipe_id:[0-9]+}/ingredients", recipe_handlers.UpdateIngredients).Methods("PUT")

   ingredient_handlers := CrudController[Ingredient]{db}
   s.HandleFunc("/ingredient", ingredient_handlers.List).Methods("GET")
   //s.HandleFunc("/ingredient", ingredient_handlers.Create).Methods("POST")
   //s.HandleFunc("/ingredient/{ingredient_id:[0-9]+}", ingredient_handlers.Delete).Methods("DELETE")
   http.ListenAndServe(":8000", r)

   //measuring_unit_handlers := MeasuringUnitHandlers{db}
   //s.HandleFunc("/ingredient", ingredient_handlers.List).Methods("GET")
   //s.HandleFunc("/ingredient", ingredient_handlers.Create).Methods("POST")
   //s.HandleFunc("/ingredient/{ingredient_id:[0-9]+}", ingredient_handlers.Delete).Methods("DELETE")
   //http.ListenAndServe(":8000", r)
}
