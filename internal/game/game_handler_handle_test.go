package game

import (
	"errors"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type MockConn struct {
	ReadErr  error
	WriteErr error
	Closed   bool
}

func (m *MockConn) ReadMessage() ([]byte, error) {
	if m.ReadErr != nil {
		return nil, m.ReadErr
	}
	// Simulate an idle connection that eventually errors
	time.Sleep(50 * time.Millisecond)
	return nil, errors.New("connection closed")
}

func (m *MockConn) WriteMessage(p []byte) error           { return m.WriteErr }
func (m *MockConn) SendPacket(p protocol.Encodable) error { return nil }
func (m *MockConn) EnableXTEA(key [4]uint32)              {}
func (m *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
}
func (m *MockConn) Close() error { m.Closed = true; return nil }

// --- Tests ---

func TestHandle_Lifecycle(t *testing.T) {
	clientMock := &MockConn{}
	serverMock := &MockConn{}

	var capturedState *state.GameState

	handler := &GameHandler{
		TargetAddr: "127.0.0.1:7171",
		SessionInitializer: func(addr string, conn protocol.Connection) (*packets.LoginRequest, protocol.Connection, error) {
			require.Equal(t, "127.0.0.1:7171", addr)

			return &packets.LoginRequest{
				CharacterName: "TestPlayer",
			}, serverMock, nil
		},
		OnSessionStart: func(s *GameSession) {
			capturedState = s.State
		},
	}

	done := make(chan struct{})
	go func() {
		handler.Handle(clientMock)
		close(done)
	}()

	<-done

	require.NotNil(t, capturedState, "GameState should have been created")
	actualName := capturedState.CaptureFrame().Player.Name
	require.Equal(t, "TestPlayer", actualName)
}

func TestHandle_InitSessionFailure(t *testing.T) {
	clientMock := &MockConn{}

	handler := &GameHandler{
		SessionInitializer: func(addr string, conn protocol.Connection) (*packets.LoginRequest, protocol.Connection, error) {
			return nil, nil, errors.New("auth failed")
		},
	}

	done := make(chan struct{})
	go func() {
		handler.Handle(clientMock)
		close(done)
	}()

	select {
	case <-done:
		// Success: Handle exited early due to auth error
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Handle should have exited immediately on auth failure")
	}
}
