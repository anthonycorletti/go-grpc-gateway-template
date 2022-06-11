package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apipb "github.com/anthonycorletti/go-grpc-gateway-template/proto/api"
	_ "github.com/anthonycorletti/go-grpc-gateway-template/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
)

type Server struct {
	apipb.UnimplementedMessengerServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SendMessage(ctx context.Context, req *apipb.RequestContent) (*apipb.ResponseContent, error) {
	return &apipb.ResponseContent{Message: "Hello " + req.Name}, nil
}

func ServeSwaggerDocs(mux *http.ServeMux) {
	prefix := "/docs/"
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(statikFS)))
}

func main() {
	log.Printf("starting api")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}

	s := grpc.NewServer()
	// Register the messenger server with the gRPC server.
	apipb.RegisterMessengerServer(s, NewServer())

	log.Println("starting grpc server...")
	go func() {
		log.Fatal(s.Serve(lis))
	}()
	log.Println("server started")

	// Create a client
	log.Print("starting grpc client...")
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("failed to dial:", err)
	}

	log.Println("registering handlers on api gateway...")
	gwmux := runtime.NewServeMux()
	// Register Messenger
	err = apipb.RegisterMessengerHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("failed to register gateway:", err)
	}
	log.Println("completed handler registration on api gateway")

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	ServeSwaggerDocs(mux)

	log.Println("starting api server")
	gwServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	log.Println("http server started on port 8081")
	log.Fatalln(gwServer.ListenAndServe())
}
