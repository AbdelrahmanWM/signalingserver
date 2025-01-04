# WebRTC Mesh Network Client (WebAssembly) – Signaling Server Integration

This example demonstrates how to establish a full mesh WebRTC connection between multiple clients using WebAssembly (WASM) and a signaling server. Once the WebRTC connections are established, the clients can disconnect from the signaling server and continue communicating via WebRTC data channels.

## Features
- **WebRTC Peer-to-Peer Connections**: Establish peer-to-peer connections using WebRTC between multiple clients.
- **Signaling Server Communication**: Use the signaling server for initial offer/answer exchange and ICE candidate exchange.
- **Full Mesh Network**: Supports a full mesh WebRTC connection where every client connects to every other client.
- **Data Channel Communication**: After establishing connections, clients use WebRTC data channels to send/receive messages directly.
- **WebAssembly Support**: The client runs in the browser using WebAssembly (WASM), eliminating the need for a native application.
- **Log and Debugging**: View detailed logs of WebRTC connection states and data channel messages.

## Prerequisites
- Go 1.16+ installed.
- Web browser that supports WebAssembly and WebRTC.
- A WebRTC signaling server (e.g., the one from this project) running.

## Setup and Usage

### Step 1: Run the Signaling Server
Clone the signaling server repository:

```bash
git clone https://github.com/AbdelrahmanWM/signalingserver
```

Install dependencies and run the signaling server:

```bash
cd examples/webrtc_mesh_client
go run main.go
```

This will start the signaling server at `ws://localhost:8090/signalingserver`.

### Step 2: Set Up the WebRTC Client
Navigate to the `wasm/` directory:

```bash
cd examples/webrtc_mesh_client/wasm
```

Build the WASM client:

```bash
GOARCH=wasm GOOS=js go build -o main.wasm main.go
```

This will generate the `main.wasm` file.

Open the `index.html` file in more than one browser tab or window. This will simulate multiple clients connecting to the signaling server and interacting with each other. (ensure `main.wasm` is in the same directory as `index.html`).

### Step 3: Interact with the WebRTC Client

#### Connect to the Signaling Server
Click **"Connect to signaling server"** to establish a WebSocket connection to the signaling server.

- The signaling server handles the offer/answer exchange and ICE candidate exchange.

#### Get Peer IDs
Click **"Get all peer IDs"** to retrieve the list of connected peers from the signaling server.

#### Create Peer Connections

- Click **"New peer connection"** to create a new WebRTC connection to another peer by providing their peer ID.
- The signaling server will handle the exchange of offer/answer and ICE candidates to establish the peer-to-peer WebRTC connection.

#### Send Messages

- **Send to Peer**: Type a message in the message box and click **"Send"** to send a message to a specific peer.
- **Send to All**: Broadcast a message to all connected peers via WebRTC data channels.

#### Disconnect from the Signaling Server
Once all peer connections are established, you can disconnect from the signaling server by clicking **"Disconnect"**. The WebRTC data channels will continue to function without the signaling server.

### Step 4: Full Mesh WebRTC Communication

- In a full mesh network, each client connects to every other client, forming a complete peer-to-peer network.
- The WebRTC peer connections are established between all peers using the signaling server to exchange the necessary offer/answer and ICE candidates.
- After the initial connection setup via the signaling server, data channels are used for direct communication between peers.

### Step 5: Debugging and Logs

- View detailed logs in the **Logs** section to monitor the WebRTC connection process, including offer/answer exchange, ICE candidates, and message transmission.
- Use the **Clear Logs** button to reset the log display.

## WebRTC Details

### Peer-to-Peer Connections
Each client creates a WebRTC `RTCPeerConnection` object. This object manages the lifecycle of the peer connection, including:

- **Signaling**: The signaling server is responsible for facilitating the exchange of offer and answer messages between peers, as well as the exchange of ICE candidates.
- **Offer/Answer**: The initiating peer creates an offer containing an SDP message, which is sent to the receiving peer. The receiving peer then generates an answer in response.
- **ICE Candidates**: Both peers gather ICE candidates to facilitate the connection. These candidates are exchanged via the signaling server.
- **Data Channel**: Once the connection is established, a data channel is created between peers for message exchange.

### Data Channels

- **Single Channel**: A single data channel is created per WebRTC connection. This channel allows both peers to send and receive messages.
- **Multiple Channels**: You can create multiple data channels per connection if needed, though this example uses only a single channel per connection.

### ICE Candidate Gathering
WebRTC uses ICE candidates to discover the best network path between peers. These candidates are gathered and sent to the remote peer via the signaling server. The signaling server acts as a relay for exchanging ICE candidates during the initial connection setup phase.

### Example SDP Message
Here’s an example of an SDP offer message that could be sent between peers during the connection establishment phase:

```plaintext
v=0
o=- 1389558159263351010 2 IN IP4 127.0.0.1
s=-
t=0 0
a=group:BUNDLE 0
a=msid-semantic: WMS
m=application 9 UDP/DTLS/SCTP webrtc-datachannel
c=IN IP4 0.0.0.0
a=ice-ufrag:c9wH
a=ice-pwd:UcWfX7SBzpAq+/Ne/ukGbZer
a=ice-options:trickle
a=fingerprint:sha-256 65:3C:C0:C8:B4:85:A2:4C:69:AC:7A:BD:E4:2C:D4:94:9D:E0:E7:99:A3:0F:13:87:7C:77:CE:A4:71:9E:62:E6
a=setup:active
a=mid:0
a=sctp-port:5000
a=max-message-size:262144
```

## Debugging WebRTC Connection

Use the **Logs** section in the UI to monitor the WebRTC connection’s status, including:

- ICE connection state
- Offer/Answer exchange status
- Data channel status
- Any errors or issues during the WebRTC connection setup

### Troubleshooting

- **No Common Codecs**: If you encounter issues where peers cannot connect due to codec mismatches, ensure that both peers support compatible WebRTC codecs.
- **ICE Candidate Issues**: If you have trouble with ICE candidates or connection failures, ensure that the signaling server is correctly relaying ICE candidates between peers. Use browser developer tools to inspect the ICE connection states.
- **SDP/Offer-Answer Issues**: Double-check that the offer/answer exchange is working correctly, and verify that the SDP messages are properly formatted and transmitted.

## Files

- **signalingserverconn/main.go**: Go code for handling the signaling server's WebSocket connections and peer exchange logic.
- **webrtcpeerconn/webrtc_peer_conn.go**: Go code for managing WebRTC peer connections and related WebRTC signaling logic.
- **wasm/index.html**: HTML page for interacting with the WebAssembly client.
- **wasm/main.wasm**: WebAssembly binary generated from `main.go`, enabling WebRTC functionality in the browser.
- **wasm/wasm_exec.js**: WebAssembly runtime required by Go to run in the browser.
- **utils/**: Utility files that may contain helper functions used by other parts of the application.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
