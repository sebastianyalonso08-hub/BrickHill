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
  "time"
  "strings"
  "syscall"
)

type redeemResp struct { OK bool `json:"ok"`; User struct{Username string `json:"username"`} `json:"user"`; GameID string `json:"gameId"`; WS string `json:"ws"`; Error string `json:"error"` }
func main(){
  if len(os.Args)<2 { return }
  u,err:=url.Parse(os.Args[1]);if err!=nil||u.Scheme!="brickhill"{return}
  q:=u.Query(); ticket:=q.Get("ticket"); game:=q.Get("game"); if ticket==""||game==""{return}
  base:=os.Getenv("BRICKHILL_API");if base==""{base="https://brickhill.onrender.com"}
  body,_:=json.Marshal(map[string]string{"ticket":ticket})
  r,err:=http.Post(base+"/api/client/redeem","application/json",bytes.NewReader(body));if err!=nil{return};defer r.Body.Close()
  var out redeemResp;json.NewDecoder(r.Body).Decode(&out);if !out.OK||out.WS==""{return}
  self,_:=os.Executable(); dir:=filepath.Dir(self); bridge:=filepath.Join(dir,"BrickHillNetworkBridge.exe"); gameExe:=filepath.Join(dir,"Brick_Hill_Multiplayer.exe")
  b:=exec.Command(bridge,"-ws",out.WS,"-token",ticket);b.SysProcAttr=&syscall.SysProcAttr{HideWindow:true};_ = b.Start()
  time.Sleep(300*time.Millisecond)
  // Preserve the original client; only provide the launch URI as an optional argument.
  g:=exec.Command(gameExe);g.Dir=dir;g.Env=append(os.Environ(),"BRICKHILL_GAME="+game,"BRICKHILL_USER="+out.User.Username,"BRICKHILL_LOCAL_SERVER=127.0.0.1:6510");g.SysProcAttr=&syscall.SysProcAttr{}
  _=g.Start()
  _=fmt.Sprintf("%s",strings.TrimSpace(game))
}
