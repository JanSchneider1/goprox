package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	app = cli.NewApp()
	defaultPort = 8080
	defaultMessage = "Hello World"
)

func info() {
	app.Name = "goprox"
	app.Usage = "Simple proxy server made with Go"
	app.Version = "1.0.0"
}

func startServer(port int, message string) {
	log.Println("Starting server...")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		mapping := filepath.Clean(request.URL.Path)
		log.Printf("Client requested resource! (URL: localhost:%d%v)", port, mapping)
		_, _ = writer.Write([]byte(message))
	})
	log.Printf("Server listens now on localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}
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
			Action:    commandStartServer,
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
