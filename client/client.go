package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"path/filepath"
	"strconv"
	"github.com/xumc/host-to-host-file-transport/protocol"
	"github.com/xumc/host-to-host-file-transport/util"
)

type upload struct {
	filePath string
	total int
	read int
}

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

	progress := make(chan upload, 10)
	go update(progress)


	if !fileInfo.IsDir() {
		err = sendFile(conn, fileInfo.Name(), file, progress)
		util.CheckErr(err)
		return
	}

	walk := func(path string, info os.FileInfo, err error) error {
		util.CheckErr(err)
		
		if !info.IsDir() {
			file, err := os.Open(path)
			util.CheckErr(err)
			defer file.Close()

			err = sendFile(conn, path, file, progress)
			util.CheckErr(err)
			return nil
		}
		return nil
	}

	filepath.Walk(filePath, walk)
}

func sendFile(conn net.Conn, filePath string, file *os.File, progress chan upload) error {
	fileInfo, err := file.Stat()
	util.CheckErr(err)
	fileSize := int(fileInfo.Size())

	conn.Write([]byte(protocol.ConstHeader))

	contentLen := len([]byte(filePath)) + fileSize
	conn.Write(protocol.IntToBytes(contentLen))

	filePathBytes := append([]byte(filePath), '\n')
	_, err = conn.Write(filePathBytes)
	if err != nil {
		util.Log(err)
		return err
	}

	var totalRead int = 0
	up := upload{filePath: filePath, total: fileSize,}

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

		nr, err = conn.Write(buffer[:nr])
		if err != nil {
			util.Log(err)
			return err
		}

		totalRead += nr
		up.read = totalRead
		progress <- up

	}
	return nil
}

func update(progress chan upload) {
	ticker := time.NewTicker(500 * time.Millisecond)
	up := upload{}

	for {
		select {
		case up = <- progress:
		case <-ticker.C:
			fmt.Printf("\r%s => %d/%d", up.filePath, up.read, up.total)
		}
	}
}
