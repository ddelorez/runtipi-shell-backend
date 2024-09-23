// cmd/main.go

package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

// ... existing upgrader code ...

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://frontend":             true,
			"http://localhost":            true,
			"http://your-frontend-domain": true,
		}
		return allowedOrigins[origin]
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Start a local shell
	cmd := exec.Command("bash") // Use "cmd" for Windows
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("PTY start error:", err)
		return
	}
	defer func() { _ = ptmx.Close() }() // Best effort

	// Handle input from WebSocket to PTY
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				cmd.Process.Kill()
				break
			}
			_, err = ptmx.Write(msg)
			if err != nil {
				log.Println("PTY write error:", err)
				break
			}
		}
	}()

	// Handle output from PTY to WebSocket
	buf := make([]byte, 1024)
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("PTY read error:", err)
			}
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
