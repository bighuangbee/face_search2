package storage

import (
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"time"
)

type RegisteInfo face_wrapper.RegisteInfo

//type RegisteInfo struct {
//	Time      string `json:"time"` //注册时间
//	Ok        bool   `json:"ok"`
//	Filename  string `json:"filename"` //拍摄时间
//	ShootTime time.Time
//}

type FaceDb interface {
	Update(key string, value *RegisteInfo) error
	Read(key string) (value *RegisteInfo, ok bool)
	ReadBatch() (values []*RegisteInfo, err error)
	Delete(key string) error
	DeleteExpired(effectiveTime time.Duration) (values []*RegisteInfo, err error)
	Count() int64
}
