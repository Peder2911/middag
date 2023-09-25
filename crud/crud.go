
package crud

import (
   "gorm.io/gorm"
   "github.com/gorilla/mux"
   "net/http"
   "log"
   "encoding/json"
   "io"
   "strconv"
)

type Entity interface {
   GetPrimaryKey() uint
   SetPrimaryKey(id uint)
}

type CrudController [T Entity] struct {
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

func (cc *CrudController[T]) Delete (w http.ResponseWriter, r *http.Request){
   idstring,ok := mux.Vars(r)["id"]
   if ! ok {
      w.WriteHeader(http.StatusInternalServerError)
      log.Printf("Delete controller method called without ID path parameter.")
      return
   }
   id,err := strconv.Atoi(idstring)
   if err != nil {
      log.Printf("Error when converting path parameter: %s\n", err)
      w.WriteHeader(http.StatusNotFound)
      return
   }


   var model T
   result := cc.DB.Where("id = ?", id).Delete(&model)
   if result.Error != nil {
      if result.Error == gorm.ErrRecordNotFound {
         w.WriteHeader(http.StatusNotFound)
         return
      } else {
         w.WriteHeader(http.StatusInternalServerError)
         log.Printf("Error when deleting record: %s\n", result.Error)
         return
      }
   }
   w.WriteHeader(http.StatusNoContent)
   return
}

func (cc *CrudController[T]) Patcher(fieldname string, validate func(v string)(any, error)) http.HandlerFunc{
   return func(w http.ResponseWriter, r *http.Request){
      idstring,ok := mux.Vars(r)["id"]
      if !ok {
         w.WriteHeader(http.StatusNotFound)
         return
      }
      id,err := strconv.Atoi(idstring)
      if err != nil {
         w.WriteHeader(http.StatusNotFound)
         return
      }
      var model T
      request_data, err := io.ReadAll(r.Body)
      if err != nil {
         w.WriteHeader(http.StatusInternalServerError)
         return
      }
      new_value,err := validate(string(request_data))
      if err != nil {
         w.WriteHeader(http.StatusBadRequest)
         return
      }
      cc.DB.Find(&model, id).Update(fieldname, new_value)
      w.WriteHeader(http.StatusAccepted)
   }
}

func IdentityValidate(v string) (any, error) {
   return v, nil
}
