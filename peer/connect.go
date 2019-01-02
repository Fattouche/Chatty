package main

import (
	"crypto/tls"
	"net"
	"time"

	"log"

	util "github.com/Fattouche/Chatty/util"
	quic "github.com/lucas-clemente/quic-go"
)

func connect(server *net.UDPConn, destinationPrivIP, destinationPubIP string, dialer bool) (stream quic.Stream, err error) {
	pubAddr, err := net.ResolveUDPAddr("udp4", destinationPubIP)
	if err != nil {
		log.Println(err)
		return
	}
	privAddr, err := net.ResolveUDPAddr("udp4", destinationPrivIP)
	if err != nil {
		log.Println(err)
		return
	}
	err = holePunch(server, pubAddr)
	if err == nil {
		//can connect publicly
		stream, err = attemptConnection(server, pubAddr, destinationPubIP, dialer)
	} else {
		//should connect privately
		addr, _ := net.ResolveUDPAddr("udp4", server.LocalAddr().String())
		server.Close()
		time.Sleep(time.Millisecond * 500)
		server, _ = net.ListenUDP("udp", addr)
		stream, err = attemptConnection(server, privAddr, destinationPrivIP, dialer)
	}
	return
}

func attemptConnection(server *net.UDPConn, udpAddr *net.UDPAddr, addr string, dialer bool) (stream quic.Stream, err error) {
	if dialer {
		stream, err = connectAsSender(server, udpAddr, addr)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		stream, err = connectAsReciever(server)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

// holePunch punches a hole through users NATs if they exist in different networks
func holePunch(server *net.UDPConn, addr *net.UDPAddr) (err error) {
	connected := false
	var recvAddr *net.UDPAddr
	// Try to send udp packets
	go func() {
		for connected != true {
			server.WriteToUDP([]byte("1"), addr)
			time.Sleep(10 * time.Millisecond)
		}
	}()
	buff := make([]byte, 100)
	// If you don't recieve a udp packet in 2 seconds return
	server.SetReadDeadline(time.Now().Add(time.Second * 2))
	for {
		// Listen for udp packets
		_, recvAddr, err = server.ReadFromUDP(buff)
		if err != nil {
			connected = true
			return
		}
		if recvAddr.String() == addr.String() {
			connected = true
			time.Sleep(time.Millisecond * 500)
			return
		}
	}
}

// This peer is the reciever so listen on quic for the other peer dialing
func connectAsReciever(server *net.UDPConn) (stream quic.Stream, err error) {
	server.SetReadDeadline(time.Now().Add(time.Second * 10))
	connection, err := quic.Listen(server, util.GenerateTLSConfig(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	session, err := connection.Accept()
	server.SetReadDeadline(time.Now().Add(time.Hour * 24))
	if err != nil {
		log.Println(err)
		server.Close()
		return
	}

	stream, err = session.AcceptStream()
	if err != nil {
		log.Println(err)
		return
	}
	garbage := make([]byte, 10)
	stream.Read(garbage)
	return
}

// This peer is the sender so quic dial to the reciever(other peer)
func connectAsSender(server *net.UDPConn, udpAddr *net.UDPAddr, addr string) (stream quic.Stream, err error) {
	session, err := quic.Dial(server, udpAddr, addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Println(err)
		server.Close()
		return
	}
	stream, err = session.OpenStreamSync()
	if err != nil {
		return
	}
	stream.Write([]byte("\n"))
	return
}
