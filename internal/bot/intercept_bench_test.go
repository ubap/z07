package bot

import (
	"testing"
	"z07/internal/game/packets"
)

/*
goos: darwin
goarch: arm64
pkg: z07/internal/bot
cpu: Apple M4
BenchmarkInterceptC2SPacket
BenchmarkInterceptC2SPacket-10    	25687474	        47.05 ns/op
*/
func BenchmarkInterceptC2SPacket_Look(b *testing.B) {
	bot := &Bot{}
	lookPacket := []byte{0x8C, 0x69, 0x7D, 0xE5, 0x7D, 0x07, 0xBA, 0x11, 0x01}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bot.InterceptC2SPacket(lookPacket)
	}
}

/*
goos: darwin
goarch: arm64
pkg: z07/internal/bot
cpu: Apple M4
BenchmarkInterceptC2SPacket_NotIntercepted
BenchmarkInterceptC2SPacket_NotIntercepted-10    	1000000000	         0.9911 ns/op
*/
func BenchmarkInterceptC2SPacket_NotIntercepted(b *testing.B) {
	bot := &Bot{}
	unknownPacket := []byte{0xFF, 0x69, 0x7D, 0xE5, 0x7D, 0x07, 0xBA, 0x11, 0x01}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bot.InterceptC2SPacket(unknownPacket)
	}
}

/*
goos: darwin
goarch: arm64
pkg: z07/internal/bot
cpu: Apple M4
BenchmarkInterceptS2CPacket_NotIntercepted
BenchmarkInterceptS2CPacket_NotIntercepted-10    	1000000000	         0.9783 ns/op
*/
func BenchmarkInterceptS2CPacket_NotIntercepted(b *testing.B) {
	bot := &Bot{}
	unknownPacket := []byte{0xFF}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bot.InterceptS2CPacket(unknownPacket)
	}
}

/*
goos: darwin
goarch: arm64
pkg: z07/internal/bot
cpu: Apple M4
BenchmarkInterceptS2CPacket_Login
BenchmarkInterceptS2CPacket_Login-10    	37234057	        33.04 ns/op
*/
func BenchmarkInterceptS2CPacket_Login(b *testing.B) {
	bot := &Bot{}
	loginPacket := []byte{byte(packets.S2CSLoginQueue)}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bot.InterceptS2CPacket(loginPacket)
	}
}
