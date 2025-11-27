# A Man-in-the-Middle (MITM) proxy for the Tibia MMO, built to explore network programming and protocol reverse engineering in Go.

### 📦 Setup Guide

This project does **not** distribute copyrighted game assets or binaries. You must provide the files from your own legal copy of the Tibia 7.72 client.

#### 1. Configure the Bot (Server Side)
The bot needs to know about game physics (walls, stackable items).
1.  Copy `Tibia.dat` into the `data/772` folder of this project.
2.  Run the converter: `go run cmd/tools/dat_to_json.go`

#### 2. Patch your Client (Player Side)
You need a modified client to connect to the bot.
1.  Run the patcher, pointing to your Tibia installation:
    ```bash
    go run cmd/tools/client_patcher.go -binary "C:\Games\Tibia772\Tibia.exe"
    ```
2.  This creates `Tibia_patched.exe` inside `C:\Games\Tibia772\`.
3.  **Run `Tibia_patched.exe`** from that folder to play.

---

### 💡 Developer Note: The RSA Key
The patcher above replaces the client's public key with the standard "OTPublicRSA". For your Proxy to successfully decrypt the login packet, your Go code (`internal/crypto`) **must** be using the matching Private Key.

If you haven't already, ensure your server code uses this Private Key (Standard OT Key):

```go
// internal/crypto/rsa.go

// Matches the Public Key injected by client_patcher.go
var OTServPrivateKey = []byte{
    // ... (The big private key definition) ...
}
```

If these keys do not match, `ParseCredentialsPacket` will fail with "RSA encryption failed" or garbage data.


----
Package structure (outdated)

    Login --> Packets
    Login --> Protocol

    Packets --> Protocol

    Protocol --> Nothing
    Model --> Nothing