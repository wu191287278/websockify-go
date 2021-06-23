package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	httpAddress := flag.String("address", ":8888", "Server port")
	targetAddress := flag.String("target", "localhost:5900", "Target port")
	flag.Parse()
	http.Handle("/websockify", websocket.Handler(func(wsconn *websocket.Conn) {
		defer wsconn.Close()
		var d net.Dialer
		var address = *targetAddress
		conn, err := d.DialContext(wsconn.Request().Context(), "tcp", address)
		if err != nil {
			log.Printf("[%s] [ERROR] [%v]", address, err)
			return
		}
		defer conn.Close()
		wsconn.PayloadType = websocket.BinaryFrame
		go func() {
			io.Copy(wsconn, conn)
			wsconn.Close()
			log.Printf("[%s] [SESSION_CLOSED]", address)
		}()
		io.Copy(conn, wsconn)
		log.Printf("[%s] [CLIENT_DISCONNECTED]", address)
	}))
	log.Printf("Http listening on %s \n", *httpAddress)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))

}
