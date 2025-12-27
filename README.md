# gotp

gotp is a lightweight, custom HTTP server implementation written in Go. It operates over raw TCP connections, featuring a custom HTTP parser and a concurrent routing system without relying on the standard `net/http` library for request handling.

## Features

* **Custom Parser:** Manually processes HTTP 1.1 request lines, headers, and bodies
* **Routing:** Supports exact-match routing based on HTTP verbs and paths
* **Concurrency:** Handles each incoming connection in a separate goroutine

## Screenshots
![](https://github.com/goldeneas/gotp/blob/main/screenshots/index.png?raw=true)
