// +build windows
package mdwiki

import (
  "os"
  "time"
  "syscall"

)

// determine the OS and make correct syscall
// refer to https://github.com/djherbis/times/blob/master/times_linux.go
func accessTime(fi os.FileInfo) string {
  return tfmt(getWindowsAtime(fi))
}

//  Linux
//  func getTimespec(fi os.FileInfo) Timespec {
//	  var t timespec
//	  stat := fi.Sys().(*syscall.Win32FileAttributeData)
//	  t.atime.v = time.Unix(0, stat.LastAccessTime.Nanoseconds())
//	  return t
//  }

// WINDOWS
//https://github.com/djherbis/times/blob/master/times_windows.go
func getWindowsAtime(fi os.FileInfo) (time.Time) {
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
  t := time.Unix(0, stat.LastAccessTime.Nanoseconds())
  //	t = time.Unix(0, stat.LastWriteTime.Nanoseconds())
  //	t = time.Unix(0, stat.CreationTime.Nanoseconds())
  return t
}
