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
    go run tools/client-patcher/main.go -binary "C:\Games\Tibia772\Tibia.exe"
    ```
2.  This creates `Tibia_patched.exe` inside `C:\Games\Tibia772\`.
3.  **Run `Tibia_patched.exe`** from that folder to play.

---

### 🔑 RSA Key Finder (`rsa_finder.go`)

This tool scans a Tibia client binary to locate the RSA Public Key. It is useful if you are using a modified client or a different protocol version and need to find the memory offset for patching, or if you need to extract the original key to use in a server config.

#### Usage

**Basic Scan:**
Searches `Tibia.exe` in the current folder and saves the key to `rsa_key.txt`.
```bash
go run tools/rsa-finder/main.go
```

**Custom Paths:**
Specify a specific binary path and output location.
```bash
go run tools/rsa-finder/main.go --binary "C:\Games\Tibia772\Tibia.exe" --output "my_key.txt"
```

#### Output Explanation
The tool will output the **File Offset** (Address) where the key starts.

```text
---------------------------------------------------
Found RSA Key!
File Offset (Decimal): 1422880
File Offset (Hex):     0x15B620  <-- This is the address used in the Patcher
Key Length:            309 digits
---------------------------------------------------
Key successfully saved to: RSA.txt
```

*   **For Patcher Developers:** Use the **Hex Offset** (`0x15B620`) to update the `RSAAddress` constant in `client_patcher.go` if you are supporting a new client version.
*   **For Proxy Users:** The content of `rsa_key.txt` is the Public Key the client uses. The proxy must use this key to properly encrypt the communication with the server.


----
Package structure (outdated)

    Login --> Packets
    Login --> Protocol

    Packets --> Protocol

    Protocol --> Nothing
    Model --> Nothing