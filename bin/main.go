// Package main starts up a gRPC server which serves the EventSpreadService RPC.
// This is an example -- separate implementations will define custom dispatchers
// (which will call custom EventSpreadHandlers).
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/cripplet/event-spread/lib/core/handlers"
	"github.com/cripplet/event-spread/lib/core/spread"
	espb "github.com/cripplet/event-spread/lib/proto/event_spread_go_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	host = flag.String("host", "localhost", "EventSpreadService gRPC host")
	port = flag.Int("port", 8080, "EventSpreadService gRPC port")

	dispatcher = map[espb.SpreadType]handlers.EventSpreadHandler{
		espb.SpreadType_SPREAD_TYPE_INSTANT_GLOBAL: &handlers.InstantGlobalEventSpreadHandler{},
	}
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	s := grpc.NewServer()
	srv := spread.NewEventSpreadService(dispatcher)
	espb.RegisterEventSpreadServiceServer(s, srv)
	reflection.Register(s)

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start server on %s with error %v", addr, err)
	}

	log.Printf("serving on %s", addr)

	log.Fatal(s.Serve(conn))
}
