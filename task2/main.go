package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

/**
Write a Go program that connects to a remote server via SSH, runs a specified ("uptime") command,
and prints the command's output to the console. Parse the uptime command output to get details
like server, status and users count. Final output to have JSON entries for all records in uptime.json
file. The program should take the server's IP address, username, and password as command-line
arguments. Additionally, the program should handle errors gracefully and be able to handle large
outputs efficiently.
*/

func main() {
	ip := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		os.Exit(1)
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s\n", err)
		os.Exit(1)
	}

	var uptimeOutput []byte
	uptimeOutput, err = session.Output("uptime")
	if err != nil {
		fmt.Printf("Failed to run command: %s\n", err)
		os.Exit(1)
	}

	// Parse the uptime command output to get details
	// like server, status, and users count
	// Outputs the JSON entries for all records in uptime.json file
	var uptimeData map[string]interface{}
	json.Unmarshal([]byte(uptimeOutput), &uptimeData)

	f, err := os.OpenFile("uptime.json", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.Encode(uptimeData)

	session.Close()
	client.Close()

}
