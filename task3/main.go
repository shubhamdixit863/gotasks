package task3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
	"sync"
)

type Configuration struct {
	Username string
	Password string
	Ip       string
}

func readFile(filename string) ([]Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var configurations []Configuration
	for _, line := range lines {
		fields := strings.Fields(line)
		configurations = append(configurations, Configuration{
			Username: fields[0],
			Password: fields[1],
			Ip:       fields[2],
		})
	}

	return configurations, nil
}

// This login function will be running using goroutines workerpools

func LoginOverSsh(username string, password string, ip string) {
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

func main() {

	// Create a worker pool with a size equal to the number of lines
	lines, err := readFile("ssh.config")
	if err != nil {
		log.Println("Error Reading file")
	}
	numWorkers := len(lines) // We can spawn thousands of go routines but they will be managed by go scheduler
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			username := lines[i].Username
			password := lines[i].Password
			ip := lines[i].Ip
			LoginOverSsh(username, password, ip)
		}(i)
	}

	wg.Wait() // wait for all the workers to complete their job

}
