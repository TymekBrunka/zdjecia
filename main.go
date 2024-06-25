package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"zdj/webs"
)

func GetIPV4() net.IP {
	var (
		ret    net.IP
		err    error
		ifaces []net.Interface
		addrs  []net.Addr
	)
	if ifaces, err = net.Interfaces(); err == nil {
		for _, i := range ifaces {
			if addrs, err = i.Addrs(); err == nil {
				for _, a := range addrs {
					if ipnet, ok := a.(*net.IPNet); ok {
						if ipv4 := ipnet.IP.To4(); ipv4 != nil {
							if ipv4.IsGlobalUnicast() {
								ret = ipv4
							}
						}
					}
				}
			}
		}
	}
	return ret
}

func retIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(GetIPV4().String()))
}

func main() {

	fmt.Println("Clearing cache folders")
	os.RemoveAll("./zasoby_serwera/u_0/")
	os.Mkdir("./zasoby_serwera/u_0/", 0750)
	os.RemoveAll("./zasoby_serwera/u_1/")
	os.Mkdir("./zasoby_serwera/u_1/", 0750)
	fmt.Println("Starting server at port 8000")
	fmt.Printf("-> Hosting at \x1b[32mhttp://%s:8000\x1b[0m <-\n", GetIPV4())

	fileServer := http.FileServer(http.Dir("./static")) // New code
	http.Handle("/", fileServer)
	http.Handle("/ws", websocket.Handler(webs.HandleWebsocket))

	http.HandleFunc("/sendfiles", webs.SendFilesToClient)
	http.HandleFunc("/givfilesplz", webs.UcanHasCheeseburger)
	http.HandleFunc("/ipv4", retIP)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
