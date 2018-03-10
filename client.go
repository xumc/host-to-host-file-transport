package main

import (
	"fmt"
	"os"
	"net"
	"io"
	"strconv"
	"flag"
	"path/filepath"
)

func sendFile(serverAddr string, filePath string, file *os.File) error {
	conn, err := net.Dial("tcp", serverAddr)
	checkErr(err)
	defer conn.Close()

	filePathBytes := append([]byte(filePath), '\n')
	_, err = conn.Write(filePathBytes)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		buffer := make([]byte, 1024)
		nr, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		conn.Write(buffer[:nr])
	}

	return nil
}

func main() {
	serverIP := *flag.String("ip", "127.0.0.1", "remote IP")
	serverHost := *flag.Int("host", 3333, "remote host")
	filePath := os.Args[3]

	fmt.Printf("serverIP: %s serverHost: %d path: %s \n", serverIP, serverHost, filePath)

	serverAddr := serverIP + ":" + strconv.Itoa(serverHost)

	file, err := os.Open(filePath)
	defer file.Close()
	checkErr(err)

	fileInfo, err := file.Stat()
	checkErr(err)

	if !fileInfo.IsDir() {
		err = sendFile(serverAddr, fileInfo.Name(), file)
		checkErr(err)
		return
	}

	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fmt.Println(path)
			file, err := os.Open(path)
			checkErr(err)
			defer file.Close()
			
			err = sendFile(serverAddr, path, file)
			checkErr(err)
			return nil
		}
		return nil
	}

	filepath.Walk(filePath, walk)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}