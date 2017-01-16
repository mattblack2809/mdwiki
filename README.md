# MD Wiki


#### Purpose (what is it for):
* Provides browser access to a directory tree of (primarily) markdown documents using pandoc to render markdown as HTML
* Works with a document set that may be less than perfectly cross-referenced -
by providing directory listings and optional auto-generation of a table-of-content
for each document
* Example use: mdwiki -port=8008 -root=/home/matt/docs -toc=true


#### What (does it do)
* A clickable path is displayed at the top of the page as part of any html
formatted output
* HTML fragments are output wrapped with a header, the (generated clickable)
  path, and footer, while HTML fully structured files are output with the path
  injected at the start of the body
  * .md files display without the .md, and when clicked generate HTML output
* Non-html files are displayed 'as-is' with no injected path information
* The access time stamp is updated per view, so that the user sees the last
  modified and last accessed time for each file

#### How (does it do it)
* Options may be specified one per line in mdwiki.conf (read from cwd), and
  overridden using command line options
    * run mdwiki -help for list of options
    * Be aware:
        * Use of the browser back-button may display stale directory listings
        * mdwiki does not follow symbolic links
* A URL requesting _file.md_ outputs the pandoc conversion to HTML.  The HTML
file is cached named _file.md.html_
    * Links between the .md files can (should) remain as links
     to .md and not to .md.html
    * The timestamps of the _.md_ and _.md.html_ files are compared - stale
    _.md.html_ files are re-generated
    * Note: Pandoc expects 4 spaces to generate a sublist - annoying...
    * Also Note: pandoc -f markdown_github generates &lt;br> where the source
    .md file has a newline - also annoying
* Technical notes:
    * File access time 'atime' is OS dependent and also 'broken' by design: for
    performance reasons - to avoid a disk write - both linux and Windows do not
    update atime when a file is read in most circumstances.  Reading the atime is
    OS specific - hence use of build tags to compile for linux or Windows,
    and update of atime is achieved
    by executing 'touch -a <file>' which on Windows assumes installation of a
    touch executable such as is intalled part of git bash...
    * Any anchor referring to a directory must end with a "/" character in order
    to have the web browser prepend the path from the site root

TODO

* Git integration - for example to
    * Enable display of git log information associated with the file
    * Potentially enable URL to specify a git tag to view the tagged version
* Open a shell at the relevant directory (or even invoke editor on the
     .md file in that shell: web based editing - let's not go there!)
