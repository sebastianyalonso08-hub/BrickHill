# Brick Hill Play Fix v12

This build keeps `Brick_Hill_Multiplayer.exe` and `BRICK.dll` unchanged.

## Critical network fix
The legacy `wsock32.dll` shim was rewriting the client's destination to `127.0.0.1:26137`, while the bridge listens on `127.0.0.1:6510`. v12 patches only those two port bytes in the shim so the unchanged client reaches the local bridge.

## Selected game flow
The website creates a launch ticket for the selected game. The launcher redeems it and starts the bridge with that game's connection context. The bridge connects to `/ws/legacy` with the selected game/connection, isolating traffic by website game room.

The original client still speaks its legacy GameMaker `mplay` protocol. The bridge transports that TCP stream over WSS. Full in-game session creation/joining still requires the server to implement the exact legacy `mplay` handshake.
