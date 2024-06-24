package data

import (
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/go-kratos/kratos/v2/log"
)

type Data struct {
}

func NewData(bc *conf.Bootstrap, logger log.Logger) (*Data, error) {

	data := &Data{}

	return data, nil
}
