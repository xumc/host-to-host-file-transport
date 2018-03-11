package server

import (
	"strings"
	"bytes"
	"io"
	"net"
	"os"
	"path/filepath"
	"github.com/xumc/host-to-host-file-transport/protocol"
	"github.com/xumc/host-to-host-file-transport/util"
	"regexp"
)

const (
	ConnType = "tcp"
	Addr = ":0"
)

func Start() {
	potentialIPs := getIPs()

	l, err := net.Listen(ConnType, Addr)
	if err != nil {
		util.Log("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	port := strings.Replace(l.Addr().String(), "[::]:", "", -1)
	util.Log("potential IPs: ", strings.Join(potentialIPs, ", "), "port: ", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			util.Log("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()

	var file *os.File
	defer file.Close()

	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			util.Log(conn.RemoteAddr().String(), "disconnected")
			return
		}

		if err != nil {
			util.Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			i := bytes.IndexByte(data, '\n')
			if i == -1 {
				util.Log("data format error")
				continue
			}

			// Todo remove server_ prefix, we add prefix for testing convenience
			filePath := "server_" + string(data[:i])
			ensureFileDirExist(filePath)

			file, err := os.Create(filePath)
			util.CheckErr(err)
			defer file.Close()

			_, err = file.Write(data[(i + 1):])
			util.CheckErr(err)
		}
	}
}

func ensureFileDirExist(filePath string) error {
	dir := filepath.Dir(filePath)
	exist, err := isExists(dir)
	util.CheckErr(err)

	if !exist {
		err = os.MkdirAll(dir, 0777)
		util.CheckErr(err)
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

func getIPs() []string {
	ifaces, err := net.Interfaces()
	util.CheckErr(err)

	ips := make([]string, 0)
	for _, i := range ifaces {  
        addrs, _ := i.Addrs()  
		util.CheckErr(err)
  
        for _, addr := range addrs {
			if addr.Network() == "ip+net" {
				match, err := regexp.MatchString("", addr.String())
				util.CheckErr(err)
				if match {
					ip := strings.Split(addr.String(), "/")[0]
					if !strings.Contains(ip, ":") {
						ips = append(ips, ip)
					}
				}
			}
        }  
	}
	return ips
}
