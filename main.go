package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/TheLazarusNetwork/TunnelClient/client"
	"github.com/TheLazarusNetwork/TunnelClient/cmd"
	"github.com/TheLazarusNetwork/TunnelClient/core"
	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

func main() {
	// Load enviornment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error in reading the config file: %v", err)
	}

	// get command line argument values
	port, host := cmd.Command()

	// local service to be forwarded
	portInt, _ := strconv.Atoi(port)
	localEndpoint := Endpoint{
		Host: host,
		Port: portInt,
	}

	// remote SSH server
	serverEndpoint := Endpoint{
		Host: os.Getenv("NETWORK_DOMAIN"),
		Port: 22,
	}

	// get user api key
	var apiKey string
	fmt.Print("Enter API Key Provided: ")
	fmt.Scanf("%s", &apiKey)

	// check validity of user
	status := core.CheckUser(apiKey)
	if status == "invalid" {
		fmt.Println("API Key invalid")
		return
	} else {
		fmt.Println("Welcome " + os.Getenv("USER_NAME"))
	}

	// display the tunnels available
	fmt.Println("Available Tunnels are: ")
	fmt.Println(os.Getenv("TUNNELS"))

	// choose tunnel to be used
	var name string
	fmt.Print("Enter from above Tunnel Name you want to used: ")
	fmt.Scanf("%s", &name)

	// remote forwarding port (on remote SSH server network)
	tunnelPort, _ := strconv.Atoi(os.Getenv("TUNNEL_PORT"))
	remoteEndpoint := Endpoint{
		Host: "localhost",
		Port: tunnelPort,
	}

	// create a ssh configuration
	sshConfig := client.Conn()

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, local)
	}

}
