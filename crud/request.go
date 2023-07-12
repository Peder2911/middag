package crud 

import (
   "net/http"
   "strconv"
)

func get_page_from_request(r *http.Request) int {
   page,err := strconv.Atoi(r.URL.Query().Get("page"))
   if err != nil {
      page = 0
   }
   return page
}
