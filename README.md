# thermosphere

A small Go tool for discovering and health-checking [Aurae](https://aurae.io) runtime nodes on a
local network.

[Aurae](https://github.com/aurae-runtime/aurae) is a Rust-based systems runtime for managing
processes, pods, and VMs, exposed over gRPC with mTLS authentication. thermosphere acts as a
gRPC client to that API: it scans a CIDR block, calls each host's `Discovery.Health` RPC, and
reports which nodes are healthy over a simple HTTP endpoint.

## How it works

1. Scans every host address in a given CIDR block (`-cidr`, default `192.168.178.0/24`).
2. Opens an mTLS gRPC connection to each host and calls the Aurae `DiscoveryService.Health` RPC.
3. Tracks which hosts responded healthy.
4. Serves the current set of healthy nodes as plain text at `/nodes` over HTTP (`-port`,
   default `4321`).

## Layout

- `cmd/main.go` — the CLI/daemon entrypoint: CIDR scanning, health checks, and the `/nodes` HTTP
  handler.
- `pkg/client` — builds an mTLS gRPC connection to an Aurae instance.
- `pkg/config` — connection configuration (protocol/socket, TLS cert paths), with defaults for a
  local Aurae install (Unix socket `/var/run/aurae/aurae.sock`, PKI under `~/.aurae/pki/`).
- `pkg/discovery` — wraps the Aurae `DiscoveryService.Health` RPC.
- `pkg/api/v0/{discovery,observe,runtime}` — generated protobuf/gRPC client stubs for Aurae's API.

## Building

Protobuf/gRPC stubs under `pkg/api` are generated from Aurae's upstream API definitions using
[`buf`](https://docs.buf.build/installation):

```sh
make proto   # regenerate pkg/api from aurae-runtime/aurae
make test    # regenerate protos, then run tests
make mod     # go mod tidy / vendor / download
make clean   # remove generated pkg/api
```

## Running

Requires mTLS client credentials for the target Aurae node(s), matching `pkg/config`'s defaults
(`~/.aurae/pki/ca.crt`, `~/.aurae/pki/_signed.client.nova.crt`, `~/.aurae/pki/client.nova.key`).

```sh
go run ./cmd -cidr 192.168.1.0/24 -port 4321
```

Then query discovered nodes:

```sh
curl http://localhost:4321/nodes
```

## Status

Early/experimental. Discovery is currently a naive sequential scan over every host in the CIDR
block (no concurrency yet). A `TODO` in `cmd/main.go` notes a possible future move to Aurae's own
`ae` CLI/SDK.
