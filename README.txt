goload
======

Simple upload service using cURL or form.


BUILD:
  go build goload.go

RUN:
  ./goload

TODO:
  * Random file names,
  * Respond with file name,
  * Hash files and check for duplicates during upload (return old file if it
    already exist)
  * Probably integrate with Redis for keeping the hashes of each file (?)
