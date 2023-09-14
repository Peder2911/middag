
package main

import (
   "net/http"
   "github.com/peder2911/kitchen/models"
   "github.com/peder2911/kitchen/crud"
   "github.com/gorilla/mux"
   "gorm.io/gorm"
   "gorm.io/driver/sqlite"
)


func main(){
   db,err := gorm.Open(sqlite.Open("./foo.db"))
   if err != nil {
      panic(err)
   }
   db.AutoMigrate(models.Ingredient{}, models.MeasuringUnit{}, models.RecipeIngredient{}, models.Recipe{})
   
   router := mux.NewRouter()
   ingredient_handler := crud.CrudController[models.Ingredient]{DB: db}
   router.HandleFunc("/ingredient", ingredient_handler.List).Methods("GET")
   router.HandleFunc("/ingredient", ingredient_handler.Post).Methods("POST")
   http.ListenAndServe(":8000",router)
}
