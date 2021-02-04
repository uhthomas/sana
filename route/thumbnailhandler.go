package route

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DAddYE/vips"
	"github.com/golang/groupcache"
	"github.com/gorilla/mux"
)

var thumbnailCache = groupcache.NewGroup("thumbnail", 1024<<10, groupcache.GetterFunc(
	func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
		var (
			parsedContext     = ctx.(map[string]string)
			id                = parsedContext["id"]
			size              = parsedContext["size"]
			width         int = 1920
			height        int = 1080
		)

		switch size {
		case "large":
			width, height = 1920, 1080
		case "medium":
			width, height = 1280, 720
		case "small":
			width, height = 640, 360
		case "mini":
			width, height = 300, 169
		case "tiny":
			width, height = 120, 68
		}

		options := vips.Options{
			Width:        width,
			Height:       height,
			Crop:         true,
			Enlarge:      true,
			Extend:       vips.EXTEND_WHITE,
			Interpolator: vips.BILINEAR,
			Gravity:      vips.CENTRE,
			Quality:      100,
		}

		data, err := ioutil.ReadFile(filepath.Join("_", "thumbnail", id+".jpg"))
		if err != nil {
			return err
		}

		data, err = vips.Resize(data, options)
		if err != nil {
			return err
		}

		dest.SetBytes(data)
		return nil
	}))

func thumbnailHandler(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	size := strings.Split(mux.Vars(req)["size"], ".")[0]

	thumbnailContext := map[string]string{
		"size": size,
		"id":   id,
	}

	var data []byte
	err := thumbnailCache.Get(thumbnailContext, id+"-"+size, groupcache.AllocatingByteSliceSink(&data))
	if err == os.ErrNotExist {
		res.WriteHeader(404)
		return
	} else if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	res.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(res, req, id, time.Time{}, bytes.NewReader(data))
}
