package api

import (
	"log"
	"net/http"
)

func HandleIndex(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	if req.Method != http.MethodGet {
		http.Error(res, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
}
