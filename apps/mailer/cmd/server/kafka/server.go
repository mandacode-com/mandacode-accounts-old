package kafkaserver

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaServer interface {
	Start() error
	Stop() error
}

type ReaderHandler struct {
	Reader  *kafka.Reader
	Handler KafkaHandler
}

type kafkaServer struct {
	readerHandlers []*ReaderHandler
	logger         *zap.Logger
	wg             *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

// Start implements KafkaServer.
func (k *kafkaServer) Start() error {
	k.logger.Info("Starting Kafka server")

	for _, readerHandler := range k.readerHandlers {
		k.wg.Add(1)
		go func(rh *ReaderHandler) {
			defer k.wg.Done()
			k.logger.Info("Starting reader", zap.String("topic", rh.Reader.Config().Topic))
			for {
				m, err := rh.Reader.ReadMessage(k.ctx)
				if err != nil {
					if k.ctx.Err() != nil {
						k.logger.Info("Context cancelled, stopping reader", zap.String("topic", rh.Reader.Config().Topic))
						return
					}
					k.logger.Error("Error reading message", zap.Error(err), zap.String("topic", rh.Reader.Config().Topic))
					continue
				}

				if err := rh.Handler.HandleMessage(m); err != nil {
					k.logger.Error("Error handling message", zap.Error(err), zap.String("topic", rh.Reader.Config().Topic), zap.ByteString("key", m.Key), zap.ByteString("value", m.Value))
					continue
				}
			}
		}(readerHandler)
	}

	k.wg.Wait()
	return nil
}

// Stop implements KafkaServer.
func (k *kafkaServer) Stop() error {
	k.logger.Info("Stopping Kafka server")
	k.cancel() // Cancel the context to stop all readers

	// Wait for all goroutines to finish
	k.wg.Wait()
	k.logger.Info("Kafka server stopped")
	return nil
}

func NewKafkaServer(logger *zap.Logger, readerHandlers []*ReaderHandler) KafkaServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &kafkaServer{
		readerHandlers: readerHandlers,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		ctx:            ctx,
		cancel:         cancel,
	}
}
