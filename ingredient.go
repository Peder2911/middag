package main
import (
   "gorm.io/gorm"
   "net/http"
   "encoding/json"
   "io"
   "log"
   "github.com/gorilla/mux"
   "strconv"
)

type IngredientHandlers struct {
   db *gorm.DB
}

func (rh IngredientHandlers) List(w http.ResponseWriter, r *http.Request){
   page := get_page_from_request(r)
   var ingredients []Ingredient
   rh.db.Where("id >= ?", page).Limit(10).Find(&ingredients)
   response, err := json.Marshal(IngredientListResponseFromOrm(ingredients))
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
   }
   w.Write(response)
}
type IngredientListResponse struct {
   Ingredients []IngredientListing `json:"ingredients"`
}
type IngredientListing struct {
   ID uint `json:"id"`
   Name string `json:"name"`
}
func IngredientListResponseFromOrm(ingredients []Ingredient) IngredientListResponse{
   response := IngredientListResponse {
      Ingredients: make([]IngredientListing, len(ingredients)),
   }
   for i := range ingredients {
      response.Ingredients[i].ID = ingredients[i].ID
      response.Ingredients[i].Name = ingredients[i].Name

   }
   return response
}

func (rh IngredientHandlers) Create(w http.ResponseWriter, r *http.Request){
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
   ingredient := Ingredient{Name: new_ingredient.Name}
   result := rh.db.Create(&ingredient)
   if result.Error != nil {
      log.Printf("Failed to create new ingredient: %s\n", result.Error)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
   response_data, err := json.Marshal(ingredient)
   w.Write(response_data)
}

type NewIngredient struct {
   Name string `json:"name"`
}

func (rh IngredientHandlers) Delete(w http.ResponseWriter, r *http.Request){
   vars := mux.Vars(r)
   ingredient_id, err := strconv.Atoi(vars["ingredient_id"])
   if err != nil {
      log.Printf("Failed to parse ingredient id as int: \"%s\"\n", vars["ingredient_id"])
   }
   result := rh.db.Delete(&Ingredient{}, ingredient_id)
   if result.Error != nil {
      log.Printf("Failed to delete ingredient: %s\n", result.Error)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
}
