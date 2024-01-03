package files

import (
	"bytes"
	"fmt"
	"image"
	"os"
)

type Filer interface {
	Save(file []byte, uuid string, idx int) error
	Delete(uuid string) error

	GetImage(uuid string, idx int) (image.Image, error)
}

func NewFiler(relativeDir string) Filer {
	return &filer{relativeDir: relativeDir}
}

type filer struct {
	relativeDir string
}

func (flr filer) Delete(uuid string) error {
	dir := fmt.Sprintf("%s/%s", flr.relativeDir, uuid)

	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("folder with this uuid does not exist")
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove dir: %v", err)
	}

	return nil
}

func (flr filer) Save(file []byte, uuid string, idx int) error {
	if _, err := os.Stat(flr.relativeDir); err != nil {
		if err = os.MkdirAll(flr.relativeDir, 0644); err != nil {
			return fmt.Errorf("failed to make relative dir: %v", err)
		}
	}

	dir := fmt.Sprintf("%s/%s/", flr.relativeDir, uuid)

	if err := os.MkdirAll(dir, 0644); err != nil {
		return fmt.Errorf("failed to make dir: %v", err)
	}

	f, err := os.OpenFile(fmt.Sprintf("%s%d", dir, idx), os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(file)
	if err != nil {
		return fmt.Errorf("failed to write into file: %v", err)
	}

	return nil
}

func (flr filer) GetImage(uuid string, idx int) (image.Image, error) {
	dir := fmt.Sprintf("%s/%s/%d", flr.relativeDir, uuid, idx)

	b, err := os.ReadFile(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(b))
	return img, nil
}
