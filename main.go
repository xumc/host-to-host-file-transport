package main

import (
	"flag"
	"os"
	"github.com/xumc/host-to-host-file-transport/client"
	"github.com/xumc/host-to-host-file-transport/server"
	"github.com/xumc/host-to-host-file-transport/util"
)

func main() {
	mode := flag.String("mode", "s", "mode")
	serverIP := flag.String("ip", "127.0.0.1", "remote IP")
	serverPort := flag.Int("port", 3333, "remote host")
	flag.Parse()
	filePath := os.Args[(len(os.Args) -1)]

	switch *mode {
	case "s":
		server.Start()
	case "c":
		client.Start(*serverIP, *serverPort, filePath)
	default:
		util.Log("invaldate args")
	}
}