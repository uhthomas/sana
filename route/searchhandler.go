package route

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	elastic "gopkg.in/olivere/elastic.v3"
	"io.6f.sana/db"
)

func searchHandler(res http.ResponseWriter, req *http.Request) {
	var offset int
	if s := req.URL.Query().Get("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			offset = v
		}
	}
	query := req.URL.Query().Get("q")
	results, err := db.ElasticClient.Search().
		Index("sana").
		Type("media").
		Query(elastic.NewMultiMatchQuery(query, "id", "originalId", "title", "description", "provider")).
		From(offset).Size(30).
		Do()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	var payload []db.Media
	for _, item := range results.Each(reflect.TypeOf(db.Media{})) {
		payload = append(payload, item.(db.Media))
	}
	json.NewEncoder(res).Encode(&payload)
}
