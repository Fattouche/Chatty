package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/Fattouche/Chatty/protobuff"
	util "github.com/Fattouche/Chatty/util"
	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/grpc"
)

const (
	NOTIFY_OTHERS      = 1
	DONT_NOTIFY_OTHERS = 2
	PORT               = "50003"
	LB_IP              = "127.0.0.1:5004"
)

type server struct{}

var nodeMap = make(map[string]int)
var peerMap map[string]*pb.Peer

func (s *server) NodeDown(ctx context.Context, nodes *pb.Node) (*pb.Response, error) {
	var status string
	for _, node := range nodes.GetNodes() {
		delete(nodeMap, node)
	}
	if nodes.GetNotifyOthers() == 1 {
		err := notifyDeadNode(nodes.GetNodes())
		if err != nil {
			status = err.Error()
		} else {
			status = "ok"
		}
	}
	return &pb.Response{Status: status}, nil
}

func notifyDeadNode(deadNodes []string) error {
	var err error
	for node := range nodeMap {
		conn, err := grpc.Dial(node, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to %s with %v", node, err)
		}
		defer conn.Close()
		c := pb.NewServerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.NodeDown(ctx, &pb.Node{Nodes: deadNodes, NotifyOthers: DONT_NOTIFY_OTHERS})
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
		c := pb.NewServerClient(conn)
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

func (s *server) PeerArrival(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	peerMap[peer.GetName()] = peer
	return &pb.Response{Status: "ok"}, nil
}

func (s *server) PeerRemoval(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	delete(peerMap, peer.GetName())
	return &pb.Response{Status: "ok"}, nil
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
	c := pb.NewServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	nodeList, err = c.registerAsNode(util.privateIP())
	if err != nil {
		log.Printf("Error from %s of %v", node, err)
	}
	for _, node := range nodeList {
		nodeMap[node] = 1
	}
}

func main() {
	var err error
	registerAsNode()
	connection, err := quic.ListenAddr(":"+PORT, util.generateTLSConfig(), nil)
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
