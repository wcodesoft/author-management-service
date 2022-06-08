package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"service/database"
	"service/routes"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	authorManagement "github.com/wcodesoft/author-management-service/grpc/go/author-management.proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
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

	dbConnectorString, ok := os.LookupEnv("DB_CONNECTOR")
	if !ok {
		dbConnectorString = "postgres://postgres:postgrespw@localhost:55000"
	}
	log.Printf("Trying to connect to database at: %s\n", dbConnectorString)
	postgresDialector := postgres.Open(dbConnectorString)
	db := database.NewConnection(postgresDialector)
	authorManagement.RegisterAuthorManagementServer(gRPCServer, routes.NewRoutes(db))

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

	log.Printf("Starting gRPC Web server. http port %d\n", *webPort)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed starting http server: %v", err)

	}

	log.Printf("gRPC server initialized at port %d\n", *port)
	log.Printf("gRPC Web server started at http port %d\n", *webPort)
}
