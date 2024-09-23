// cmd/main.go

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade an HTTP connection to a WebSocket connection

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost"
		}
		origins := strings.Split(allowedOrigins, ",")
		for _, o := range origins {
			if origin == o {
				return true
			}
		}
		return false
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Start a local shell - account for OS differences
	shell := os.Getenv("SHELL_COMMAND")
	if shell == "" {
		shell = "bash" //Default Shell
	}
	cmd := exec.Command(shell) // Use "cmd" for Windows
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
