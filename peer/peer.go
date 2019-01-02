package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/Fattouche/Chatty/protobuff"
	util "github.com/Fattouche/Chatty/util"
	quic "github.com/lucas-clemente/quic-go"
	"google.golang.org/grpc"
)

// Peer used to keep track of peer information.
type Peer struct {
	Priv_ip string
	Pub_ip  string
	Name    string
	Friend  string
	Dialer  bool
}

var myPeerInfo *Peer
var friend Peer
var lbIP string

const (
	BUFFERSIZE = 48000
)

// Dials to the loadbalancer and gets an IP address of one of the live rendezvous servers
func grpcRendezvousAddr() string {
	conn, err := grpc.Dial(lbIP, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to dial to %s with %v\n", lbIP, err)
	}
	defer conn.Close()
	// Create a client that can access the node functions
	c := pb.NewLoadBalancerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	node, err := c.RendevouszServerIP(ctx, &pb.Request{})
	if err != nil {
		log.Println("Error: ", err.Error())
		os.Exit(1)
	}
	return node.Pub_IP
}

// quic dials to the assigned rendezvous server and exchanges peer information
func getPeerInfo(server *net.UDPConn) error {
	CentServerAddr := grpcRendezvousAddr()
	buff, err := json.Marshal(myPeerInfo)
	if err != nil {
		log.Println("Error:" + err.Error())
		return err
	}
	centUDPAddr, err := net.ResolveUDPAddr("udp", CentServerAddr)
	if err != nil {
		log.Println("Error:" + err.Error())
		return err
	}
	session, err := quic.Dial(server, centUDPAddr, CentServerAddr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Println("Error:" + err.Error())
		return err
	}
	defer session.Close()
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Println("Error:" + err.Error())
		return err
	}
	defer stream.Close()
	stream.Write(buff)
	recvBuff := make([]byte, BUFFERSIZE)
	len, _ := stream.Read(recvBuff)
	log.Println("Attempting to connect to " + myPeerInfo.Friend)
	err = json.Unmarshal(recvBuff[:len], &friend)
	if friend.Dialer {
		myPeerInfo.Dialer = false
	} else {
		myPeerInfo.Dialer = true
	}
	if err != nil {
		log.Println("Error:" + err.Error())
		return err
	}

	return nil
}

// reads from the quick stream and outputs it to the console with their name
func read(stream quic.Stream) {
	defer stream.Close()
	buff := make([]byte, 1000)
	for {
		len, err := stream.Read(buff)
		if err != nil {
			log.Println("Error: " + err.Error())
		}
		fmt.Print(friend.Name + ": " + string(buff[:len]))
	}
}

// reads from input and writes to the quic stream
func write(stream quic.Stream) {
	defer stream.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		stream.Write([]byte(text))
	}
}

// Connects to a peer and starts reading/writing
func chatWithPeer(server *net.UDPConn) error {
	stream, err := connect(server, friend.Priv_ip, friend.Pub_ip, myPeerInfo.Dialer)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Connected to: " + myPeerInfo.Friend + " you can start typing!\n")
	go read(stream)
	write(stream)
	return nil
}

// takes in myName and friend which is the only thing that is needed to connect to the friend.
func main() {
	var friendName = flag.String("friend", "Empty", "The name of the peer you wish to talk to")
	var name = flag.String("myName", "Alex", "My name that peer will connect to")
	var lb_addr = flag.String("lb", "68.183.175.69", "The ip  of the loadbalancer")
	flag.Parse()
	lbIP = *lb_addr + ":50000"
	myPeerInfo = new(Peer)
	myPeerInfo.Name = strings.ToLower(*name)
	myPeerInfo.Friend = strings.ToLower(*friendName)
	machineIP, _ := util.PrivateIP()
	addr, err := net.ResolveUDPAddr("udp", ":0")
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("Error: " + err.Error())
		server.Close()
		return
	}
	defer server.Close()
	port := strings.Split(server.LocalAddr().String(), ":")
	myPeerInfo.Priv_ip = machineIP + ":" + port[len(port)-1]
	log.Println("Listening on: " + myPeerInfo.Priv_ip)
	err = getPeerInfo(server)
	if err != nil {
		log.Println("Error :" + err.Error())
		return
	}
	chatWithPeer(server)
}
