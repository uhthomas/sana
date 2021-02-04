package main

import (
	"net/http"
	"runtime"

	_ "io.6f.sana/daemon"
	"io.6f.sana/route"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// daemon.Launch()
	r := route.NewRouter()
	http.ListenAndServe(":9001", r)
}
