package main

import (
    "context"
    "fmt"
    "net"

    "google.golang.org/grpc"
    _ "google.golang.org/grpc/encoding/gzip"

    proto "github.com/danslimmon/yikes/proto/go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type traceServiceServer struct {
    proto.UnimplementedTraceServiceServer
}

// Export is the server-side implementation of the Export method.
//
// It accepts trace spans as they come in from the client.
func (s *traceServiceServer) Export(ctx context.Context, in *proto.ExportTraceServiceRequest) (*proto.ExportTraceServiceResponse, error) {
    spans := in.GetResourceSpans()
    fmt.Printf("spans: %v\n", spans)
    return new(proto.ExportTraceServiceResponse), nil
}

func newServer() (s *traceServiceServer) {
    return new(traceServiceServer)
}

func main() {
    lis, err := net.Listen("tcp", ":44317")
    if err != nil {
        panic(err)
    }
    fmt.Println("Listening on port 44317")

    grpcServer := grpc.NewServer()
    proto.RegisterTraceServiceServer(grpcServer, newServer())
    fmt.Printf("service info: %v\n", grpcServer.GetServiceInfo())
    grpcServer.Serve(lis)
    fmt.Println("Exiting")
}
