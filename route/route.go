package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"io.6f.sana/template"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", template.HttpHandler).Methods("GET")
	r.HandleFunc("/recent", template.CustomHandler("", nil)).Methods("GET")
	r.HandleFunc("/random", template.CustomHandler("", nil)).Methods("GET")
	r.HandleFunc("/media/{id}", template.CustomHandler("", nil)).Methods("GET")
	r.HandleFunc("/resolve/{id}", resolveHandler).Methods("GET")
	r.PathPrefix("/s").Handler(http.StripPrefix("/s", http.FileServer(http.Dir("_/public"))))
	r.HandleFunc("/_/search", searchHandler).Methods("GET")
	r.HandleFunc("/_/content/recent", recentHandler).Methods("GET")
	r.HandleFunc("/_/content/random/{type}", randomHandler).Methods("GET")
	r.HandleFunc("/_/content/{id}/thumbnail-{size}", thumbnailHandler).Methods("GET")
	r.HandleFunc("/_/content/{id}/media", mediaHandler).Methods("GET")
	r.HandleFunc("/_/content/{id}/media{_}", mediaHandler).Methods("GET")
	r.HandleFunc("/_/content/{id}/info", infoHandler).Methods("GET")
	r.HandleFunc("/_/content/{id}/info{_}", infoHandler).Methods("GET")
	return r
}
