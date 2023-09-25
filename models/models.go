
package models
import "gorm.io/gorm"

type Ingredient struct {
   gorm.Model
   Name string
}
func (i Ingredient) GetPrimaryKey() uint {
   return i.ID
}
func (i Ingredient) SetPrimaryKey(id uint) {
   i.ID = id
}

type MeasuringUnit struct {
   gorm.Model
   Name string `gorm:"unique"`
   RecipeIngredient []RecipeIngredient
}

type RecipeIngredient struct {
   gorm.Model
   RecipeId uint `gorm:"primary_key"` 
   IngredientId uint `gorm:"primary_key"`
   Amount uint 
   MeasuringUnitId uint
}

type Recipe struct {
   gorm.Model
   Name string `gorm:"unique"`
   Ingredients []Ingredient `gorm:"many2many:recipe_ingredients"`
}
