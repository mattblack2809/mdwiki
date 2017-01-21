package mdwiki

import (
	"fmt"
	//"path/filepath"
	"io"
	"io/ioutil"
	"log"
	"strings"
	//"os"
	"html/template"
	"time"
)

// dirDisplay filled out be <func> and used during running the HTML template
// to display a directory
type dirDisplay struct {
	DisplayName string
	Anchor      string
	Size        string
	Modified    string
	Accessed    string
}

// FilteredReadDir removes files named .md.html from the FileInfo slice
// where there is a coresponding .md file
func filteredReadDir(path string) []dirDisplay {
	// Where a .md and an accompanying .md.html exist
	// suppress the .md.html output
	var retList []dirDisplay
	dlist, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	// remember you can't futz with a map whilst iterating it
	dmap := make(map[string]bool) // key on file name, value immaterial
	var mds []string
	var ads []string
	for _, file := range dlist {
		fname := file.Name()
		dmap[fname] = true
		if strings.HasSuffix(fname, ".md") {
			mds = append(mds, fname)
		}
		if strings.HasSuffix(fname, ".ad") {
			ads = append(mds, fname)
		}
	}
	// remove the entries for .md.html files where there is a .md source
	for _, md := range mds {
		delete(dmap, md+".html") // deleting non-existent key is OK
	}
	for _, ad := range ads {
		delete(dmap, ad+".html") // deleting non-existent key is OK
	}
	for _, file := range dlist {
		if _, ok := dmap[file.Name()]; ok {
			var display dirDisplay
			display.DisplayName = file.Name()
			display.Size = fmt.Sprintf("%d", file.Size())
			display.Modified = tfmt(file.ModTime())
			display.Accessed = accessTime(file)
			if file.IsDir() {
				display.Anchor = file.Name() + "/" // tell browser to prepend relative URL
			} else {
				display.Anchor = file.Name()
			}
			retList = append(retList, display)
		}
	}
	return retList
}

// dirReport is the compiled version of a static HTML template used to
// output directory listings.  It registers functions that are called
// within the template.
var dirReport = template.Must(template.New("dirlist").Funcs(
	template.FuncMap{
		"ShortenName": ShortenName,
	}).Parse(dirTempl))

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
<td><a href='{{.Anchor}}'>{{.DisplayName | ShortenName}}</a></td>
<td>{{.Size}}</td>
<td>{{.Modified}}</td>
<td>{{.Accessed}}</td>
</tr>
{{end}}
</table>
`

// tfmt uses the 'canonical' base time example in Go to format a time.Time
// as a string
func tfmt(t time.Time) string {
	//return t.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	return t.Format("2006 01 02 (Mon) 15:04:05")
}

// ShortenName stips any .md.html/ad.html or .md/.ad sufffix from the passed name
func ShortenName(name string) string {
	if strings.HasSuffix(name, ".md.html") { // hopefully already stripped
		return name[:len(name)-8]
	} else if strings.HasSuffix(name, ".md") {
		return name[:len(name)-3]
	} else if strings.HasSuffix(name, ".ad.html") { // hopefully already stripped
			return name[:len(name)-8]
		} else if strings.HasSuffix(name, ".ad") {
			return name[:len(name)-3]
		} else {
		return name
	}
}

// HTMLDirList lists directory content.  Where a .md and an accompanying .md.html exist
// suppress the .md.html output: clicking on the .md link will cause any
// cached .md.html to be returned provided that file is not stale (in which
// case it is re-generted at the point of access)
func HTMLDirList(w io.Writer, absPath string, urlPath string) {

	log.Println("printing directory", absPath)
	printHTMLHeader(w)
	fmt.Fprintln(w, PrintPath(urlPath))
	defer printHTMLFooter(w)
	dlist := filteredReadDir(absPath)
	err := dirReport.Execute(w, dlist)
	if err != nil {
		fmt.Fprintf(w, "Error executing html template", err)
	}
}
