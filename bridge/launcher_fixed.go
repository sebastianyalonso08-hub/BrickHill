//go:build windows
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "os"
  "os/exec"
  "path/filepath"
  "syscall"
  "time"
  "unsafe"
)

type redeemResp struct { OK bool `json:"ok"`; User struct{Username string `json:"username"`} `json:"user"`; GameID string `json:"gameId"`; WS string `json:"ws"`; Error string `json:"error"` }

func msg(title, text string) { dll:=syscall.NewLazyDLL("user32.dll"); p:=dll.NewProc("MessageBoxW"); t,_:=syscall.UTF16PtrFromString(title); m,_:=syscall.UTF16PtrFromString(text); p.Call(0,uintptr(unsafe.Pointer(m)),uintptr(unsafe.Pointer(t)),0x10) }
func logLine(s string){ p:=filepath.Join(os.TempDir(),"BrickHillLauncher.log"); f,_:=os.OpenFile(p,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644); if f!=nil { defer f.Close(); fmt.Fprintf(f,"%s %s\n",time.Now().Format(time.RFC3339),s) } }
func main(){
  if len(os.Args)<2 { msg("Brick Hill Launcher","This launcher must be opened by the Brick Hill website.\n\nIf you are installing the client, run BrickHillInstaller.exe first."); return }
  raw:=os.Args[1]; u,err:=url.Parse(raw); if err!=nil || u.Scheme!="brickhill" { logLine("invalid URI: "+raw); msg("Brick Hill Launcher","Invalid Brick Hill launch link."); return }
  q:=u.Query(); ticket:=q.Get("ticket"); game:=q.Get("game");
  if ticket=="" || game=="" { msg("Brick Hill Launcher","The launch link is missing the game or session ticket.\n\nGo back to Brick Hill and click Play again."); return }
  base:=os.Getenv("BRICKHILL_API"); if base=="" {base="https://brickhill.onrender.com"}
  logLine("redeeming game="+game)
  body,_:=json.Marshal(map[string]string{"ticket":ticket}); r,err:=http.Post(base+"/api/client/redeem","application/json",bytes.NewReader(body)); if err!=nil { logLine("redeem error: "+err.Error()); msg("Brick Hill Launcher","Could not connect to the Brick Hill server.\n\n"+err.Error()); return }; defer r.Body.Close()
  var out redeemResp; if err=json.NewDecoder(r.Body).Decode(&out); err!=nil {msg("Brick Hill Launcher","The server returned an invalid response.");return}; if !out.OK || out.WS=="" {msg("Brick Hill Launcher","Brick Hill could not start this game.\n\n"+out.Error);return}
  self,_:=os.Executable(); dir:=filepath.Dir(self); bridge:=filepath.Join(dir,"BrickHillNetworkBridge.exe"); gameExe:=filepath.Join(dir,"Brick_Hill_Multiplayer.exe")
  if _,err=os.Stat(bridge);err!=nil {msg("Brick Hill Launcher","BrickHillNetworkBridge.exe is missing.\n\nPlease reinstall the Brick Hill Client.");return}
  if _,err=os.Stat(gameExe);err!=nil {msg("Brick Hill Launcher","Brick_Hill_Multiplayer.exe is missing.\n\nPlease reinstall the Brick Hill Client.");return}
  b:=exec.Command(bridge,"-ws",out.WS,"-token",ticket,"-listen","127.0.0.1:6510"); b.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}; if err=b.Start();err!=nil {msg("Brick Hill Launcher","Could not start the network bridge.\n\n"+err.Error());return}
  time.Sleep(500*time.Millisecond)
  g:=exec.Command(gameExe); g.Dir=dir; g.Env=append(os.Environ(),"BRICKHILL_GAME="+game,"BRICKHILL_USER="+out.User.Username,"BRICKHILL_LOCAL_SERVER=127.0.0.1:6510"); if err=g.Start();err!=nil {msg("Brick Hill Launcher","Could not start Brick Hill.\n\n"+err.Error());return}
  logLine("game started pid="+fmt.Sprint(g.Process.Pid));
}
