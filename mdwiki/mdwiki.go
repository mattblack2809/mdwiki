package main

import (
  "fmt"
  "flag"
  "log"
  "matt/mdwiki"
  "net/http"
  "os"
//  "time"
)
var toc *bool
func main() {
  port := flag.Int("port", 8000, "port to listen on")
  toc = flag.Bool("toc", false, "true to automatically generate tables of content")
  flag.Parse()
  //mdwiki.HTMLDirList(os.Stdout, ".")
  //os.Exit(0)
  http.HandleFunc("/", handler)
  hostPort := fmt.Sprintf("localhost:%d",*port)
  log.Fatal(http.ListenAndServe(hostPort, nil))  // port 80 access perm error
}

func handler(w http.ResponseWriter, r *http.Request) {
  path := "." + r.URL.Path
  finfo, err := os.Stat(path)
  if err != nil {
    fmt.Fprintf(w, "Error tying os.Stat(%s) : %q", path, err)
    return
  }
  if finfo.IsDir() {
    //mdwiki.FmtDir(w, path)
    mdwiki.HTMLDirList(w, path)
  } else {
    mdwiki.PrintFile(w, path, *toc)
  }
}
