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

  You should already have setup a nginx web server.

    * Create and setup the neckup user,
        - $ useradd -m neckup
        - $ su neckup
        - $ cd ~

    * Get the latest version of neckup and cd into it,
        - $ git clone git@github.com:willeponken/neckup.git
        - $ cd neckup

    * Comepile neckup.go and show the different flags available,,
        - $ go build neckup.go
        - $ ./neckup --help

    * Create a nginx server block that proxies requests to
      localhost:8080 or w/e you specified,
        - see second block in examples/nginx/neckup
    
    * Create a nginx server block that has the ./files as root
      or w/e you specified,
        - see first block in examples/nginx/neckup

    * Optionally add an init script for the process.
      Feel free to add more scripts in "examples/".
        - see examples/upstart/neckup.conf (upstart)

    * Reload nginx and start neckup and you should be good to go!

DEMO:
  http://u.wiol.io

LICENSE:
  GPL-3.0 (can be found at /LICENSE)
