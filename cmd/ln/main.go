// Package ln hard links all files from the given JSON array to the
// target directory.
//
// The array can be generated using findunique and jq
//	jq 'map(select(length <= 2)) | add' < out.json > real.json
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Main(ctx context.Context) error {
	in := flag.String("in", "", "")
	target := flag.String("target", "/mnt/d/appsrc/io.6f.sana/_/media_real", "")
	flag.Parse()

	var r io.Reader
	switch in := *in; in {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(in)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()
		r = f
	}

	var paths []string
	if err := json.NewDecoder(r).Decode(&paths); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	for _, p := range paths {
		abs, err := filepath.Abs(filepath.Join(*target, filepath.Base(p)))
		if err != nil {
			return fmt.Errorf("abs: %w", err)
		}
		if err := os.Link(p, abs); err != nil {
			return fmt.Errorf("symlink: %w", err)
		}
	}
	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
