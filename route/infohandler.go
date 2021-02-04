package route

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io.6f.sana/db"
	uppdb "upper.io/db"
)

func infoHandler(res http.ResponseWriter, req *http.Request) {
	strId := mux.Vars(req)["id"]

	id, err := hex.DecodeString(strId)
	if err != nil {
		res.WriteHeader(400)
		return
	}

	info, err := db.GetMedia(uppdb.Cond{"_id": bson.ObjectId(string(id))})
	if err != nil {
		res.WriteHeader(500)
		return
	}

	data, err := json.Marshal(info)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}
