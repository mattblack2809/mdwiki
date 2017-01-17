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
  "path/filepath"
//  "time"
)

var absRoot string // the top level 'root' displayed

func main() {
  port := flag.String("port", "8000", "port to listen on")
  toc := flag.Bool("toc", false, "-toc=true automatically generates table of content")
  silent := flag.Bool("silent", false, "-silent=true suppresses log output")
  noCache := flag.Bool("no-cache", false, "-no-cache=true invokes pandoc for every access")
  pandocArgs :=flag.String("pandoc-args", "", "additional options passed to pandoc")
  root :=flag.String("root", "", "path to document root (absolute or relative to cwd)")
  flag.Parse()

  mdwiki.ReadOptions() // load options from file, then over-write based on flags
  if *toc == true { mdwiki.Options["toc"] = "true"}
  if *noCache == true { mdwiki.Options["no-cache"] = "true"}
  if *silent == true { mdwiki.Options["silent"] = "true"}
  if *pandocArgs != "" { mdwiki.Options["pandoc-args"] = *pandocArgs}
  var err error
  if *root != "" {
    absRoot, err = filepath.Abs(*root)
    if err != nil {log.Fatal("absroot error")}
  } else if r,ok := mdwiki.Options["root"]; ok{
    absRoot, err = filepath.Abs(r)
    if err != nil {log.Fatal("absroot error")}
  } else {
    absRoot,err = filepath.Abs(".")
    if err != nil {log.Fatal("absroot error")}
  }
  if *port != "8000" { mdwiki.Options["port"] = *port}

  //fmt.Println(filepath.Abs(*root))
  //os.Exit(0)

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
  urlPath := r.URL.Path
  var absPath string
  var err error
  absPath, err = filepath.Abs(absRoot+"/"+urlPath)
  if err != nil {
    log.Println("Error generating path", err)
  }
  log.Println("Processing absPath, url", absPath, urlPath)
  finfo, err := os.Stat(absPath)
  if err != nil {
    fmt.Fprintf(w, "Error trying os.Stat(%s) : %q", absPath, err)
    return
  }
  if finfo.IsDir() {
    //mdwiki.FmtDir(w, path)
    mdwiki.HTMLDirList(w, absPath, urlPath)
  } else {
    mdwiki.PrintFile(w, absPath, urlPath)
  }
}
