package main

import (
	"context"
	"flag"
	"log"
	"testing"
	"time"

	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func TestClient(t *testing.T) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWorkflowClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetInstList(ctx, &pb.GetInstListReq{
		BasicReq: &pb.BasicReq{Operator: "hunter"},
		DefID:    "254",
	})
	if err != nil {
		t.Fatalf("Could not get inst list: %v", err)
	}
	log.Printf("Details: %s", r.GetInstDetails())
}
