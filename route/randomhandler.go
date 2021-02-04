package route

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"io.6f.sana/db"
)

func randomHandler(res http.ResponseWriter, req *http.Request) {
	result := db.MediaCollection.Find(nil)

	count, err := result.Count()
	if err != nil {
		res.WriteHeader(500)
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())
	var m db.Media
	if err := result.Skip(uint(rand.Intn(int(count)))).One(&m); err != nil {
		res.WriteHeader(500)
		return
	}

	http.Redirect(res, req, filepath.Join("/_/content/", fmt.Sprintf("%x", string(m.Id)), mux.Vars(req)["type"]), 302)
}
