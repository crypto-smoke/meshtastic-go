package transport

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"sync"
)

const wildcardHandler = "__WILDCARD_HANDLER__"

// MessageHandler defines the function signature for a handler that processes a protobuf message.
type MessageHandler func(msg proto.Message)

// HandlerRegistry holds registered handlers for protobuf messages.
type HandlerRegistry struct {
	errorOnNoHandlers bool
	mu                sync.RWMutex
	handlers          map[string][]MessageHandler
}

// New creates a new instance of HandlerRegistry. Set errorOnNoHandler to true if you want HandleMessage to return
// an error if there are no handlers registered for a given msg when HandleMessage is called.
func NewHandlerRegistry(errorOnNoHandler bool) *HandlerRegistry {
	return &HandlerRegistry{
		errorOnNoHandlers: errorOnNoHandler,
		handlers:          make(map[string][]MessageHandler),
	}
}

// RegisterHandler registers a handler for a specific protobuf message type. If a nil msg is passed, then handler will
// be called for all received messages.
func (r *HandlerRegistry) RegisterHandler(msg proto.Message, handler MessageHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if msg == nil {
		r.handlers[wildcardHandler] = append(r.handlers[wildcardHandler], handler)
		return
	}
	msgName := proto.MessageName(msg)
	if msgName == "" {
		return // Could not get message name; consider logging or handling the error
	}
	name := string(msgName)
	r.handlers[name] = append(r.handlers[name], handler)
}

// HandleMessage invokes all registered handlers for the provided protobuf message, in the order they were registered.
func (r *HandlerRegistry) HandleMessage(msg proto.Message) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgName := proto.MessageName(msg)
	if msgName == "" {
		return fmt.Errorf("failed to get message name for type: %T", msg) // Could not get message name; consider logging or handling the error
	}
	name := string(msgName)
	//fmt.Println(name)
	var hasHandler bool

	// call specific handlers
	for _, handler := range r.handlers[name] {
		go handler(msg)
		hasHandler = true
	}

	// call wildcard handlers
	for _, handler := range r.handlers[wildcardHandler] {
		go handler(msg)
		hasHandler = true
	}

	if !hasHandler && r.errorOnNoHandlers {
		return fmt.Errorf("no handlers registered for message: %s", msgName)
	}

	return nil
}
