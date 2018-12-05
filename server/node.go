package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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
	GRPC_PORT          = ":50000"
	QUIC_PORT          = ":50001"
)

// This server implements the protobuff Node type
type server struct{}

var gcprAddr string
var lbAddr string
var quicAddr string

// nodeMap keeps track of all the rendezvous servers(nodes)
var nodeMap = make(map[string]int)

// peerMap keeps track of all the peer information
var peerMap map[string]*pb.Peer

// Removes a node from our nodeMap and notifies the other nodes if neccesary
func (s *server) NodeRemoval(ctx context.Context, node *pb.Node) (*pb.Response, error) {
	var status string
	delete(nodeMap, node.Grpc_IP)
	log.Println("Node removed: ", node.Grpc_IP)
	if node.GetNotifyOthers() == 1 {
		err := notifyNodeUpdate(node, false)
		if err != nil {
			status = err.Error()
		} else {
			status = "ok"
		}
	}
	return &pb.Response{Status: status}, nil
}

// Adds a new node to the nodeMap and notifies other nodes if necessary
func (s *server) NodeArrival(ctx context.Context, node *pb.Node) (*pb.Response, error) {
	nodeMap[node.Grpc_IP] = 1
	var status string
	log.Println("New Node arrived: ", node.Grpc_IP)
	if node.GetNotifyOthers() == 1 {
		err := notifyNodeUpdate(node, true)
		if err != nil {
			status = err.Error()
		} else {
			status = "ok"
		}
	}
	return &pb.Response{Status: status}, nil
}

// A new peer has arrived
func (s *server) PeerArrival(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	log.Println("New peer has arrived: ", peer.GetName())
	peerMap[peer.GetName()] = peer
	return &pb.Response{Status: "ok"}, nil
}

// A peer can be deleted
func (s *server) PeerRemoval(ctx context.Context, peer *pb.Peer) (*pb.Response, error) {
	log.Println("Peer needs to be removed: ", peer.GetName())
	delete(peerMap, peer.GetName())
	return &pb.Response{Status: "ok"}, nil
}

// Healthcheck to satisfy the loadbalancers checks
func (s *server) HealthCheck(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Status: "ok"}, nil
}

// A general function to grpc dial to other nodes within the cluster to tell them of a node arrival or removal
func notifyNodeUpdate(nodeToUpdate *pb.Node, arrival bool) error {
	var err error
	nodeToUpdate.NotifyOthers = DONT_NOTIFY_OTHERS
	for addr := range nodeMap {
		if addr == gcprAddr || nodeToUpdate.Grpc_IP == addr {
			continue
		}
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to %s with %v", addr, err)
		}
		defer conn.Close()
		c := pb.NewRendezvousClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if arrival {
			_, err = c.NodeArrival(ctx, nodeToUpdate)
		} else {
			_, err = c.NodeRemoval(ctx, nodeToUpdate)
		}
		if err != nil {
			log.Printf("Error from %s of %v", addr, err)
		}
	}
	//only returns the last error unfortunately
	return err
}

// A general function to grpc dial to other nodes within the cluster to tell them of a peer arrival or removal
func notifyPeerUpdate(peer *pb.Peer, arriving bool) error {
	var err error
	for addr := range nodeMap {
		if addr == gcprAddr {
			continue
		}
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to %s with %v", addr, err)
		}
		defer conn.Close()
		c := pb.NewRendezvousClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if arriving {
			peerMap[peer.GetName()] = peer
			_, err = c.PeerArrival(ctx, peer)
		} else {
			delete(peerMap, peer.Name)
			_, err = c.PeerRemoval(ctx, peer)
		}
		if err != nil {
			log.Printf("Error from %s of %v", addr, err)
		}
	}
	return err
}

// Creates a new peer
func createPeer(length int, buff []byte, pubIP string) (*pb.Peer, error) {
	peer := new(pb.Peer)
	err := json.Unmarshal(buff[:length], &peer)
	if err != nil {
		fmt.Println("Error in createPeer: " + err.Error())
		return nil, err
	}
	peer.PubIp = pubIP
	return peer, nil
}

// Registeres the current node as a node with the loadbalancer, this information will be broadcasted out
func registerAsNode() {
	conn, err := grpc.Dial(lbAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v", lbAddr, err)
	}
	defer conn.Close()
	c := pb.NewLoadBalancerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	myNode := &pb.Node{Grpc_IP: gcprAddr, Pub_IP: quicAddr, NotifyOthers: DONT_NOTIFY_OTHERS}
	fmt.Println("Dialining to :", lbAddr)
	nodeList, err := c.RegisterNode(ctx, myNode)
	if err != nil {
		log.Printf("Error from %s of %v", lbAddr, err)
		os.Exit(1)
	}
	log.Println("Registered as node, recieved current list: ", nodeList.Nodes)

	for _, Node := range nodeList.GetNodes() {
		nodeMap[Node.Grpc_IP] = 1
	}
}

// Starts a generic GRPC server
func StartGRPCServer() {
	lis, err := net.Listen("tcp", GRPC_PORT)
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

func watchForFriend(stream quic.Stream, session quic.Session, peer *pb.Peer, dialer bool) {
	for {
		if _, ok := peerMap[peer.GetFriend()]; ok {
			peerFriend := peerMap[peer.Friend]
			peerFriend.Dialer = dialer
			msgForPeer, err := json.Marshal(peerFriend)
			if err != nil {
				fmt.Println("Error marshalling in checkpeer: " + err.Error())
			}
			log.Printf("Notifying %s of %s info", peer.Name, peer.Friend)
			stream.Write(msgForPeer)

			time.Sleep(time.Millisecond * 1000)
			notifyPeerUpdate(peer, false)
			stream.Close()
			session.Close()
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	var pub_ip = flag.String("public", "68.183.175.129", "The ip the quic server will listen on")
	var lb_ip = flag.String("lb", "68.183.175.69", "The ip of the loadbalancer")
	flag.Parse()
	lbAddr = *lb_ip + GRPC_PORT
	gcprAddr = *pub_ip + GRPC_PORT
	quicAddr = *pub_ip + QUIC_PORT
	var err error
	go StartGRPCServer()
	time.Sleep(time.Microsecond * 5000)
	registerAsNode()
	// Listen for peer connections
	addr, err := net.ResolveUDPAddr("udp4", QUIC_PORT)
	server, err := net.ListenUDP("udp4", addr)
	connection, _ := quic.Listen(server, util.GenerateTLSConfig(), nil)
	log.Println("Listening for quic connections on :", connection.Addr().String())
	buff := make([]byte, 1000)
	if err != nil {
		log.Fatal(err)
		return
	}
	peerMap = make(map[string]*pb.Peer)
	for {
		//Blocks waiting for a connection
		session, err := connection.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		defer session.Close()
		pubIP := session.RemoteAddr()
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
		// Try to create a peer from the quic connection
		peer, err := createPeer(len, buff, pubIP.String())
		fmt.Println("Got a peer connection from: ", peer.Name)
		if err != nil {
			fmt.Println("Error parsing peer info: " + err.Error())
			continue
		} else {
			notifyPeerUpdate(peer, true)
			dialer := true
			if _, ok := peerMap[peer.GetFriend()]; !ok {
				dialer = false
			}
			go watchForFriend(stream, session, peer, dialer)
		}
	}
}
