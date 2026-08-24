# Brick Hill v18 — 3 Games

This build keeps Render and the existing website. It defines the three official game sessions and routes WebSocket players by gameId:

- 1001 — Brick Hill Classic
- 1002 — Brick Obby
- 1003 — Build & Chill

The launcher continues to receive `gameId` from the website ticket and starts the original client. The current bridge remains a compatibility relay; a complete DirectPlay 4 broker is still required for native MPlay session discovery/join.
