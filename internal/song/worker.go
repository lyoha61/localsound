package song

import (
	"context"

	"go.uber.org/zap"
)

type Worker struct {
	service *Service
	logger  *zap.SugaredLogger
}

func NewWorker(service *Service, logger *zap.SugaredLogger) *Worker {
	return &Worker{
		service: service,
		logger:  logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Infof("Starting song worker")

	for {
		select {
		case <-ctx.Done():
			w.logger.Infof("Stopping song worker")
			return
		case task, ok := <-w.service.importQueue:
			if !ok {
				w.logger.Infof("Song import queue closed")
				return
			}
			w.processTask(ctx, task)
		}
	}
}

func (w *Worker) processTask(ctx context.Context, task ImportTask) {
	w.logger.Infof("Processings import task: %s - %s", task.Artist, task.Title)
}
