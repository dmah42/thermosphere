package discovery

import (
	"context"
	"log"

	"google.golang.org/grpc"

	discoveryv0 "github.com/dmah42/thermosphere/pkg/api/v0/discovery"
)

type Discovery interface {
	Health(context.Context, *discoveryv0.HealthRequest) (*discoveryv0.HealthResponse, error)
}

type discovery struct {
	client discoveryv0.DiscoveryServiceClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) (Discovery, error) {
	d := &discovery{
		client: discoveryv0.NewDiscoveryServiceClient(conn),
	}

	return d, nil
}

func (d *discovery) Health(ctx context.Context, req *discoveryv0.HealthRequest) (*discoveryv0.HealthResponse, error) {
	rsp, err := d.client.Health(ctx, req)
	if err != nil {
		log.Fatal("Cannot call Health ", err)

		return nil, err
	}

	return rsp, nil
}
