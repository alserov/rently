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

func (s serverConverter) ImageBytesToPb(bytes []byte) *fstorage.GetImageRes {
	//TODO implement me
	panic("implement me")
}

func (s serverConverter) LinksToPb(strings []string) *fstorage.GetLinksRes {
	//TODO implement me
	panic("implement me")
}
