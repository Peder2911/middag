
package crud 

// Database models

type Recipe struct {
   Id int      `json:"id"`
   Name string `json:"name"`
}

type RecipeIngredient struct {
   ingredient int
   amount float32
   measuring_unit string 
}

type Ingredient struct {
   Id int64    `json:"id" db:"id"`
   Name string `json:"name" db:"name"`
}
