package main

import (
	"context"
	"log"
	"math"
	"net"
	"time"

	pb "github.com/Fattouche/Chatty/protobuff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	GRPC_PORT              = ":50000"
	HEALTH_PORT            = ":51000"
	MAX_HEALTHCHECK_MISSES = 5
)

type loadbalancer struct{}

var nodeList *pb.NodeList
var roundRobinIndex int

func (s *loadbalancer) RegisterNode(ctx context.Context, node *pb.Node) (*pb.NodeList, error) {
	nodeList.Nodes = append(nodeList.Nodes, node)
	return nodeList, nil
}

func (s *loadbalancer) RendevouszServerIP(ctx context.Context, req *pb.Request) (*pb.Node, error) {
	roundRobinIndex := math.Mod(float64(roundRobinIndex+1), float64(len(nodeList.Nodes)))
	node := nodeList.Nodes[int(roundRobinIndex)]
	return node, nil
}

// It would probably be better to use UDP for this healthcheck because it puts less strain on both sides
// But this is much easier to implement than trying to handle several different connections.
func healthCheck() {
	for {
		for _, node := range nodeList.Nodes {
			conn, err := grpc.Dial(node.IP, grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to dial to %s with %v\n", node, err)
			}
			defer conn.Close()
			c := pb.NewRendezvousClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err = c.HealthCheck(ctx, &pb.Request{})
			if err != nil {
				log.Printf("Node %v failed healthcheck, removing from list", node.IP)
				removeNodeFromList(node)
			}
			//Don't want to dos the server
			time.Sleep(1 * time.Second)
		}
	}
}

func removeNodeFromList(deadNode *pb.Node) {
	for i, node := range nodeList.Nodes {
		if node.IP == deadNode.IP {
			nodeList.Nodes = append(nodeList.Nodes[:i], nodeList.Nodes[i+1:]...)
		}
	}
}

// Starts a generic GRPC server
func StartGRPCServer() {
	nodeList = new(pb.NodeList)
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLoadBalancerServer(s, &loadbalancer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go healthCheck()
	StartGRPCServer()
}
