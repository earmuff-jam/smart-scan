
Duplicate Finder — 10-Step TODO

 1. Set up the Go project

Create the Go module
Create main.go
Make dupfinder <directory> work
 2. Find all files

Recursively walk the directory
Collect each file's path and size
 3. Group files by size

Files with different sizes can't be duplicates
Only keep groups where 2+ files have the same size
 4. Hash candidate files

Use SHA-256
Read files without loading the entire file into memory
 5. Group by hash

Same size + same hash = duplicate
Print the duplicate groups
 6. Make the output nice

Show file names, sizes, number of groups, and wasted space
Sort the results
 7. Add tests + safety

Test identical files, same-size/different-content files, empty files, errors, etc.
No deleting yet
 8. Add interactive deletion

Let me choose which copy to keep
Require confirmation before deleting
Show how much space was recovered
 9. Make it fast

Add concurrent hashing with a worker pool
Add progress reporting
Benchmark it
 10. Add one "cool feature"

Pick one:
Persistent hash cache
TUI interface
Near-duplicate image detection
JSON output
Smart "which copy should I keep?" suggestions