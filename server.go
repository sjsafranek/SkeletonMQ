/*=======================================*/
//	project: gospatial
//	author: stefan safranek
//	email: sjsafranek@gmail.com
/*=======================================*/

package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	// "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	// "time"
)

var (
	port          int
	tcp_port      int
	bind          string
	versionReport bool
)

const (
	// DEFAULT_HTTP_PORT int = 8080
	//DEFAULT_TCP_PORT    int    = 3333
	VERSION string = "0.0.1"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
	flag.IntVar(&port, "p", DEFAULT_HTTP_PORT, "http server port")
	//flag.IntVar(&tcp_port, "tcp_port", DEFAULT_TCP_PORT, "tcp server port")
	flag.BoolVar(&versionReport, "V", false, "App Version")

	flag.Parse()
	if versionReport {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

}

func main() {

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// source: http://patorjk.com/software/taag/#p=display&f=Slant&t=SkeletonMQ
	// HyperCube Platforms
	fmt.Println(`
    _____ __        __     __              __  _______
   / ___// /_____  / /__  / /_____  ____  /  |/  / __ \
   \__ \/ //_/ _ \/ / _ \/ __/ __ \/ __ \/ /|_/ / / / /
  ___/ / ,< /  __/ /  __/ /_/ /_/ / / / / /  / / /_/ /
 /____/_/|_|\___/_/\___/\__/\____/_/ /_/_/  /_/\___\_\

	`)

	// Graceful shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for sig := range sigs {
			// sig is a ^C, handle it
			log.Println("Recieved ", sig)
			log.Println("Gracefully shutting down")
			log.Println("Waiting for sockets to close...")
			for {
				// 	geo_skeleton_server.ServerLogger.Info("Shutting down...")
				os.Exit(0)
			}
		}
	}()

	http_server := HttpServer{Port: port}
	http_server.Start()

	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)
}
