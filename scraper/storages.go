package cloudrun

import (
	"context"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	cloud_storage "cloud.google.com/go/storage"
)

type Storage interface {
	write(filePath string)
}

type GCloudStorage struct {
	bucket          string
	destinationPath string
}

func (storage GCloudStorage) write(source string) {
	//var r io.Reader
	f, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//r = f
	f.Sync()

	//TODO: check how to run selected code parts depending on local or cloud enviroment
	//ctx := context.Background()
	//if err = storage.upload(ctx, r, source); err != nil {
	//	log.Fatal(err)
	//}
}

func (storage GCloudStorage) upload(ctx context.Context, r io.Reader, source string) error {
	client, err := cloud_storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bh := client.Bucket(storage.bucket)
	name := path.Join(storage.destinationPath, filepath.Base(source))
	obj := bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
