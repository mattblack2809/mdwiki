package main

import (
  "fmt"
  "flag"
  "log"
  "matt/mdwiki"
  "net/http"
  "os"
)

func main() {
  port := flag.Int("port", 8000, "port to listen on")
  http.HandleFunc("/", handler)
  hostPort := fmt.Sprintf("localhost:%d",*port)
  log.Fatal(http.ListenAndServe(hostPort, nil))  // port 80 access perm error
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "<html><head><title>MD Wiki</title></head><body>")
  defer fmt.Fprintln(w, "</body></html>")
  fmt.Fprintln(w, mdwiki.PrintPath(r.URL.Path))
  path := "." + r.URL.Path
  finfo, err := os.Stat(path)
  if err != nil {
    fmt.Fprintf(w, "Error tying os.Stat(%s) : %q", path, err)
    return
  }
  if finfo.IsDir() {
    fmt.Fprintln(w, mdwiki.FmtDir(path))
  } else {
    mdwiki.PrintFile(w, path)
  }
}
