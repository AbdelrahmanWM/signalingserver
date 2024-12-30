# WebSocket-based Signaling Server

A basic websocket-based signaling server written in Go, designed to facilitate peer-to-peer communication in a mesh network. This server helps WebRTC peers discover each other, exchange session descriptions (SDP), and handle ICE candidates. It can be used in projects requiring direct browser-to-browser communication without relying on central media servers.

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Use Cases](#use-cases)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [License](#license)

## Description

This signaling server implements basic signaling functionalities required for WebRTC mesh communication, including:

- Peer discovery and connection management.
- Exchanging SDP offers and answers.
- ICE candidate exchange for NAT traversal.

The server supports multiple clients and enables them to connect and communicate directly. It works by passing signaling messages between peers and can handle client disconnections.

## Features

- Peer-to-peer WebRTC communication in a full mesh topology.
- Supports WebSocket-based signaling for real-time messaging.
- Allows appending of sender IDs in messages for better traceability.
- Optional logging of peer connections and interactions.
- Graceful handling of peer disconnects and connection cleanup.

## Use Cases

- **Real-time Communication**: Use the server as the foundation for building WebRTC-based real-time communication applications such as video calls, audio calls, or file sharing.
- **Mesh Networks**: Suitable for decentralized, mesh-style networks where peers directly communicate without the need for a central server.
- **Collaborative Tools**: Perfect for applications that require direct communication between users, such as collaborative whiteboards or document editing.

## Getting Started

To use this signaling server, follow the steps below:

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/) 1.18 or higher
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket handling
- A WebSocket client (browser, or custom client in your application)

To install the Gorilla WebSocket package, run:

```bash
go get github.com/gorilla/websocket
```

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/webrtc-signaling-server.git
   cd webrtc-signaling-server
   ```

2. Install the dependencies:

   ```bash
   go mod tidy
   ```

3. Build and run the server:

   ```bash
   go run main.go
   ```
