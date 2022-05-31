package service

import (
	"flag"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	authorManagement "github.com/wcodesoft/author-management-service/grpc/go/author-management.proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"service/database"
	"service/routes"
)

var (
	port    = flag.Int("port", 9000, "Port on which gRPC server should listen TCP conn.")
	webPort = flag.Int("webPort", 9001, "Port on which web gRPC server should listen TCP conn.")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Initializing gRPC server on port %d", *port)

	gRPCServer := grpc.NewServer()

	db := database.NewDatabase()
	authorManagement.RegisterAuthorManagementServer(gRPCServer, routes.NewRouteServer(db))

	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Now wrapping to server for web clients using Typescript gRPC
	wrappedServer := grpcweb.WrapServer(gRPCServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		if wrappedServer.IsGrpcWebRequest(req) {
			wrappedServer.ServeHTTP(resp, req)
			return
		}
		http.DefaultServeMux.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", *webPort),
		Handler: cors.AllowAll().Handler(http.HandlerFunc(handler)),
	}

	log.Printf("Starting gRPC Web server. http port :%d\n", *webPort)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed starting http server: %v", err)

	}
}
