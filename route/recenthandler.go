package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"io.6f.sana/db"
	uppdb "upper.io/db"
)

func recentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragma", "no-cache")
	var offset uint
	if s := r.URL.Query().Get("offset"); s != "" {
		n, err := strconv.ParseUint(s, 10, 0)
		if err == nil {
			offset = uint(n)
		}
	}
	var m []db.Media
	if err := db.MediaCollection.
		Find(uppdb.Cond{
			"uploaded": uppdb.Cond{
				"$gt": time.Date(2013, time.January, 0, 0, 0, 0, 0, time.UTC),
			},
		}).
		Sort("uploaded").
		Skip(offset).
		Limit(20).
		All(&m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(m)
}
