package data

import (
	"fmt"
	"github.com/bighuangbee/face_search2/pkg/conf"
	basePb "github.com/bighuangbee/kratos-vue-admin/api/admin/v1"
	"github.com/go-kratos/kratos/v2/log"
	grpcDial "google.golang.org/grpc"
)

type Data struct {
	mircoSevices map[string]*grpcDial.ClientConn

	sysuserClient basePb.SysuserClient
}

func NewData(bc *conf.Bootstrap, logger log.Logger) (*Data, error) {

	data := &Data{}

	if err := data.initMicroServices(bc, logger); err != nil {
		panic(err)
	}

	return data, nil
}

func (data *Data) SysuserGrpcClient() basePb.SysuserClient {
	return data.sysuserClient
}

func (data *Data) initMicroServices(bc *conf.Bootstrap, logger log.Logger) error {
	grpcConnMap := make(map[string]*grpcDial.ClientConn)

	fmt.Println("---bc.MicroServices", bc.MicroServices)

	data.mircoSevices = grpcConnMap

	data.sysuserClient = basePb.NewSysuserClient(data.mircoSevices["baseService"])

	return nil
}
