package client

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"github.com/xumc/host-to-host-file-transport/protocol"
	"github.com/xumc/host-to-host-file-transport/util"
)

func Start(serverIP string, serverPort int, filePath string) {
	serverAddr := serverIP + ":" + strconv.Itoa(serverPort)

	file, err := os.Open(filePath)
	util.CheckErr(err)
	defer file.Close()

	conn, err := net.Dial("tcp", serverAddr)
	util.CheckErr(err)
	defer conn.Close()

	fileInfo, err := file.Stat()
	util.CheckErr(err)

	if !fileInfo.IsDir() {
		err = sendFile(conn, fileInfo.Name(), file)
		util.CheckErr(err)
		return
	}

	walk := func(path string, info os.FileInfo, err error) error {
		util.CheckErr(err)
		
		if !info.IsDir() {
			file, err := os.Open(path)
			util.CheckErr(err)
			defer file.Close()

			err = sendFile(conn, path, file)
			util.CheckErr(err)
			return nil
		}
		return nil
	}

	filepath.Walk(filePath, walk)
}

func sendFile(conn net.Conn, filePath string, file *os.File) error {
	fileInfo, err := file.Stat()
	util.CheckErr(err)

	conn.Write([]byte(protocol.ConstHeader))

	contentLen := len([]byte(filePath)) + int(fileInfo.Size())
	conn.Write(protocol.IntToBytes(contentLen))

	filePathBytes := append([]byte(filePath), '\n')
	_, err = conn.Write(filePathBytes)
	if err != nil {
		util.Log(err)
		return err
	}

	for {
		buffer := make([]byte, 1024)
		nr, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			util.Log(err)
			return err
		}
		conn.Write(buffer[:nr])
	}

	return nil
}
