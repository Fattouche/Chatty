package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Fattouche/Chatty/protobuff"
	util "github.com/Fattouche/Chatty/util"
	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	NOTIFY_OTHERS      = 1
	DONT_NOTIFY_OTHERS = 2
	LB_IP              = "127.0.0.1:50000"
)

type server struct{}

var localGRPCAddr string

var nodeMap = make(map[string]int)
var peerMap map[string]*pb.Peer

func (s *server) NodeDown(ctx context.Context, node *pb.Node) (*pb.Response, error) {
	var status string
	delete(nodeMap, node.IP)
	if node.GetNotifyOthers() == 1 {
		err := notifyDeadNode(node.GetIP())
		if err != nil {
			status = err.Error()
		} else {
			status = "ok"
		}
	}
	return &pb.Response{Status: status}, nil
}

func (s *server) PeerArrival(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	peerMap[peer.GetName()] = peer
	return &pb.Response{Status: "ok"}, nil
}

func (s *server) PeerRemoval(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	delete(peerMap, peer.GetName())
	return &pb.Response{Status: "ok"}, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Status: "ok"}, nil
}

func notifyDeadNode(deadNodeIP string) error {
	var err error
	for node := range nodeMap {
		conn, err := grpc.Dial(node, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to %s with %v", node, err)
		}
		defer conn.Close()
		c := pb.NewRendezvousClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.NodeDown(ctx, &pb.Node{IP: deadNodeIP, NotifyOthers: DONT_NOTIFY_OTHERS})
		if err != nil {
			log.Printf("Error from %s of %v", node, err)
		}
	}
	//only returns the last error unfortunately
	return err
}

func notifyPeerUpdate(peer *pb.Peer, arriving bool) error {
	var err error
	for node := range nodeMap {
		conn, err := grpc.Dial(node, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to %s with %v", node, err)
		}
		defer conn.Close()
		c := pb.NewRendezvousClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if arriving {
			_, err = c.PeerArrival(ctx, peer)
		} else {
			_, err = c.PeerRemoval(ctx, peer)
		}
		if err != nil {
			log.Printf("Error from %s of %v", node, err)
		}
	}
	return err
}

func createPeer(length int, buff []byte) (*pb.Peer, error) {
	peer := new(pb.Peer)
	err := json.Unmarshal(buff[:length], &peer)
	if err != nil {
		fmt.Println("Error in createPeer: " + err.Error())
		return nil, err
	}
	return peer, nil
}

func registerAsNode() {
	conn, err := grpc.Dial(LB_IP, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v", LB_IP, err)
	}
	defer conn.Close()
	c := pb.NewLoadBalancerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	privIp, err := util.PrivateIP()
	if err != nil {
		log.Printf("Error when trying to find private ip: ", err.Error())
	}
	myNode := &pb.Node{IP: privIp, NotifyOthers: DONT_NOTIFY_OTHERS}
	nodeList, err := c.RegisterNode(ctx, myNode)
	if err != nil {
		log.Printf("Error from %s of %v", LB_IP, err)
	}
	for _, Node := range nodeList.GetNodes() {
		nodeMap[Node.IP] = 1
	}
}

// Starts a generic GRPC server
func StartGRPCServer() {
	lis, err := net.Listen("tcp", localGRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRendezvousServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	var port = flag.String("port", "50001", "The port the grpc server will listen on")
	flag.Parse()
	localGRPCAddr = "127.0.0.1:" + *port
	var err error
	StartGRPCServer()
	registerAsNode()
	return
	connection, err := quic.ListenAddr(":", util.GenerateTLSConfig(), nil)
	buff := make([]byte, 1000)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		//Blocks waiting for a connection
		session, err := connection.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		defer session.Close()
		stream, err := session.AcceptStream()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		defer stream.Close()
		len, err := stream.Read(buff)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		peer, err := createPeer(len, buff)
		if err != nil {
			fmt.Println("Error parsing peer info: " + err.Error())
			continue
		} else {
			if _, ok := peerMap[peer.GetFriend()]; ok {
				delete(peerMap, peer.GetFriend())
				notifyPeerUpdate(peerMap[peer.GetFriend()], false)
				fmt.Println("Connecting " + peer.Name + " and " + peer.Friend)
				//TODO, send peer information to both peers
				continue
			}
			peerMap[peer.GetName()] = peer
			notifyPeerUpdate(peer, true)
		}
	}
}
