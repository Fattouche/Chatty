package main

import (
	"log"
	"net"
	"sync"
	"testing"
)

func initializeServer(ipAddr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", ipAddr)
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("Error: " + err.Error())
		server.Close()
		return nil, err
	}
	return server, nil
}

func TestConnect(t *testing.T) {
	peer1IP := "127.0.0.1:8080"
	peer2IP := "127.0.0.1:8081"
	pubIP := ":8082"
	testStr := "testing"
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		server, err := initializeServer(peer1IP)
		if err != nil {
			log.Println("Error: " + err.Error())
			server.Close()
			t.Error(err)
		}
		stream, err := connect(server, peer2IP, pubIP, true)
		if err != nil {
			t.Error(err)
			wg.Done()
			return
		}
		buf := make([]byte, 100)
		stream.Write([]byte(testStr))
		n, _ := stream.Read(buf)
		if string(buf[:n]) != testStr {
			t.Fail()
		}
		wg.Done()
	}()

	go func() {
		server, err := initializeServer(peer2IP)
		if err != nil {
			log.Println("Error: " + err.Error())
			server.Close()
			t.Error(err)
		}
		stream, err := connect(server, peer1IP, pubIP, false)
		if err != nil {
			t.Error(err)
			wg.Done()
			return
		}
		buf := make([]byte, 100)
		stream.Write([]byte(testStr))
		n, _ := stream.Read(buf)
		if string(buf[:n]) != testStr {
			t.Fail()
		}
		wg.Done()
	}()
	wg.Wait()
}
