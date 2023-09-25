package crud

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/peder2911/middag/models"
	"gorm.io/gorm"
        "strconv"
        "log"
        "encoding/json"
)

type IngredientCrudController struct {
   CrudController[models.Ingredient]
}

func NewIngredientController(db *gorm.DB) IngredientCrudController {
   controller := IngredientCrudController{}
   controller.DB = db
   return controller
}

func (icc *IngredientCrudController) Detail(w http.ResponseWriter, r *http.Request){
   idstring,ok := mux.Vars(r)["id"]
   if ! ok {
      w.WriteHeader(http.StatusNotFound)
      return
   }
   id,err := strconv.Atoi(idstring)
   if err != nil {
      w.WriteHeader(http.StatusNotFound)
      return
   }
   var ingredient models.Ingredient
   result := icc.DB.First(&ingredient, id)
   if result.Error != nil {
      if result.Error == gorm.ErrRecordNotFound {
         w.WriteHeader(http.StatusNotFound)
         return
      } else {
         w.WriteHeader(http.StatusInternalServerError)
         log.Printf("Error when fetching ingredient %v: %s\n", id, result.Error)
      }
   }
   data,err := json.Marshal(ingredient)
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Error when serializing ingredient: %s\n", err)
   }
   w.Write(data)
}
