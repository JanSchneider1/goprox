package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	app            = cli.NewApp()
	defaultPort    = 8080
	defaultMessage = "Hello World"
)

func info() {
	app.Name = "goprox"
	app.Usage = "Simple proxy server made with Go"
	app.Version = "1.0.0"
}

func startServer(port int, message string) {
	ip := getLocalIP()
	log.Println("Starting server...")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		mapping := filepath.Clean(request.URL.Path)
		log.Printf("Client requested resource! (URL: %v:%v%v)", ip, port, mapping)
		_, _ = writer.Write([]byte(message))
	})
	log.Printf("Server listens now on %v:%v\n", ip, port)
	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", ip, port), nil); err != nil {
		panic(err)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func commandStartServer(c *cli.Context) error {
	if c.NArg() > 0 {
		return cli.Exit("usage: goprox start [<Options>]", 87)
	}
	port := c.Int("port")
	message := c.String("message")
	startServer(port, message)
	return nil
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name:      "start",
			Aliases:   []string{"s"},
			Usage:     fmt.Sprintf("Start server on given port (default: %v), echoing given message (default: '%v'", defaultPort, defaultMessage),
			ArgsUsage: "",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort},
				&cli.StringFlag{Name: "message", Aliases: []string{"m", "echo", "e"}, Value: defaultMessage},
			},
			Action: commandStartServer,
		},
	}
}

func main() {
	info()
	commands()
	// Run CLI App and catch errors
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
