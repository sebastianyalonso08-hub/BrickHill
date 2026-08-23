package main

import (
  "bufio"
  "crypto/tls"
  "encoding/base64"
  "flag"
  "fmt"
  "io"
  "log"
  "math/rand"
  "net"
  "net/http"
  "os"
  "strings"
)

func wsConnect(url string) (net.Conn, *bufio.Reader, error) {
  // Minimal RFC6455 client handshake. Server-to-client frames are unmasked.
  u := strings.TrimPrefix(strings.TrimPrefix(url, "wss://"), "ws://")
  hostPath := strings.SplitN(u, "/", 2)
  host := hostPath[0]
  path := "/"
  if len(hostPath) == 2 { path += hostPath[1] }
  conn, err := tls.Dial("tcp", host, &tls.Config{ServerName: strings.Split(host, ":")[0], MinVersion: tls.VersionTLS12})
  if err != nil { return nil,nil,err }
  keyBytes := make([]byte,16); rand.Read(keyBytes); key:=base64.StdEncoding.EncodeToString(keyBytes)
  req := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path,host,key)
  if _,err=conn.Write([]byte(req)); err!=nil {conn.Close();return nil,nil,err}
  r:=bufio.NewReader(conn); status,err:=r.ReadString('\n'); if err!=nil{return nil,nil,err}
  if !strings.Contains(status,"101") {return nil,nil,fmt.Errorf("websocket handshake failed: %s",strings.TrimSpace(status))}
  for { line,err:=r.ReadString('\n'); if err!=nil{return nil,nil,err}; if line=="\r\n" {break} }
  return conn,r,nil
}

func writeWS(c net.Conn, payload []byte) error {
  n:=len(payload); hdr:=[]byte{0x82}
  if n<126 { hdr=append(hdr,byte(n)) } else if n<=65535 { hdr=append(hdr,126,byte(n>>8),byte(n)) } else { return fmt.Errorf("frame too large") }
  _,err:=c.Write(append(hdr,payload...)); return err
}
func readWS(r *bufio.Reader) ([]byte,bool,error) {
  b1,err:=r.ReadByte(); if err!=nil{return nil,false,err}; b2,err:=r.ReadByte(); if err!=nil{return nil,false,err}
  opcode:=b1&0x0f; fin:=b1&0x80!=0; masked:=b2&0x80!=0; n:=int(b2&0x7f)
  if n==126 {a,e:=r.ReadByte();if e!=nil{return nil,false,e};b,e:=r.ReadByte();if e!=nil{return nil,false,e};n=int(a)<<8|int(b)}
  if n==127{return nil,false,fmt.Errorf("large websocket frame unsupported")}
  mask:=make([]byte,4); if masked {if _,err=io.ReadFull(r,mask);err!=nil{return nil,false,err}}
  p:=make([]byte,n);if _,err=io.ReadFull(r,p);err!=nil{return nil,false,err};for i:=range p {if masked{p[i]^=mask[i%4]}}
  if !fin{return nil,false,fmt.Errorf("fragmented frame unsupported")}
  if opcode==8{return nil,true,nil}; if opcode==9 { return nil,false,nil }; if opcode!=2 && opcode!=1 {return nil,false,nil}
  return p,false,nil
}

func pipeTCPToWS(t net.Conn, ws net.Conn) error { buf:=make([]byte,65535); for {n,e:=t.Read(buf);if n>0{if we:=writeWS(ws,buf[:n]);we!=nil{return we}};if e!=nil{return e}} }
func pipeWSToTCP(r *bufio.Reader, ws net.Conn, t net.Conn) error {for {p,closed,e:=readWS(r);if e!=nil{return e};if closed{return io.EOF};if len(p)>0{if _,e=t.Write(p);e!=nil{return e}}}}

func main(){
  listen:=flag.String("listen","127.0.0.1:6510","legacy TCP listener")
  ws:=flag.String("ws","wss://brickhill.onrender.com/ws/legacy","WebSocket endpoint")
  token:=flag.String("token","","session token")
  logPath:=flag.String("log","brickhill-bridge.log","capture log")
  flag.Parse(); _=http.MethodGet
  if *token!="" { sep := "?"; if strings.Contains(*ws,"?") { sep = "&" }; *ws += sep+"bridge="+*token }
  f,_:=os.OpenFile(*logPath,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644);if f!=nil{defer f.Close();log.SetOutput(f)}
  ln,err:=net.Listen("tcp",*listen);if err!=nil{log.Fatal(err)}
  log.Printf("Brick Hill legacy->WebSocket bridge listening on %s -> %s",*listen,*ws)
  for { c,err:=ln.Accept();if err!=nil{continue};go func(t net.Conn){defer t.Close();wsconn,r,e:=wsConnect(*ws);if e!=nil{log.Printf("WS connect: %v",e);return};defer wsconn.Close();log.Printf("client %s connected",t.RemoteAddr());
      done:=make(chan error,2);go func(){done<-pipeTCPToWS(t,wsconn)}();go func(){done<-pipeWSToTCP(r,wsconn,t)}();e=<-done;log.Printf("client %s closed: %v",t.RemoteAddr(),e)
    }(c)}
}
