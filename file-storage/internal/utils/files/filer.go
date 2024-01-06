package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const RelativeImageDir = "internal/images"

type Filer interface {
	Save(file []byte, uuid string, idx int) error
	Delete(uuid string) error

	GetLinks(uuid string) ([]string, error)
	GetImage(link string) ([]byte, error)
}

func NewFiler(relativeDir string) Filer {
	return &filer{
		mu:          sync.RWMutex{},
		relativeDir: relativeDir,
	}
}

type filer struct {
	mu sync.RWMutex

	relativeDir string
}

func (flr *filer) GetLinks(uuid string) ([]string, error) {
	dir := fmt.Sprintf("%s/%s", flr.relativeDir, uuid)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read images dir: %v", err)
	}

	var links []string
	for _, f := range files {
		links = append(links, fmt.Sprintf("images?uuid=%s&idx=%s", uuid, f.Name()[:len(f.Name())-4]))
	}

	return links, nil
}

func (flr *filer) Delete(uuid string) error {
	dir := fmt.Sprintf("%s/%s", flr.relativeDir, uuid)

	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("folder with this uuid does not exist")
	}

	flr.mu.Lock()
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove dir: %v", err)
	}
	flr.mu.Unlock()

	return nil
}

func (flr *filer) Save(file []byte, uuid string, idx int) error {
	dir := fmt.Sprintf("%s/%s/", flr.relativeDir, uuid)

	fmt.Println(uuid)
	flr.mu.RLock()
	if _, err := os.Stat(flr.relativeDir); err != nil {
		if err = os.MkdirAll(flr.relativeDir, 0644); err != nil {
			return fmt.Errorf("failed to make relative dir: %v", err)
		}
	}
	flr.mu.RUnlock()

	flr.mu.Lock()
	if err := os.MkdirAll(dir, 0644); err != nil {
		return fmt.Errorf("failed to make dir: %v", err)
	}

	f, err := os.OpenFile(fmt.Sprintf("%s%d.bin", dir, idx), os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	flr.mu.Unlock()

	_, err = f.Write(file)
	if err != nil {
		return fmt.Errorf("failed to write into file: %v", err)
	}

	return nil
}

func (flr *filer) GetImage(link string) ([]byte, error) {
	file := fmt.Sprintf("%s/%s", flr.relativeDir, link)

	b, err := os.ReadFile(file)
	if err != nil {

		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	return b, nil
}
