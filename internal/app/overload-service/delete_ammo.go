// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: overload.proto

package overload

import (
	"context"

	desc "github.com/x-yield/over-api/pkg/overload-service"
)

func (i *Implementation) DeleteAmmo(ctx context.Context, req *desc.DeleteAmmoRequest) (*desc.DeleteAmmoResponse, error) {
	return i.Service.DeleteAmmo(req)
}
