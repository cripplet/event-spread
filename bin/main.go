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
)

var (
	host = flag.String("host", "localhost", "EventSpreadService gRPC host")
	port = flag.Int("port", 8080, "EventSpreadService gRPC port")
	serveGRPC = flag.Bool("grpc", true, "Specify if the server should run using the gRPC or REST protocol")

	dispatcher = map[espb.SpreadType]handlers.EventSpreadHandler{
		espb.SpreadType_SPREAD_TYPE_INSTANT_GLOBAL: &handlers.InstantGlobalEventSpreadHandler{},
	}
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start server on %s with error %v", addr, err)
	}

	s := grpc.NewServer()
	espb.RegisterEventSpreadServiceServer(s, spread.NewEventSpreadService(dispatcher))

	log.Printf("serving %s traffic on %s", map[bool]string{false: "REST", true: "gRPC"}[*serveGRPC], addr)

	if *serveGRPC {
		s.Serve(conn)
	} else {
		// TODO(cripplet): Implement this.
		// s.serveHTTP(...)
	}
}
