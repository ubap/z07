7.72
MITM Proxy

----
graph TD
Main[cmd/main.go] --> Game[internal/game]
Game --> Packets[internal/packets]
Game --> Model[internal/model]
Game --> Protocol[internal/protocol]

    Packets --> Protocol
    
    Model --> Nothing
    Protocol --> Nothing