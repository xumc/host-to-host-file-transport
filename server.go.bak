package main

import (
    "fmt"
    "net"
	"os"
	"io"
	"bytes"
	"path/filepath"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {
	var (
		ip string
		host int
	)

	serverTip := fmt.Sprintf("server started: ip: %s host: %d", ip, host)
	fmt.Println(serverTip)

    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

	defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()

	filePathBytes := make([]byte, 0)
	fileNameFound := false

	var file *os.File	
	defer file.Close()

	for {
		buffer := make([]byte, 1024)
		nr, err := conn.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		i := bytes.IndexByte(buffer, '\n')
		if !fileNameFound {
			if i == -1 {
				fileNameFound = false
				filePathBytes = append(filePathBytes, buffer[:nr]...)
				continue
			}

			fileNameFound = true
			filePathBytes = append(filePathBytes, buffer[:i]...)

			// Todo remove server_ prefix, we add prefix for testing convenience
			filePath := "server_" + string(filePathBytes)
			ensureFileDirExist(filePath)
			file, err = os.Create(filePath)
			checkErr(err)

			_, err = file.Write(buffer[(i + 1):nr])
			checkErr(err)
			continue	
		}

		_, err = file.Write(buffer[:nr])
		checkErr(err)
	}
}

func ensureFileDirExist(filePath string) error {
	dir := filepath.Dir(filePath)
	exist, err := isExists(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !exist {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return nil
}

func isExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}