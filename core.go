package mdwiki

import (
  "fmt"
  "path/filepath"
  "strings"
  "io"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
)


func PrintPath(path string) string {
  res := ""
	e := strings.Split(path, "/")
	// get rid of empty elemements
	var elems []string
	for _, elem := range e {
		if elem != "" && elem != "." {
			elems = append(elems, elem)
		}
	}
	link := ""
	pp := ""
  for n, elem := range elems {
		pp += "/" + elem
		if n < len(elems) -1 {
			link += fmt.Sprintf("&nbsp;/&nbsp;<a href=\"%s\">%s</a>", pp, elem)
		} else {
			link += fmt.Sprintf("&nbsp;/&nbsp;<strong>%s</strong>", elem)
		}
	}
	if link == "" {
		link = "<strong>root</strong>"
	} else {
		link = "<a href=\"/\">root</a>" + link
	}
	res = link + "<br><hr>"
  return res
}

func FmtDir(w io.Writer, path string) {
  // List directory content.  Where a .md and an accompanying .md.html exist
  // suppress the .md.html output: clicking on the .md link will cause any
  // cached .md.html to be returned provided that file is not stale (in which
  // case it is re-generted at the point of access)

  PrintHTMLHeader(w)
  fmt.Fprintln(w, PrintPath(path))
  defer PrintHTMLFooter(w)

  dlist, err := ioutil.ReadDir(path)
  if err != nil {
    fmt.Fprintf(w, "Error readind directory listing %q\n", err)
    return
  }
  // want just the name of the directory in the anchor so it reads
  // dirname/filename (not the full path as the browser does path stuff too)
  _, dir := filepath.Split(path) // get directory name
  // remember you can't futz with a map whilst iterating it
  dmap := make(map[string]string) // key on file name, value the display name
  var mds []string
  for _, file := range dlist {
    fname := file.Name()
    if strings.HasSuffix(fname, ".md.html") {
      dmap[fname] = fname[:len(fname) - 8]
    } else if strings.HasSuffix(fname, ".md") {
      mds = append(mds, fname)
      dmap[fname] = fname[:len(fname) - 3]
    } else {
      dmap[fname] = fname
    }
  }
  // remove the entries for .md.html files where there is a .md source
  for _, md := range mds {
    delete(dmap, md+".html")
  }

  // should sort alpha (or whatever) - for now use dlist again

  for _, file := range dlist {
    f, ok := dmap[file.Name()]
    if ok {
      fmt.Fprintf(w, "<a href=\"%s\"> %s</a>\n",
        dir+"/"+file.Name(), f)
    }
  }
}

func PrintFile(w io.Writer, path string, toc bool) {
  path = mdToHTML(path, toc)
  f, err := os.Open(path)
  if err != nil {
    PrintHTMLHeader(w)
    fmt.Fprintln(w, PrintPath(path))
    fmt.Fprintf(w, "Error opening file at path %s: %q\n", path, err)
    PrintHTMLFooter(w)
    return
  }
  defer f.Close()
  // if the file ends .html, check if it is a complete web page
  // or a fragment - and either inject the PrintPath in to the full page
  // or surround the fragment.
  // For other file types just copy the file content without any path
  if strings.HasSuffix(path, ".html") {
    d, err := ioutil.ReadFile(path)
    if err != nil {
      PrintHTMLHeader(w)
      fmt.Fprintln(w, PrintPath(path))
      fmt.Fprintf(w, "Error reading file at path %s: %q\n", path, err)
      PrintHTMLFooter(w)
      return
    }
    s := string(d) // horrible but easy
    idx := strings.Index(s, "<body>")
    if idx == -1 { // a fragment
      PrintHTMLHeader(w)
      fmt.Fprintln(w, PrintPath(path))
      fmt.Fprint(w, s)
      PrintHTMLFooter(w)
      return
    } else { // A complete web page: inject the path
      fmt.Fprint(w, s[:idx+6])
      fmt.Fprintln(w, PrintPath(path))
      fmt.Fprint(w, s[idx+6:])
      return
    }
  } else {
    _, err = io.Copy(w, f)
    if err != nil {
      fmt.Fprintf(w, "Error copying path $s, %q\n", path, err)
    }
  }
}

// mdToHTML does the following:
// - return if path does not end .md
// - check if html version exists with timestamp >= than that of the .md file
// - return path to valid existing cached html file; otherwise
// - generate html file and return path to that
// if toc is true, generate a stand-alone html file with toc
// (pandoc won't generate a toc for a fragment)
func mdToHTML(path string, toc bool) string {
  // TODO allow valid extensions additional to .md
  stat, err := os.Stat(path)
  if err != nil {
    log.Printf("Error tying os.Stat(%s) : %q", path, err)
    return path
  }
  // Update the atime to now - of the source file - not any .md.html file.
  // On linux, reading the file does NOT update the atime unless
  // the mtime is greater than the mtime - for performance to avoid
  // disk writes.  Therefore, for this app, force an update on atime.
  cmd := exec.Command("touch", path, "-a")
  err = cmd.Run()
  if err != nil {
    log.Printf("Error updating access time on file %s : %q", path, err)
  }
  if path[len(path)-3:] != ".md" {
    return path
  }

  htmlPath := path+".html"
  htmlStat, err := os.Stat(htmlPath)
  if err == nil { // file exists then
    if htmlStat.ModTime().After(stat.ModTime()) {
      return htmlPath
    }
  }
  //var cmd *exec.Cmd
  if toc {
    cmd = exec.Command("pandoc", path, "-s", "--toc", "--toc-depth=6",
      "-o", path+".html")
  } else {
    cmd = exec.Command("pandoc", path,
      "-o", path+".html")
  }
  err = cmd.Run()
  if err != nil {
    log.Printf("Error running pandoc on file %s : %q", path, err)
    return path
  }
  return path+".html"
}


func PrintHTMLHeader(w io.Writer,) {
  fmt.Fprintln(w, "<html><head><title>MD Wiki</title></head><body>")
}

func PrintHTMLFooter(w io.Writer,) {
  fmt.Fprintln(w, "</body></html>")
}
