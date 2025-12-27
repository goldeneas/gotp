# gotp

gotp is a lightweight, custom HTTP server implementation written in Go. It operates over raw TCP connections, featuring a custom HTTP parser and a concurrent routing system without relying on the standard `net/http` library for request handling.

## Features

* **Custom HTTP/1.1 Parser:** Manually processes request lines, headers, query parameters, and bodies using `bufio` readers
* **TCP Connection Management:** Handles keep-alive connections with configurable timeouts and distinct goroutines per connection
* **Exact-Match Routing:** Simple, thread-safe routing mechanism for mapping HTTP verbs and paths to handler functions
* **Static File Serving:** Built-in support for serving static assets
* **Smart Caching:** In-memory file caching with an integrated file watcher that automatically invalidates cache entries when files on disk are modified
