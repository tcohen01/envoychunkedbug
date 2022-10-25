package main

import (
	"fmt"
	srv "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type ExtProcServer struct {
}

func (s *ExtProcServer) Process(stream srv.ExternalProcessor_ProcessServer) error {
	log.Print("process called")
	for {
		request, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Print("stream closed by envoy")
				return nil
			}
			log.Printf("stream closed with error %s", err)
			return err
		}
		response := srv.ProcessingResponse{}
		switch request.GetRequest().(type) {

		case *srv.ProcessingRequest_RequestHeaders:
			log.Print("received request headers")
			response.Response = &srv.ProcessingResponse_RequestHeaders{}
		case *srv.ProcessingRequest_RequestBody:
			log.Print("received request body")
			response.Response = &srv.ProcessingResponse_RequestBody{}
		case *srv.ProcessingRequest_RequestTrailers:
			log.Print("received request trailers")
			response.Response = &srv.ProcessingResponse_RequestTrailers{}
		case *srv.ProcessingRequest_ResponseHeaders:
			log.Print("received response headers")
			response.Response = &srv.ProcessingResponse_ResponseHeaders{}
		case *srv.ProcessingRequest_ResponseBody:
			log.Print("received response body, sending immediate response")
			response.Response = &srv.ProcessingResponse_ImmediateResponse{ImmediateResponse: &srv.ImmediateResponse{
				Status:     &typev3.HttpStatus{Code: 200},
				Headers:    nil,
				Body:       "Immediate Response Body",
				GrpcStatus: nil,
				Details:    "Immediate Response Details",
			}}
		case *srv.ProcessingRequest_ResponseTrailers:
			log.Print("received request trailers")
			response.Response = &srv.ProcessingResponse_RequestTrailers{}
		}
		if err := stream.Send(&response); err != nil {
			log.Printf("stream send has error %s", err)
			return err
		}
	}
}

const PORT = 8080

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalf("could not listen to port %d: %v", PORT, err)
	}

	server := grpc.NewServer()
	srv.RegisterExternalProcessorServer(server, &ExtProcServer{})
	if err = server.Serve(lis); err != nil {
		log.Fatalf("error serving grpc server: %v", err)
		return
	}
}
