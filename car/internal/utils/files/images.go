package files

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	global "github.com/alserov/rently/car/internal/models"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/alserov/rently/car/internal/utils/broker"
	fstorage "github.com/alserov/rently/proto/gen/file-storage"
	"log/slog"
	"net/http"
	"sync"
)

type Imager interface {
	Save(images [][]byte, uuid string) error
	Delete(uuid string)
	GetImages(ctx context.Context, cars []models.Car[string]) error
}

func NewImager(p sarama.AsyncProducer, client fstorage.FileStorageClient, topics broker.ImageTopics, log *slog.Logger) Imager {
	return &imager{
		p:      p,
		client: client,
		topics: topics,
		log:    log,
		mu:     sync.Mutex{},
	}
}

type imager struct {
	p sarama.AsyncProducer

	log *slog.Logger

	client fstorage.FileStorageClient

	mu sync.Mutex

	topics broker.ImageTopics
}

func (i *imager) GetImages(ctx context.Context, cars []models.Car[string]) error {
	var (
		errCounter = 0
		wg         = sync.WaitGroup{}
	)

	wg.Add(len(cars))
	for idx, c := range cars {
		go func(uuid string, idx int, wg *sync.WaitGroup) {
			defer wg.Done()
			images, err := i.client.GetLinks(ctx, &fstorage.GetLinksReq{
				UUID: uuid,
			})
			if err != nil {
				i.mu.Lock()
				errCounter++
				i.mu.Unlock()
				return
			}
			cars[idx].Images = images.Links
		}(c.UUID, idx, &wg)
	}

	wg.Wait()
	if errCounter > len(cars)/2 {
		return &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to fetch images for %d cars from %d", errCounter, len(cars)),
		}
	}

	return nil
}

func (i *imager) Save(images [][]byte, uuid string) error {
	var (
		chErr = make(chan error)
		wg    = sync.WaitGroup{}
	)

	wg.Add(len(images))

	for idx, img := range images {
		go func(img []byte, idx int, wg *sync.WaitGroup) {
			defer wg.Done()
			b, err := json.Marshal(broker.SaveImageMessage{
				Value: img,
				UUID:  uuid,
				Idx:   idx,
			})
			if err != nil {
				chErr <- err
			}

			m := &sarama.ProducerMessage{
				Value: sarama.StringEncoder(b),
				Topic: i.topics.Save,
			}
			i.p.Input() <- m
		}(img, idx, &wg)
	}

	go func() {
		wg.Wait()
		close(chErr)
	}()

	if err := <-chErr; err != nil {
		return &global.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to save image: %v", err),
		}
	}

	return nil
}

func (i *imager) Delete(carUUID string) {
	m := &sarama.ProducerMessage{
		Value: sarama.StringEncoder(carUUID),
		Topic: i.topics.Delete,
	}
	i.p.Input() <- m
}
