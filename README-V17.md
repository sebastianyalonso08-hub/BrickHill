# Brick Hill v17 — Auto-Join + MPlay/DirectPlay Broker Foundation

This build keeps the original game executable and adds:

- website launch tickets with gameId;
- launcher starts the original client;
- launcher attempts the legacy startup flow automatically: blank IP, username, Host/Join keyboard sequence;
- per-game WebSocket rooms;
- corrected local bridge target to 127.0.0.1:6510;
- MPlay/DirectPlay broker foundation (TCP relay + per-game routing).

## Important

The TCP broker is a compatibility foundation, not a full reimplementation of the DirectPlay 4 wire protocol. The original GameMaker 8 client expects real MPlay/DirectPlay session enumeration and joining. If the client still reports “Session does not exist”, the next task is implementing those DirectPlay packets rather than changing the website UI.
