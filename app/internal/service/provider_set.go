package service

import (
	"github.com/bighuangbee/face_search2/app/internal/service/face"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	face.NewFaceRecognizeApp,
)
