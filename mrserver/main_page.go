package mrserver

import "net/http"

func MainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"STATUS\": \"OK\"}"))
}
