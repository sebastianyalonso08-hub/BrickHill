# Brick Hill — conexión del cliente original

Esta versión mantiene `Brick_Hill_Multiplayer.exe` y `BRICK.dll` sin modificaciones. Sus SHA-256 coinciden con los archivos originales entregados.

## Cómo funciona

Website → `brickhill://` → Launcher → cliente original → `wsock32.dll` (shim de red) → `127.0.0.1:6510` → `BrickHillNetworkBridge.exe` → WSS → Render `/ws/legacy`.

El shim es la única pieza que intercepta la red del cliente. Redirige las conexiones TCP del cliente a `127.0.0.1:6510`; las demás funciones Winsock se reenvían a Windows mediante `ws2_32.dll`.

## Instalar en Windows

1. Extrae todo el contenido de esta carpeta en una misma carpeta.
2. Abre PowerShell en esa carpeta.
3. Ejecuta:

    powershell -ExecutionPolicy Bypass -File .\install-brickhill.ps1

4. Abre `https://brickhill.onrender.com/`, inicia sesión y pulsa **Play**.
5. Windows abrirá el protocolo `brickhill://`; el launcher validará el ticket, arrancará el bridge y luego el EXE original.

## Render

El servidor necesita `npm install` y `npm start`. El WebSocket público es `wss://TU-DOMINIO-RENDER/ws/legacy`.

La conexión WebSocket usa un ticket temporal emitido por `/api/client/launch` y canjeado una sola vez por `/api/client/redeem`.

## Nota técnica

El cliente es un build legacy de GameMaker que usa `mplay_init_tcpip`. La documentación de GameMaker describe esa función como la inicialización del transporte TCP/IP y las funciones `mplay_session_*` como la capa de sesiones. El servidor de esta versión funciona como un túnel binario por sala: no modifica los paquetes legacy.

La compatibilidad completa del protocolo `mplay` depende de que el tráfico TCP del build concreto pueda funcionar a través del túnel. Si el cliente usa tráfico adicional que no sea TCP, el siguiente ajuste debe ser una capa de transporte adicional; el EXE y DLL no necesitan ser reemplazados.
