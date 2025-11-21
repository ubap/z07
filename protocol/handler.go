package protocol

// PacketPayload is a type alias to make our function signatures clearer.
type PacketPayload []byte

// PacketHandler defines the contract for any object that can process
// a specific server-to-client packet. The Handle method receives the
// raw data payload (after the opcode) and returns the (potentially modified)
// payload that should be sent to the client.

type PacketHandler interface {
	Handle(payload PacketPayload) (PacketPayload, error)
}

type HandlerRegistry struct {
	handlers       map[uint8]PacketHandler
	DefaultHandler PacketHandler
}

// NewHandlerRegistry creates an empty registry.
func NewHandlerRegistry(defaultHandler PacketHandler) *HandlerRegistry {
	return &HandlerRegistry{
		handlers:       make(map[uint8]PacketHandler),
		DefaultHandler: defaultHandler,
	}
}

// Register assigns a handler to a specific opcode.
func (r *HandlerRegistry) Register(opcode uint8, handler PacketHandler) {
	r.handlers[opcode] = handler
}

func (r *HandlerRegistry) Get(opcode uint8) PacketHandler {
	handler, ok := r.handlers[opcode]
	if !ok {
		return r.DefaultHandler
	}
	return handler
}
