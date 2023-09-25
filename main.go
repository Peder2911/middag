
package main

import (
   "net/http"
   "github.com/peder2911/middag/models"
   "github.com/peder2911/middag/crud"
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
   ingredient_handler := crud.NewIngredientController(db)
   router.HandleFunc("/ingredient", ingredient_handler.List).Methods("GET")
   router.HandleFunc("/ingredient", ingredient_handler.Post).Methods("POST")
   router.HandleFunc("/ingredient/{id}", ingredient_handler.Detail).Methods("GET")
   router.HandleFunc("/ingredient/{id}", ingredient_handler.Delete).Methods("DELETE")
   router.HandleFunc("/ingredient/{id}/name", ingredient_handler.Patcher("name",crud.IdentityValidate)).Methods("POST")
   http.ListenAndServe(":8000",router)
}
