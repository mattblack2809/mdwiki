// +build windows
package mdwiki

import (
  "os"
  "time"
  "syscall"

)

// refer to https://github.com/djherbis/times/blob/master/times_linux.go
func accessTime(fi os.FileInfo) string {
  return tfmt(getWindowsAtime(fi))
}

//  Linux
//  func getLinuxAtime(fi os.FileInfo) time.Time {
//    stat := fi.Sys().(*syscall.Stat_t)
//    t := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
//    return t
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
