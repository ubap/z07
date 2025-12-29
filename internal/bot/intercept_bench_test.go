package bot

import (
	"testing"
)

/*
*
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
*
goos: darwin
goarch: arm64
pkg: z07/internal/bot
cpu: Apple M4
BenchmarkInterceptC2SPacket_NotIntercepted
BenchmarkInterceptC2SPacket_NotIntercepted-10    	1000000000	         0.9911 ns/op
*/
func BenchmarkInterceptC2SPacket_NotIntercepted(b *testing.B) {
	bot := &Bot{}
	lookPacket := []byte{0xFF, 0x69, 0x7D, 0xE5, 0x7D, 0x07, 0xBA, 0x11, 0x01}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bot.InterceptC2SPacket(lookPacket)
	}
}
