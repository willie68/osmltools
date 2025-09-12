package backup

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
)

// Restore restore all files from a zip to the sd card
func (b *Backup) Restore(zipfile, sdCardFolder string) error {
	fs, err := os.Stat(sdCardFolder)
	if err != nil {
		return err
	}
	if !fs.IsDir() {
		return errors.New("sd card is not a directory")
	}
	read, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := filepath.Join(sdCardFolder, file.Name)
		err = os.MkdirAll(filepath.Dir(name), os.ModeDir)
		if err != nil {
			return err
		}
		df, err := os.Create(name)
		if err != nil {
			return err
		}
		_, err = df.ReadFrom(open)
		if err != nil {
			return err
		}
		err = df.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
