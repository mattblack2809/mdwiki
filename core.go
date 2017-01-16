package mdwiki

import (
  "fmt"
  "strings"
  "io"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "bufio"
)

var Options = make(map[string]string)
func ReadOptions() {
  file, err :=os.Open("mdwiki.conf")
  if err != nil {
      return
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      fields := strings.Split(scanner.Text(), ",")
      Options[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
  }
}

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

func PrintFile(w io.Writer, path string) {
  path = mdToHTML(path)
  f, err := os.Open(path)
  if err != nil {
    PrintHTMLHeader(w)
    fmt.Fprintln(w, PrintPath(path))
    fmt.Fprintf(w, "Error opening file at path %s: %q\n", path, err)
    log.Printf("Error opening file at path %s: %q\n", path, err)
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
      log.Printf("Error reading file at path %s: %q\n", path, err)
      PrintHTMLFooter(w)
      return
    }
    s := string(d) // horrible but easy
    idx := strings.Index(s, "<body>")
    if idx == -1 { // a fragment
      log.Println("printing HTML fragment at path", path)
      PrintHTMLHeader(w)
      fmt.Fprintln(w, PrintPath(path))
      fmt.Fprint(w, s)
      PrintHTMLFooter(w)
      return
    } else { // A complete web page: inject the path
      log.Println("printing complete HTML file at path", path)
      fmt.Fprint(w, s[:idx+6])
      fmt.Fprintln(w, PrintPath(path))
      fmt.Fprint(w, s[idx+6:])
      return
    }
  } else {
    log.Println("printing non-HTML file", path)
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
func mdToHTML(path string) string {
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
    if Options["no-cache"] != "true" &&
        htmlStat.ModTime().After(stat.ModTime()) {
      log.Println("using path to cached file", htmlPath)
      return htmlPath
    }
  }
  //var cmd *exec.Cmd
  if Options["toc"] == "true" {
    log.Println("option -toc set: generating output file", path)
    cmd = exec.Command("pandoc", path, "-s", "--toc", "--toc-depth=6",
      "-o", path+".html")
  } else if opts, ok := Options["pandoc-args"]; ok {
      execArgs := []string{path, "-o", path+".html"}
      fields := strings.Fields(opts)
      execArgs = append(execArgs, fields...)
      log.Printf("pandoc-args set, output to path %s with arggs\n",
        path+".html", execArgs)
      cmd = exec.Command("pandoc", execArgs...)
  } else {
    log.Println("generating output file", path)
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
