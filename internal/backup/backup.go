package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/samber/do/v2"
	"github.com/willie68/osmltools/internal/logging"
)

// Backup the sd card backup service
type Backup struct {
	log logging.Logger
}

// Init init this service and provide it to di
func Init(inj do.Injector) {
	bck := Backup{
		log: *logging.New().WithName("Backup"),
	}
	do.ProvideValue(inj, bck)
}

// Backup backup all files from the sd card into a zip file
func (b *Backup) Backup(sdCardFolder, outputFolder string) error {
	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}

	of := filepath.Join(outputFolder, fmt.Sprintf("bck_%s.zip", time.Now().Format("20060102150405")))

	file, err := os.Create(of)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		name := path[len(sdCardFolder)+1:]
		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(name)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(sdCardFolder, walker)
	if err != nil {
		return err
	}
	return nil
}
