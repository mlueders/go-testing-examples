# go-testing-examples

The purpose of this project is to experiment with various Go test libraries.

The main branch uses core Go libraries (or as close to as possible) with other branches added to test specific 
libraries.  The two main aspects for comparison are differences in source and differences in output.  Source differences
may be compared via git diff - the intent is to maintain a single commit in each branch from main for this purpose.  
Output differences may be compared by running the `reports.sh` script which will run `go test -v` in each branch and copy
the output to files in the reports directory which can then be compared.

Libraries not evaluated and why
* gopkg.in/check.v1 (https://labix.org/gocheck) - no facility for subtests which would require a significant restructure
  of the main branch; this, combined with the age of the library and lack of struct comparisons exclude this lib
* gopwn (https://github.com/ToQoz/gopwt) - i couldn't get this to work... i so very badly want this to work... not sure
  what's going on here, wherever `gopwt.Empower()` is executed, output is 
  `could not import github.com/ToQoz/gopwt (can't find import: "github.com/ToQoz/gopwt")`
  no idea what's going on but i really want this to work so will probably come back to this at some point
