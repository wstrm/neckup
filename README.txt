neckup - For autist neckbeards (TM)
===================================

Simple upload service using cURL or form.

CONTRIBUTE:
  Please, please, please teach me more about Go.
  
  This is my first project written in the language and
  would love criticism! 

  It's obviously crappy code in some peoples eyes.

BUILD:
  go build

TEST:
  go test

RUN:
  ./neckup

HELP:
  ./neckup --help

SETUP:
  This setup uses a neckup user where everything is placed
  in its home directory (/home/neckup/neckup/).

  You should already have setup a Nginx web server.

    * Create and setup the neckup user,
        - $ useradd -m neckup
        - $ su neckup
        - $ cd ~

    * Get the latest version of neckup and cd into it,
        - $ git clone git@github.com:willeponken/neckup.git
        - $ cd neckup

    * Compile neckup.go and show the different flags available,
        - $ go build neckup.go
        - $ ./neckup --help

    * Create seperated or merged nginx server block(s),
        - see examples/nginx/neckup_*
    
    * Optionally add an init script for the process.
      Feel free to add more scripts in "examples/",
        - see examples/upstart/neckup_*.conf (upstart)

    * Reload Nginx and start neckup and you should be good to go!

DEMO:
  * Internet: https://nup.pw/
  * Hyperboria: http://h.nup.pw/
    - No ICANN: http://[fcd7:220b:18fe:ac59:836e:1777:5ec4:7d47]/ (note: the
      web server still uses the ICANN domain)

TODO:
  See /TODO.txt

LICENSE:
  GPL-3.0 (can be found at /LICENSE)
