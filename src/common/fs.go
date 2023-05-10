package common

import (
	"io/fs"
	"os"
	"path"
)

func ChmodRec(root string) {
	fileSystem := os.DirFS(root)

	fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
		//TODO: remove unnecesary variable
		fullPath := path.Join(filePath, d.Name())

		os.Chmod(fullPath, 0777)
		return nil
	})
}
