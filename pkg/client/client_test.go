package client

import (
	"context"
	"testing"

	discoveryv0 "github.com/dmah42/thermosphere/pkg/api/v0/discovery"
)

const ()

func TestDiscoveryHealth(t *testing.T) {
	c, err := New(context.Background())
	if err != nil {
		t.Errorf("Cannot create client")
	}

	req := &discoveryv0.HealthRequest{}

	d, err := c.Discovery()
	if err != nil {
		t.Errorf("Discovery service not available")
	}

	rsp, err := d.Health(context.Background(), req)
	if err != nil {
		t.Errorf("Health should have NOT returned an error")
	}

	if !rsp.Healthy {
		t.Errorf("Health should not be unhealthy")
	}
	if rsp.Version != "0.1.0" {
		t.Errorf("expected version %q, got %q", "0.1.0", rsp.Version)
	}
}
