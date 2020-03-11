package main

import (
	"fmt"
	"gsocket/websocket"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {

	pool := websocket.NewPool(1)
	go pool.Start()

	pool2 := websocket.NewPool(2)
	go pool2.Start()

	p := make([]*websocket.Pool, 0)
	p = append(p, pool)
	p = append(p, pool2)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	router.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		index := 0
		for indx, pp := range p {
			if pp.ID == id {
				index = indx
			}
		}
		serveWs(p[index], w, r)
	})

	http.ListenAndServe(":8080", router)
}

func main() {
	setupRoutes()

}
