// Package findunique recursively maps all files in a directory by its size.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/uhthomas/sana/internal/fastwalk"
)

func Main(ctx context.Context) error {
	dir := flag.String("dir", ".", "")
	flag.Parse()

	m := map[int64][]string{}
	var mu sync.Mutex

	if err := fastwalk.Walk(*dir, func(path string, typ os.FileMode) error {
		if !typ.IsRegular() {
			return nil
		}
		d, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
		mu.Lock()
		m[d.Size()] = append(m[d.Size()], path)
		mu.Unlock()
		return nil
	}); err != nil {
		return fmt.Errorf("walk: %w", err)
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	return e.Encode(m)
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
