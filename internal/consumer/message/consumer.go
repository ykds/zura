package message

import (
	"encoding/json"
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/pkg/discovery"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/comet"
	"github.com/ykds/zura/proto/logic"
	"github.com/ykds/zura/proto/protocol"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"strings"
	"sync"
	"time"
)

type messageConsumer struct {
	ctx    context.Context
	comets map[int32]*grpc.ClientConn
	m      sync.RWMutex
	dis    *discovery.EtcdManager
}

func (m *messageConsumer) Consume(message kafka2.Message) error {
	pushmsg := logic.PushMsg{}
	err := json.Unmarshal(message.Value, &pushmsg)
	if err != nil {
		return err
	}
	m.m.RLock()
	if conn, ok := m.comets[pushmsg.Server]; ok {
		m.m.RUnlock()
		client := comet.NewCometClient(conn)
		body, _ := json.Marshal(pushmsg.Message)
		_, err := client.PushMessage(context.Background(), &comet.PushMsgRequest{
			Op:       pushmsg.Op,
			ToUserId: pushmsg.ToUserId,
			Proto: &protocol.Protocol{
				Op:   pushmsg.Op,
				Body: body,
			},
		})
		return err
	}
	m.m.RUnlock()
	return nil
}

func NewConsumer(ctx context.Context, dis *discovery.EtcdManager) kafka.ConsumerHandler {
	c := &messageConsumer{
		ctx:    ctx,
		comets: make(map[int32]*grpc.ClientConn),
		dis:    dis,
	}
	go c.watchComet()
	return c
}

func (m *messageConsumer) watchComet() {
	fn := func(update []*endpoints.Update) {
		for _, item := range update {
			switch item.Op {
			case endpoints.Add:
				service := strings.Replace(item.Key, common.CometDiscoveryEndpoint+"/", "", -1)
				serverId, _ := strconv.ParseInt(service, 10, 64)
				ctx2, cancel2 := context.WithTimeout(m.ctx, 2*time.Second)
				cometConn, err := grpc.DialContext(ctx2,
					item.Endpoint.Addr,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
				)
				if err != nil {
					log.Errorf("new comet grpc client failed, err: %+v", err)
					cancel2()
					continue
				}
				cancel2()
				m.m.Lock()
				m.comets[int32(serverId)] = cometConn
				m.m.Unlock()
			case endpoints.Delete:
				service := strings.Replace(item.Key, common.CometDiscoveryEndpoint+"/", "", -1)
				serverId, _ := strconv.ParseInt(service, 10, 64)
				m.m.Lock()
				if conn, ok := m.comets[int32(serverId)]; ok {
					_ = conn.Close()
					delete(m.comets, int32(serverId))
				}
				m.m.Unlock()
			}
		}
	}
	err := m.dis.Watch(m.ctx, fn)
	if err != nil {
		return
	}
}
