// Package chtime resets the modtime for all files in a directory to the derived
// timestamp from its BSON ID.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/uhthomas/sana/internal/fastwalk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Main(ctx context.Context) error {
	dir := flag.String("dir", "", "")
	flag.Parse()

	return fastwalk.Walk(*dir, func(path string, m os.FileMode) error {
		if !m.IsRegular() {
			return nil
		}
		name := filepath.Base(path)
		name = name[:strings.IndexByte(name, '.')]

		id, err := primitive.ObjectIDFromHex(name)
		if err != nil {
			return err
		}
		return os.Chtimes(path, now, id.Timestamp())
	})
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal(err)
	}
}
