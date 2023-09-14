
package crud

import (
   "gorm.io/gorm"
   "net/http"
   "log"
   "encoding/json"
   "io"
)

type CrudController [T interface{}] struct {
   DB *gorm.DB
}

func (cc *CrudController[T]) List(w http.ResponseWriter, r *http.Request){
   models := make([]T,10)
   result := cc.DB.Find(&models)
   if result.Error != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when fetching from database %s\n", result.Error)
      return
   }
   log.Printf("Fetched %v rows\n", result.RowsAffected)
   response,err := json.Marshal(models)
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when serializing models %s\n", err)
      return
   }
   w.Write(response)
}

func (cc *CrudController[T]) Post(w http.ResponseWriter, r *http.Request){
   var model T
   request_data, err := io.ReadAll(r.Body)
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when reading data from request: %s\n", err)
      return
   }
   json.Unmarshal(request_data, &model)
   result := cc.DB.Create(&model)
   if result.Error != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when inserting data: %s\n", result.Error)
      return
   }

   w.WriteHeader(http.StatusCreated)
   data,err := json.Marshal(model)
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Something went wrong when marshalling inserted data: %s\n", err)
      return
   }
   w.Write(data)
}
