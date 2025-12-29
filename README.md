# z07

### The Legend of the Two Sevens

In the world of the wanderer, **Floor 7** is the surface. It is the beginning of every journey—the sunlit ground where the grass grows and the legends walk. It is the anchor of reality.

But in the world of the architect, there is another summit. At the very peak of the digital stack sits **Layer 7**—the **Zenith**. This is the Application Layer, the high-altitude realm where logic breathes and where the game’s heart truly beats.

**z07** is the phantom that lives at the intersection of these two worlds.

Forged in the speed of **Go** and hidden in the shadows of the Man-in-the-Middle, **z07** acts as a silent observer at the edge of the grid. It is more than a proxy; it is a vantage point. It sits at the **Zenith of the Seventh Floor**, translating the cold, raw binary of the wire into the vivid pulse of the game state.

While the client and server speak their private language, **z07** listens, learns, and pilots the flow. It is the invisible bridge between the pixel and the packet.

**Welcome to the middle. Welcome to z07.**

---



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
    go run cmd/client-patcher/main.go -binary "C:\Games\Tibia772\Tibia.exe" -ip 192.168.1.142
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
