package main

import (
	"context"
	"flag"
	rb "github.com/fflow-tech/fflow/api/foundation/rbac"
	"log"
	"testing"
	"time"

	pb "github.com/fflow-tech/fflow/api/foundation/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50042", "the address to connect to")
)

func TestClient(t *testing.T) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ValidateToken(ctx, &pb.ValidateTokenReq{
		BasicReq: &pb.BasicReq{Namespace: "hunter", AccessToken: "123"},
	})
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	log.Printf("Result: %s", r.String())
}

func TestClient2(t *testing.T) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := rb.NewRbacClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.HasPermission(ctx, &rb.RbacReq{User: "test"})
	if err != nil {
		t.Fatalf("Failed to has permission: %v", err)
	}
	log.Printf("Result: %s", r.String())
}
