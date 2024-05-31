package data

import (
	"context"
	"fmt"
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/bighuangbee/face_search2/pkg/util/kitGrpc"
	basePb "github.com/bighuangbee/kratos-vue-admin/api/admin/v1"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	grpcDial "google.golang.org/grpc"
	"time"
)

type Data struct {
	mircoSevices map[string]*grpcDial.ClientConn

	sysuserClient basePb.SysuserClient
}

func NewData(bc *conf.Bootstrap, logger log.Logger, naming naming_client.INamingClient) (*Data, error) {

	data := &Data{}

	if err := data.initMicroServices(bc, logger, naming); err != nil {
		panic(err)
	}

	return data, nil
}

func (data *Data) SysuserGrpcClient() basePb.SysuserClient {
	return data.sysuserClient
}

func (data *Data) initMicroServices(bc *conf.Bootstrap, logger log.Logger, naming naming_client.INamingClient) error {
	grpcConnMap := make(map[string]*grpcDial.ClientConn)

	fmt.Println("---bc.MicroServices", bc.MicroServices)
	var discovery *nacos.Registry
	var filter selector.NodeFilter
	if bc.Discovery.Enable {
		discovery = nacos.New(naming)
		filter = kitGrpc.FilterVersion(bc.Version)
	}
	for k, v := range bc.MicroServices {
		var discoveryLocal *nacos.Registry
		var filterLocal selector.NodeFilter
		endpoint := v.Endpoint
		if !v.IsLocal {
			discoveryLocal = discovery
			endpoint = "discovery:///" + v.Name
			filterLocal = filter
		}

		opt := kitGrpc.GetConnOption{
			Endpoint:     endpoint,
			Logger:       logger,
			Discovery:    discoveryLocal,
			SelectFilter: filterLocal,
			Caller:       bc.Name,
			Timeout:      time.Second * time.Duration(v.Timeout),
		}
		conn, err := kitGrpc.GetGrpcClient(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("service key", k, opt, "  ", conn)
		grpcConnMap[k] = conn
	}

	data.mircoSevices = grpcConnMap

	data.sysuserClient = basePb.NewSysuserClient(data.mircoSevices["baseService"])

	return nil
}
