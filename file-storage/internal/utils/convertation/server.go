package convertation

import fstorage "github.com/alserov/rently/proto/gen/file-storage"

type ServerConverter interface {
	ToPb
}

type ToPb interface {
	LinksToPb([]string) *fstorage.GetLinksRes
	ImageBytesToPb([]byte) *fstorage.GetImageRes
}

func NewServerConverter() ServerConverter {
	return &serverConverter{}
}

type serverConverter struct {
}

func (s serverConverter) ImageBytesToPb(imagesBytes []byte) *fstorage.GetImageRes {
	return &fstorage.GetImageRes{
		Image: imagesBytes,
	}
}

func (s serverConverter) LinksToPb(links []string) *fstorage.GetLinksRes {
	return &fstorage.GetLinksRes{
		Links: links,
	}
}
