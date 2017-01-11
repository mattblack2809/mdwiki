# MD Wiki

Provide browser access to a directory tree of markdown documents
using pandoc to render markdown as HTML

Features:

* Listens on port 8000 or the port specified using -port=nnnn
* File types handled as follows:
    * .md files are converted to HTML (and cached)
    * Setting -toc=true instructs pandoc to output a fully structered HTML
    document with a table of contents - useful if the documents are not well
    indexed
    * HTML fragments are output (wrapped with a header, the path, and footer)
    * HTML fully structured files are output with the path injected at the start of the body
    * Other file types are output to the browser 'as is' with no extra path
* Displays the path from the root (cwd) at the top of the page, with
each level being a link (HTML anchor)
* Clicking on any directory displayed in the path outputs a directory listing
    * Any _.md_ files are shown with the _.md_ suffix suppressed, and any (cached)
  HTML files named _.md.html_ are not listed (provided there is a
  corresponding _.md_ file)
* A URL requesting _file.md_ outputs the pandoc conversion to HTML.  The HTML
file is cached named _file.md.html_
    * The timestamps of the _.md_ and _.md.html_ files are compared - stale
    _.md.html_ files are re-generated
    * Note: Pandoc expects 4 spaces to generate a sublist - annoying...
    * Also Note: pandoc -f markdown_github generates &lt;br> where the source
    .md file has a newline - also annoying
* Note that links between the .md files can (should) remain as links to .md and not to .md.html

TODO

* Prettify directory listings
* Git integration - for example to
    * Enable display of git log information associated with the file
    * Potentially enable URL to specify a git tag to view the tagged version
* Open a shell at the relevant directory (or even invoke editor on the
     .md file in that shell: web based editing - let's not go there!)