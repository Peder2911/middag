package main
import (
   "gorm.io/gorm"
   "net/http"
   "log"
   "encoding/json"
)

type GormModel interface {
   PrimaryKey() uint 
}

type CrudController[T interface{}] struct {
   db *gorm.DB
}

func (cc CrudController[T]) List (w http.ResponseWriter, r *http.Request) {
   page := get_page_from_request(r)
   models := make([]T,10) 
   result := cc.db.Find(&models).Where("id >= ?", page).Limit(10)
   if result.Error != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when fetching from database: %s\n", result.Error)
      return
   }
   log.Printf("Fetched %v rows\n", result.RowsAffected)

   response, err := json.Marshal(models)
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when serializing data to JSON: %s\n", err)
      return
   }

   w.Write(response)
}
