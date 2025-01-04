# WebSocket Client Example (WebAssembly)

This example demonstrates how to connect to a WebSocket signaling server and interact with it using WebAssembly (WASM) in a Go client.

## Features

- Connect to a WebSocket signaling server.
- Send/receive signaling messages (Get peer IDs, send messages to peers, etc.).
- WebAssembly support to run in the browser.

## Prerequisites

- Go 1.16+ installed.
- Web browser that supports WebAssembly.
- A WebSocket signaling server (e.g., from this project) running.

## Setup and Usage

### Step 1: Run the Signaling Server

1. Clone the signaling server repo:
   
   ```bash
   git clone github.com/AbdelrahmanWM/signalingserver/signalingserver
   ```

2. Navigate to `examples/basic_wasm_client/`.

2. Install dependencies and run the server:

   ```bash
   go run main.go
   ```

   This will start the signaling server at `ws://localhost:8090/signalingserver`.

### Step 2: Set Up the Frontend

1. Navigate to `examples/basic_wasm_client/wasm`.

2. Build the WASM client:

   ```bash
   GOARCH=wasm GOOS=js go build -o main.wasm main.go
   ```

   This will generate `main.wasm`.

3. Open `demo.html` in your browser (ensure `main.wasm` is in the same directory).

### Step 3: Interact with the Server

1. Open `demo.html` in multiple browser tabs.
2. Click "Connect to signaling server" to establish a WebSocket connection.
3. Use the UI to send messages:
   - **Send to a Peer**: Enter the peer ID and message, then click "Send".
   - **Send to All**: Broadcast a message to all peers.
   - **Get Peer IDs**: Click "Get all peers IDs" to list connected peers.

## Files

- `main.go`: WebAssembly code to manage WebSocket connections.
- `index.html`: Basic page to interact with the WebSocket.
- `demo.html`: Example page to run the WebAssembly client.
- `main.wasm`: WebAssembly binary generated from `main.go`.
- `wasm_exec.js`: Required Go WebAssembly runtime for the client.