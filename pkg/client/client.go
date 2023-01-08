package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dmah42/thermosphere/pkg/config"
	"github.com/dmah42/thermosphere/pkg/discovery"
)

var ErrOptionNotAvailable = errors.New("option is not available")

type AuraeClient interface {
	Discovery() (discovery.Discovery, error)
}

type auraeClient struct {
	cfg       *config.Configs
	conn      grpc.ClientConnInterface
	discovery discovery.Discovery
}

func New(ctx context.Context, cfg ...config.Config) (AuraeClient, error) {
	cf, err := config.From(cfg...)
	if err != nil {
		log.Fatal("Cannot initialize config", err)
	}

	tlsCredentials, err := loadTLSCredentials(cf.Auth)
	if err != nil {
		log.Fatal("Cannot load TLS credentials", err)
	}

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := net.Dialer{}

		return d.DialContext(ctx, cf.System.Protocol, addr)
	}

	conn, err := grpc.Dial(
		cf.System.Socket,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		log.Fatal("Cannot Dial", err)

		return nil, err
	}

	d, err := discovery.New(ctx, conn)
	if err != nil {
		log.Fatal("Cannot crete discovery client", err)

		return nil, err
	}

	c := &auraeClient{
		cfg:       cf,
		conn:      conn,
		discovery: d,
	}

	return c, nil
}

func loadTLSCredentials(auth config.Auth) (credentials.TransportCredentials, error) {
	caPEM, err := os.ReadFile(auth.CaCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	clientKeyPair, err := tls.LoadX509KeyPair(auth.ClientCert, auth.ClientKey)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientKeyPair},
		RootCAs:      certPool,
		ServerName:   auth.ServerName,
	}

	return credentials.NewTLS(config), nil
}

func (c *auraeClient) Discovery() (discovery.Discovery, error) {
	if c.discovery == nil {
		return nil, fmt.Errorf("configuration: %w", ErrOptionNotAvailable)
	}

	return c.discovery, nil
}
