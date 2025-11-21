package login

import (
	"errors"
	"io"
)

// PacketPayload is a type alias to make our function signatures clearer.
type PacketPayload []byte

// PacketHandler defines the contract for any object that can process
// a specific server-to-client packet. The Handle method receives the
// raw data payload (after the opcode) and returns the (potentially modified)
// payload that should be sent to the client.

type PacketHandler interface {
	Handle(r io.Reader) error
}

type HandlerRegistry struct {
	handlers map[uint8]PacketHandler
}

// NewHandlerRegistry creates an empty registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[uint8]PacketHandler),
	}
}

// Register assigns a handler to a specific opcode.
func (r *HandlerRegistry) Register(opcode uint8, handler PacketHandler) {
	r.handlers[opcode] = handler
}

func (r *HandlerRegistry) Get(opcode uint8) (PacketHandler, error) {
	handler, ok := r.handlers[opcode]
	if !ok {
		return nil, errors.New("no handler registered for opcode")
	}
	return handler, nil
}
