package mdwiki

import (
  "fmt"
  "path/filepath"
  "strings"
  "io"
  "io/ioutil"
  "log"
  "os"
  "html/template"
  "time"
)

func FilteredReadDir(path string) ([]os.FileInfo, error) {
  // Where a .md and an accompanying .md.html exist
  // suppress the .md.html output
  var retList []os.FileInfo
  dlist, err := ioutil.ReadDir(path)
  if err != nil {
    return nil, err
  }
  // remember you can't futz with a map whilst iterating it
  dmap := make(map[string]bool) // key on file name, value immaterial
  var mds []string
  for _, file := range dlist {
    fname := file.Name()
    dmap[fname] = true
    if strings.HasSuffix(fname, ".md") {
      mds = append(mds, fname)
    }
  }
  // remove the entries for .md.html files where there is a .md source
  for _, md := range mds {
    delete(dmap, md+".html")  // deleting non-existent key is OK
  }
  for _, file := range dlist {
    if _, ok := dmap[file.Name()]; ok {
      retList = append(retList, file)
    }
  }
  return retList, nil
}

var dirReport = template.Must(template.New("dirlist").Funcs(
    template.FuncMap{
      "ShortenName": ShortenName,
      "addDir": addDir,
        "tfmt": tfmt,
        "accessTime": accessTime,
    }).Parse(dirTempl))
// dirTempl designed to work with os.FileInfo
const dirTempl = `
<style type="text/css">
 table.ex1 {border-spacing: 0}
 table.ex1 td, th {padding: 0.5em 0.5em}
 table.ex1 tr:nth-child(odd) {color: #000; background: #FFF}
 table.ex1 tr:nth-child(even) {color: #000; background: #F4F4F4}

 table.ex2 {border-spacing: 0}
 table.ex2 td, th {padding: 0 0.2em}
 table.ex2 col:first-child {background: #FF0}
 table.ex2 col:nth-child(2n+3) {background: #CCC}
</style>
<h1>Directory Listing</h1>
<table  class=ex1>
<tr style='text-align: left'>
  <th>Filename</th><th>Size</th><th>Modified</th><th>Last Acccess</th>
</tr>
{{range .}}
<tr>
<td><a href='{{.Name | addDir}}'>{{.Name | ShortenName}}</a></td>
<td>{{.Size}}</td>
<td>{{.ModTime | tfmt}}</td>
<td>{{. | accessTime}}</td>
</tr>
{{end}}
</table>
`

func tfmt(t time.Time) string {
  //return t.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
  return t.Format("2006 01 02 (Mon) 15:04:05")
}

var dir string // prepend to filename
func addDir(fname string) string {
  return dir+"/"+fname // put the directory name in front
}
func ShortenName(name string) string {
  //fmt.Println(name)
  if strings.HasSuffix(name, ".md.html")  {  // hopefully already stripped
   return name[:len(name) - 8]
  } else if strings.HasSuffix(name, ".md") {
    return name[:len(name) - 3]
  } else {
    return name
  }
}

func HTMLDirList(w io.Writer, path string) {
  // List directory content.  Where a .md and an accompanying .md.html exist
  // suppress the .md.html output: clicking on the .md link will cause any
  // cached .md.html to be returned provided that file is not stale (in which
  // case it is re-generted at the point of access)

  log.Println("printing directory", path)
  PrintHTMLHeader(w)
  fmt.Fprintln(w, PrintPath(path))
  defer PrintHTMLFooter(w)

  dlist, err := FilteredReadDir(path)
  if err != nil {
    fmt.Fprintf(w, "Error reading filtered directory listing %q\n", err)
    return
  }
  // want just the name of the directory in the anchor so it reads
  // dirname/filename (not the full path as the browser does path stuff too)
  _, dir = filepath.Split(path) // get directory name - put in package var

  // should sort alpha (or whatever) - for now use dlist again
  err = dirReport.Execute(w, dlist)
  if err != nil {
    fmt.Fprintf(w, "Error executing html template", err)
  }
//  for _, file := range dlist {
//    fName := file.Name()
//    displayName := fName
//    if strings.HasSuffix(fName, ".md.html")  {  // hopefully already stripped
//      displayName = fName[:len(fName) - 8]
//    } else if strings.HasSuffix(fName, ".md") {
//      displayName = fName[:len(fName) - 3]
//    }
//    fmt.Fprintf(w, "<a href=\"%s\"> %s</a>\n", dir+"/"+fName, displayName)
//  }
}

// superceded by HTMLDirList
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
