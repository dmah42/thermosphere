package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	// TODO: port to ae
	"github.com/dmah42/thermosphere/pkg/client"
	"github.com/dmah42/thermosphere/pkg/config"

	discoveryv0 "github.com/dmah42/thermosphere/pkg/api/v0/discovery"
)

var (
	cidr = flag.String("cidr", "192.168.178.0/24", "cidr block to scan")
	port = flag.Int("port", 4321, "port on which to listen")

	nodes map[string]bool = make(map[string]bool)
)

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func discoverNodes(ctx context.Context) {
	// TODO: concurrency: see https://gist.github.com/kotakanbe/d3059af990252ba89a82

	hosts, err := hosts(*cidr)
	if err != nil {
		log.Fatal(err)
	}

	for _, ip := range hosts {
		c, err := client.New(ctx, config.WithSystem(config.System{Protocol: "tcp4", Socket: ip}))
		if err != nil {
			log.Fatal(err)
		}

		d, err := c.Discovery()
		if err != nil {
			log.Fatal(err)
		}

		rsp, err := d.Health(ctx, &discoveryv0.HealthRequest{})
		if err != nil {
			log.Fatal(err)
		}

		if rsp.Healthy {
			nodes[ip] = true
		} else {
			delete(nodes, ip)
		}
	}
}

func nodesHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%+v", nodes)
}

func main() {
	flag.Parse()

	ctx := context.Background()

	discoverNodes(ctx)

	http.HandleFunc("/nodes", nodesHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
