# BrickHill Play v14

The Play flow no longer tries to navigate to `brickhill://` from an async JavaScript callback.
That can be blocked by Chromium/Edge because the original user activation has been lost.

Play now prepares a ticket over HTTPS and then displays a real `Launch <game>` anchor.
The user click on that anchor directly opens `brickhill://play?...`, allowing Windows to hand it to BrickHillLauncher.

After the launcher starts, the existing ticket redeem/bridge flow is unchanged.
