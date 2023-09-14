package main
import(
   "gorm.io/gorm"
   "net/http"
   "encoding/json"
)

type RecipeHandlers struct {
   db *gorm.DB
}

type RecipeListResponse struct {
   Recipes []RecipeListing `json:"recipes"`
}

type RecipeListing struct {
   ID uint `json:"id"`
   Name string `json:"name"`
}

func RecipeListResponseFromOrm(recipes []Recipe) RecipeListResponse {
   response := RecipeListResponse{
      Recipes: make([]RecipeListing, len(recipes)),
   }
   for i := range recipes {
      response.Recipes[i].ID = recipes[i].ID
      response.Recipes[i].Name = recipes[i].Name
   }
   return response
}

func (rh RecipeHandlers) List(w http.ResponseWriter, r *http.Request){
   page := get_page_from_request(r)
   var recipes []Recipe
   rh.db.Where("id >= ?", page).Limit(10).Find(&recipes)
   response, err := json.Marshal(RecipeListResponseFromOrm(recipes))
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
   }
   w.Write(response)
}

func (rh RecipeHandlers) Create(w http.ResponseWriter, r *http.Request){}
func (rh RecipeHandlers) Detail(w http.ResponseWriter, r *http.Request){}
func (rh RecipeHandlers) Delete(w http.ResponseWriter, r *http.Request){}
func (rh RecipeHandlers) UpdateIngredients(w http.ResponseWriter, r *http.Request){}
