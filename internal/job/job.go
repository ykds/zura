package job

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/comet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type Server struct {
	cfg            Config
	uncommitMsgMap sync.Map
	cometClient    comet.CometClient
	consumersMap   map[string]*kafka.Reader
}

func NewJobServer(cfg Config) *Server {
	server := &Server{
		cfg:            cfg,
		uncommitMsgMap: sync.Map{},
		consumersMap:   map[string]*kafka.Reader{},
	}

	for k, v := range cfg.Kafka.GroupTopicMap {
		consumer := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.Kafka.Brokers,
			GroupID:        k,
			Topic:          v,
			CommitInterval: 5 * time.Second,
			Logger:         kafka.LoggerFunc(log.GetGlobalLogger().Infof),
			ErrorLogger:    kafka.LoggerFunc(log.GetGlobalLogger().Errorf),
		})
		server.consumersMap[v] = consumer
		go func() {
			for {
				message, err := consumer.FetchMessage(context.Background())
				if err != nil {
					log.Errorf("fetch message err: %+v", err)
					continue
				}
				server.uncommitMsgMap.Store(string(message.Key), message)
			}
		}()
		return server
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	cometConn, err := grpc.DialContext(ctx2, fmt.Sprintf("%s:%s", cfg.CometServer.Host, cfg.CometServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("new comet grpc client failed, err: %+v", err)
	}
	cancel2()
	server.cometClient = comet.NewCometClient(cometConn)

	return server
}
