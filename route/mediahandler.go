package route

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(mux.Vars(r)["id"], ".")[0]
	f, err := os.Open(fmt.Sprintf("_/media/%s.mp4", id))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=%q", id+".mp4"))
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, id, time.Time{}, f)
}
