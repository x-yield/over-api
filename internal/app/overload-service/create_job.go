// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: overload.proto

package overload

import (
	"context"

	desc "github.com/x-yield/over-api/pkg/overload-service"
)

func (i *Implementation) CreateJob(ctx context.Context, req *desc.CreateJobRequest) (*desc.CreateJobResponse, error) {
	return i.Service.CreateJob(req)
}
