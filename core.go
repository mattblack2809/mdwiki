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
		if elem != "" {
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

func FmtDir(path string) string {
  // List directory content.  Where a .md and an accompanying .md.html exist
  // suppress the .md.html output: clicking on the .md link will cause any
  // cached .md.html to be returned provided that file is not stale (in which
  // case it is re-generted at the point of access)
  res := ""

  dlist, err := ioutil.ReadDir(path)
  if err != nil {
    res = fmt.Sprintf("Error readind directory listing %q\n", err)
    return res
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
      res += fmt.Sprintf("<a href=\"%s\"> %s</a>\n",
        dir+"/"+file.Name(), f)
    }
  }
  return res
}

func PrintFile(w io.Writer, path string) {
  path = mdToHTML(path)
  f, err := os.Open(path)
  if err != nil {
    fmt.Fprintf(w, "Error opening file at path %s: %q\n", path, err)
    return
  }
  defer f.Close()

  _, err = io.Copy(w, f)
  if err != nil {
    fmt.Fprintf(w, "Error copying README.md %q\n", err)
    return
  }
}

// mdToHTML does the following:
// - return if path does not end .md
// - check if html version exists with timestamp >= than that of the .md file
// - return path to valid existing cached html file; otherwise
// - generate html file and return path to that
func mdToHTML(path string) string {
  // TODO allow valid extensions additional to .md
  if path[len(path)-3:] != ".md" {
    return path
  }
  mdStat, err := os.Stat(path)
  if err != nil {
    log.Printf("Error tying os.Stat(%s) : %q", path, err)
    return path
  }
  htmlStat, err := os.Stat(path+".html")
  if err == nil { // file exists then
    if htmlStat.ModTime().After(mdStat.ModTime()) {
      return path+".html"
    }
  }
//  cmd := exec.Command("pandoc", path, "-f", "markdown_github",
//    "-o", path+".html")
  cmd := exec.Command("pandoc", path,
    "-o", path+".html")
  err = cmd.Run()
  if err != nil {
    log.Printf("Error running pandoc on file %s : %q", path, err)
    return path
  }
  return path+".html"
}
