package lyrics

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	service  *Service
	interval time.Duration
	logger   *zap.SugaredLogger
}

func NewWorker(service *Service, interval time.Duration, logger *zap.SugaredLogger) *Worker {
	return &Worker{
		service:  service,
		interval: interval,
		logger:   logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Infof("Starting lyrics worker with interval %v", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Stopping lyrics worker")
			return
		case <-ticker.C:
			if err := w.service.ProcessPending(ctx); err != nil && ctx.Err() == nil {
				w.logger.Errorf("Error in lyrics worker processPending: %v", err)
			}
		}
	}
}
