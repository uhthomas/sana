package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"io.6f.sana/db"
	uppdb "upper.io/db"
)

func resolveHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	url := "/"
	info, err := db.GetMedia(uppdb.Cond{"originalId": id})
	if err == nil {
		url = "/media/" + fmt.Sprintf("%x", string(info.Id))
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
