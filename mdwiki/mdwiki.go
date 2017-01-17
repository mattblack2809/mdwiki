// package main provides the top level structure of mdwiki (Markdown widi).
// It reads options allowing defaults to be over-written firstly by the
// content of mdwiki.conf (read from the current directory) and then by
// command line options.  Run mdwiki -help for information on the options.
// An http server is started with a single handler function that stats the
// file pointed to by the URL to determine whether to ouput a directory listing
// or the file itself.
// Non-HTML files are output 'as is' while markwoen / HTML files are output with a
// clickable file path displayed at the top of the screen.
// Markdown files are converted to HTML on-the-fly using pandoc - with caching,
// and optionally with a table of content for each .md file generated out of
// the pandoc conversion.
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

// main reads options, registers a handler, and starts an http server
func main() {
  port := flag.String("port", "8000", "port to listen on")
  toc := flag.Bool("toc", false, "true to automatically generate tables of content")
  silent := flag.Bool("silent", false, "true supress log output")
  noCache := flag.Bool("no-cache", false, "true force pandoc each access")
  pandocArgs :=flag.String("pandoc-args", "", "addition options to pandoc")
  flag.Parse()

  mdwiki.ReadOptions() // load options from file, then over-write based on flags
  if *toc == true { mdwiki.Options["toc"] = "true"}
  if *noCache == true { mdwiki.Options["no-cache"] = "true"}
  if *silent == true { mdwiki.Options["silent"] = "true"}
  if *pandocArgs != "" { mdwiki.Options["pandoc-args"] = *pandocArgs}
  if *port != "8000" { mdwiki.Options["port"] = *port}


  if *silent {
    out,_ := os.Open(os.DevNull)
    log.SetOutput(out)
  }
  listenPort := *port
  optPort, ok := mdwiki.Options["port"]
  if *port == "8000" && ok {
    listenPort = optPort
  }
  http.HandleFunc("/", handler)
  log.Println("Starting mdwiki using options:", mdwiki.Options)

  hostPort := fmt.Sprintf("localhost:%s", listenPort)
  fmt.Fprintln(os.Stderr, http.ListenAndServe(hostPort, nil))  // port 80 access perm error
}

// handler stats the file pointed to by the URL, tests if it is a directory,
// and invokes the functions to output either a directory listing or the file
func handler(w http.ResponseWriter, r *http.Request) {
  path := "." + r.URL.Path
  finfo, err := os.Stat(path)
  if err != nil {
    fmt.Fprintf(w, "Error trying os.Stat(%s) : %q", path, err)
    return
  }
  if finfo.IsDir() {
    //mdwiki.FmtDir(w, path)
    mdwiki.HTMLDirList(w, path)
  } else {
    mdwiki.PrintFile(w, path)
  }
}
