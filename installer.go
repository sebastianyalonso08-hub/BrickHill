//go:build windows
package main

import (
  _ "embed"
  "archive/zip"
  "bytes"
  "fmt"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "syscall"
  "time"
  "unsafe"
)

//go:embed BrickHillClient.zip
var clientZip []byte

func msg(title, text string, flags uintptr) { dll:=syscall.NewLazyDLL("user32.dll"); p:=dll.NewProc("MessageBoxW"); t,_:=syscall.UTF16PtrFromString(title); m,_:=syscall.UTF16PtrFromString(text); p.Call(0,uintptr(unsafe.Pointer(m)),uintptr(unsafe.Pointer(t)),flags) }
func run(name string,args ...string) error { c:=exec.Command(name,args...); c.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}; return c.Run() }
func main(){
  install:=filepath.Join(os.Getenv("LOCALAPPDATA"),"BrickHill")
  if err:=os.MkdirAll(install,0755);err!=nil {msg("Brick Hill Installer","Could not create the installation folder.\n\n"+err.Error(),0x10);return}
  zr,err:=zip.NewReader(bytes.NewReader(clientZip),int64(len(clientZip)));if err!=nil{msg("Brick Hill Installer","The bundled client package is damaged.",0x10);return}
  for _,f:=range zr.File { name:=filepath.Clean(f.Name); if name=="."||strings.HasPrefix(name,"..")||filepath.IsAbs(name)||strings.Contains(name,string(filepath.Separator)+".."+string(filepath.Separator)){msg("Brick Hill Installer","Unsafe client package path detected.",0x10);return}; dst:=filepath.Join(install,name); if !strings.HasPrefix(dst,install+string(filepath.Separator)) && dst!=install {msg("Brick Hill Installer","Invalid installation path.",0x10);return}; if f.FileInfo().IsDir(){os.MkdirAll(dst,0755);continue}; if err=os.MkdirAll(filepath.Dir(dst),0755);err!=nil{msg("Brick Hill Installer",err.Error(),0x10);return}; in,e:=f.Open();if e!=nil{msg("Brick Hill Installer",e.Error(),0x10);return}; out,e:=os.Create(dst);if e==nil{_,e=io.Copy(out,in);out.Close()};in.Close();if e!=nil{msg("Brick Hill Installer","Failed to extract the client.\n\n"+e.Error(),0x10);return} }
  launcher:=filepath.Join(install,"BrickHillLauncher.exe")
  // Register brickhill:// for the current Windows user.
  if err=run("reg.exe","ADD",`HKCU\Software\Classes\brickhill`,`/ve`,`/d`,`URL:Brick Hill Protocol`,`/f`);err!=nil{msg("Brick Hill Installer","Could not register the Brick Hill URL protocol.\n\n"+err.Error(),0x10);return}
  if err=run("reg.exe","ADD",`HKCU\Software\Classes\brickhill`,`/v`,`URL Protocol`,`/t`,`REG_SZ`,`/d`,``,`/f`);err!=nil{msg("Brick Hill Installer",err.Error(),0x10);return}
  if err=run("reg.exe","ADD",`HKCU\Software\Classes\brickhill\shell\open\command`,`/ve`,`/d`,`"`+launcher+`" "%1"`,`/f`);err!=nil{msg("Brick Hill Installer","Could not register the launcher command.\n\n"+err.Error(),0x10);return}
  // Create a desktop shortcut that opens the website.
  desktop:=filepath.Join(os.Getenv("USERPROFILE"),"Desktop","Brick Hill.lnk")
  ps:=`$s=(New-Object -ComObject WScript.Shell).CreateShortcut('`+strings.ReplaceAll(desktop,"'","''")+`');$s.TargetPath='https://brickhill.onrender.com/';$s.WorkingDirectory='`+strings.ReplaceAll(install,"'","''")+`';$s.Description='Brick Hill';$s.Save()`
  _=run("powershell.exe","-NoProfile","-ExecutionPolicy","Bypass","-Command",ps)
  // Write install marker for diagnostics.
  _=os.WriteFile(filepath.Join(install,"installed.txt"),[]byte(time.Now().Format(time.RFC3339)+"\nbrickhill:// registered\n"),0644)
  msg("Brick Hill Installer","Brick Hill Client installed successfully!\n\nLocation:\n"+install+"\n\nThe brickhill:// protocol is registered and a Brick Hill shortcut was added to your desktop.\n\nYou can now return to the website and click Play.",0x40)
  _=fmt.Sprint(install)
}
