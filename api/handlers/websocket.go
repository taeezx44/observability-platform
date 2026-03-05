package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles real-time log streaming
type WebSocketHandler struct {
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	TraceID   string `json:"trace_id,omitempty"`
	SpanID    string `json:"span_id,omitempty"`
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo
			},
		},
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.register <- conn

	// Start goroutines for this connection
	go h.writePump(conn)
	go h.readPump(conn)
}

func (h *WebSocketHandler) StartBroadcaster() {
	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				log.Printf("Client connected. Total clients: %d", len(h.clients))

				// Send welcome message
				welcome := LogMessage{
					Timestamp: time.Now().Format(time.RFC3339),
					Level:     "INFO",
					Service:   "websocket",
					Message:   "Connected to real-time log stream",
				}
				data, _ := json.Marshal(welcome)
				h.broadcast <- data

			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					client.Close()
					log.Printf("Client disconnected. Total clients: %d", len(h.clients))
				}

			case message := <-h.broadcast:
				for client := range h.clients {
					err := client.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Printf("Error sending to client: %v", err)
						h.unregister <- client
					}
				}
			}
		}
	}()

	// Start generating demo logs
	go h.generateDemoLogs()
}

func (h *WebSocketHandler) writePump(conn *websocket.Conn) {
	defer conn.Close()

	for {
		select {
		case message, ok := <-h.broadcast:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				h.unregister <- conn
				return
			}
		}
	}
}

func (h *WebSocketHandler) readPump(conn *websocket.Conn) {
	defer conn.Close()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			h.unregister <- conn
			break
		}
	}
}

func (h *WebSocketHandler) generateDemoLogs() {
	services := []string{"api", "worker", "database", "auth", "cache", "queue"}
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	messages := map[string][]string{
		"api": {
			"HTTP request processed",
			"Database query executed",
			"Authentication successful",
			"Rate limit exceeded",
			"Cache miss",
			"Response sent to client",
		},
		"worker": {
			"Job started",
			"Processing queue items",
			"Job completed successfully",
			"Retrying failed job",
			"Queue depth updated",
			"Worker scaling event",
		},
		"database": {
			"Connection established",
			"Query executed",
			"Transaction committed",
			"Index rebuilt",
			"Backup started",
			"Connection pool updated",
		},
		"auth": {
			"User login attempt",
			"Token validated",
			"Session created",
			"Permission check",
			"Token refreshed",
			"User logged out",
		},
		"cache": {
			"Cache hit",
			"Cache miss",
			"Cache entry expired",
			"Cache cleared",
			"Warm-up started",
			"Memory usage updated",
		},
		"queue": {
			"Message published",
			"Message consumed",
			"Queue size updated",
			"Dead letter queue",
			"Consumer lag",
			"Partition rebalanced",
		},
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service := services[int(time.Now().UnixNano())%len(services)]
			level := levels[int(time.Now().UnixNano())%len(levels)]
			serviceMessages := messages[service]
			message := serviceMessages[int(time.Now().UnixNano())%len(serviceMessages)]

			logMsg := LogMessage{
				Timestamp: time.Now().Format(time.RFC3339),
				Level:     level,
				Service:   service,
				Message:   message,
				TraceID:   fmt.Sprintf("trace_%x", time.Now().UnixNano()),
				SpanID:    fmt.Sprintf("span_%x", time.Now().UnixNano()/1000),
			}

			data, err := json.Marshal(logMsg)
			if err != nil {
				log.Printf("Error marshaling log message: %v", err)
				continue
			}

			h.broadcast <- data
		}
	}
}
