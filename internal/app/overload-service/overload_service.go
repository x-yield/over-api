// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: overload.proto

package overload

import (
	"github.com/utrack/clay/v2/transport"
	desc "github.com/x-yield/over-api/pkg/overload-service"
	"github.com/x-yield/over-api/service"
)

type Implementation struct {
	Service *service.OverloadService
}

// NewOverloadService create new Implementation
func NewOverloadService(overloadSrv *service.OverloadService) *Implementation {
	return &Implementation{
		Service: overloadSrv,
	}
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *Implementation) GetDescription() transport.ServiceDesc {
	return desc.NewOverloadServiceServiceDesc(i)
}
