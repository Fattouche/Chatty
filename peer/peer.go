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

// holePunch punches a hole through users NATs if they exist in different networks. test
func holePunch(server *net.UDPConn, addr *net.UDPAddr) error {
	connected := false
	go func() {
		for connected != true {
			server.WriteToUDP([]byte("1"), addr)
			time.Sleep(10 * time.Millisecond)
		}
	}()
	buff := make([]byte, 100)
	server.SetReadDeadline(time.Now().Add(time.Second * 5))
	for {
		_, recvAddr, err := server.ReadFromUDP(buff)
		if err != nil {
			connected = true
			return err
		}
		if recvAddr.String() == addr.String() {
			connected = true
			time.Sleep(time.Millisecond * 500)
			return nil
		}
	}
}

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

// getPeerInfo communicates with the centralized server to exchange information between peers.
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

func connectAsSender(server *net.UDPConn, addr string) (quic.Stream, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	session, err := quic.Dial(server, udpAddr, addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Println("Error: ", err)
		server.Close()
		return nil, err
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		return nil, err
	}
	stream.Write([]byte("\n"))
	return stream, nil
}

func connectAsReciever(server *net.UDPConn, addr string) (quic.Stream, error) {
	server.SetReadDeadline(time.Now().Add(time.Second * 10))
	connection, err := quic.Listen(server, util.GenerateTLSConfig(), nil)
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil, err
	}
	session, err := connection.Accept()
	server.SetReadDeadline(time.Now().Add(time.Hour * 24))
	if err != nil {
		log.Println("Error: " + err.Error())
		server.Close()
		return nil, err
	}

	stream, err := session.AcceptStream()
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil, err
	}
	garbage := make([]byte, 10)
	stream.Read(garbage)
	return stream, nil
}

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

func write(stream quic.Stream) {
	defer stream.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		stream.Write([]byte(text))
	}
}

func chat(server *net.UDPConn, addr string) error {
	var stream quic.Stream
	var err error
	if myPeerInfo.Dialer {
		stream, err = connectAsSender(server, addr)
		if err != nil {
			return err
		}
	} else {
		stream, err = connectAsReciever(server, addr)
		if err != nil {
			return err
		}
	}
	log.Println("Connected to: " + myPeerInfo.Friend + " you can start typing!\n")
	go read(stream)
	write(stream)
	return nil
}

func chatWithPeer(server *net.UDPConn) error {
	var err error
	addr, _ := net.ResolveUDPAddr("udp", friend.Pub_ip)
	laddr, _ := net.ResolveUDPAddr("udp", myPeerInfo.Priv_ip)
	public := true
	err = holePunch(server, addr)
	if err != nil {
		server.Close()
		time.Sleep(time.Millisecond * 500)
		server, _ = net.ListenUDP("udp", laddr)
		public = false
	}
	if public {
		return chat(server, friend.Pub_ip)
	} else {
		return chat(server, friend.Priv_ip)
	}
}

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
