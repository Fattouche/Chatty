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

// Used to implement the grpc loadbalancer server
type loadbalancer struct{}

// Keeps track of which nodes are currently alive
var nodeList *pb.NodeList
var roundRobinIndex int

// Registers a new node and notifies existing nodes of its arrival
func (s *loadbalancer) RegisterNode(ctx context.Context, node *pb.Node) (*pb.NodeList, error) {
	nodeList.Nodes = append(nodeList.Nodes, node)
	if len(nodeList.Nodes) > 1 {
		notifyNodeUpdate(node, true)
	}
	log.Println("New node arrival ", node.IP)
	return nodeList, nil
}

// Returns the IP of a node to one of the peers
func (s *loadbalancer) RendevouszServerIP(ctx context.Context, req *pb.Request) (*pb.Node, error) {
	roundRobinIndex := math.Mod(float64(roundRobinIndex+1), float64(len(nodeList.Nodes)))
	node := nodeList.Nodes[int(roundRobinIndex)]
	return node, nil
}

// When a node arrives or is removed, this will be called
func notifyNodeUpdate(updatedNode *pb.Node, arrival bool) {
	// The node that will be responsible for telling all the other nodes about the update
	node := nodeList.Nodes[0]
	// Setup GRPC connection
	conn, err := grpc.Dial(node.IP, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v\n", node, err)
	}
	defer conn.Close()
	// Create a client that can access the node functions
	c := pb.NewRendezvousClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	updatedNode.NotifyOthers = 1
	if arrival {
		_, err = c.NodeArrival(ctx, updatedNode)
	} else {
		_, err = c.NodeRemoval(ctx, updatedNode)
	}
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
			// Do the healthcheck
			_, err = c.HealthCheck(ctx, &pb.Request{})
			if err != nil {
				log.Printf("Node %v failed healthcheck, removing from list", node.IP)
				// If the node fails the healthcheck, remove it from our list of healthy nodes
				removeNodeFromList(node)
			}
			//Don't want to dos the server
			time.Sleep(1 * time.Second)
		}
	}
}

// A helper function to remove a node from our node list
func removeNodeFromList(deadNode *pb.Node) {
	for i, node := range nodeList.Nodes {
		// Since both pointers, just compares their address
		if node == deadNode {
			nodeList.Nodes = append(nodeList.Nodes[:i], nodeList.Nodes[i+1:]...)
			if len(nodeList.Nodes) > 0 {
				notifyNodeUpdate(deadNode, false)
			}
		}
	}
}

// Starts a generic GRPC server for peers to hit
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
