# go-testing-examples

The purpose of this project is to experiment with various Go test libraries.

The main branch uses core Go libraries (or as close to as possible) with other branches added to test specific 
libraries.  The two main aspects for comparison are differences in source and differences in output.  Source differences
may be compared via git diff - the intent is to maintain a single commit in each branch from main for this purpose.  
Output differences may be compared by running the `reports.sh` script which will run `go test -v` in each branch and copy
the output to files in the reports directory which can then be compared.
