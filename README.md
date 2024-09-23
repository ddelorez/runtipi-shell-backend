# Web-Based Shell Interface Setup Guide
## This guide walks you through setting up a web-based shell interface for your home server. We'll cover everything from creating the project structure to building and running the application using Go and Docker.

## Prerequisites
- Basic Knowledge: Familiarity with command-line operations.
- Installed Software:
    - Go (version 1.18 or later)
    - Docker

## Project Overview
We'll create a Go application that:
- Serves a WebSocket endpoint.
- Starts a shell session on the server.
- Relays input and output between the shell and the client.

## Directory Structure
```
go
web-shell/
├── cmd/
│   └── main.go
├── Dockerfile
├── go.mod
└── go.sum
```

## Step 1: Set Up the Project Directory
### 1.1 Create the Project Directory
Open your terminal and run:
```bash
mkdir web-shell
cd web-shell
```

### 1.2 Initialize the Go Module
```bash
go mod init github.com/yourusername/web-shell
```

## Step 2: Create the Go Application
### 2.1 Create the cmd Directory and main.go File
```bash
mkdir cmd
touch cmd/main.go
```

### 2.2 Write the Code in main.go
Open `cmd/main.go` in a text editor and add the following code:
```go
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

// Rest of the code...
```

## Step 3: Install Dependencies
In your terminal, run:
```bash
go get github.com/gorilla/websocket
go get github.com/creack/pty
```
This will create `go.mod` and `go.sum` files.

## Step 4: Create the Dockerfile
### 4.1 Create the Dockerfile
```bash
touch Dockerfile
```

### 4.2 Write the Dockerfile Content
Open `Dockerfile` and add:
```dockerfile
# Build Stage
FROM golang:1.18-alpine AS builder

# Rest of the Dockerfile content...
```

## Step 5: Build and Run the Application
### 5.1 Build the Docker Image
```bash
docker build -t web-shell .
```

### 5.2 Run the Docker Container
```bash
docker run -d -p 8080:8080 \
    -e ALLOWED_ORIGINS="http://localhost" \
    -e SHELL_COMMAND="bash" \
    --name web-shell-container \
    web-shell
```

## Step 6: Test the Application
### 6.1 Using a WebSocket Client
Since we don't have a frontend yet, you can test using a WebSocket client.

Option 1: Browser Console
Open your browser's developer console (usually by pressing F12) and run:
```javascript
const socket = new WebSocket('ws://localhost:8080/ws');
socket.onopen = () => {
    console.log('WebSocket connection opened');
    // Send a command to the shell
    socket.send('ls\n');
};
socket.onmessage = (event) => {
    console.log('Received from server:', event.data);
};
socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};
socket.onclose = () => {
    console.log('WebSocket connection closed');
};
```

Option 2: WebSocket Client Tools
Use tools like WebSocket King or Postman (v9 or later) to connect to `ws://localhost:8080/ws`.

### 6.2 Verify the Output
You should see the output of the `ls` command in your console or client.
Check the server logs for any errors.

## Step 7: Clean Up
When you're done testing, you can stop and remove the Docker container:
```bash
docker stop web-shell-container
docker rm web-shell-container
```

## Recap of Files
- `cmd/main.go`: See Step 2.2 for the full code.
- `Dockerfile`: See Step 4.2 for the full content.
- `go.mod` and `go.sum`: These files are generated and managed by Go modules.

## Next Steps
- Integrate a Frontend: Use a terminal emulator like xterm.js to build a user interface.
- Security Enhancements: Implement authentication and secure connections (HTTPS/WSS).
- Configuration: Add more environment variables for flexibility.

## Tips and Notes
### Environment Variables:
- `ALLOWED_ORIGINS`: Comma-separated list of allowed origins.
- `SHELL_COMMAND`: The shell to execute (e.g., bash, sh).

### Docker Commands:
- Build Image: `docker build -t web-shell .`
- Run Container: See Step 5.2.

### Testing:
- Ensure the port 8080 is not blocked by a firewall.
- If you encounter issues, check the server logs for errors.
