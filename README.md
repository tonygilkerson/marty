# Marty

This is an IoT project to monitor for mail delivery and watch for cars passing by on my street.

## Doc

This repo uses [mkdocs](https://www.mkdocs.org/) ([help](https://mkdocs.readthedocs.io/en/0.10/)) and [github pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) to host content at:

[https://tonygilkerson.github.io/marty/](https://tonygilkerson.github.io/marty/)

**Develop:**

```sh
mkdocs serve
# Edit content and review changes here:
open http://127.0.0.1:8000/
```

DEVTODO - use [mapfs](https://pkg.go.dev/testing/fstest) to embed config files into the running process
DEVTODO - need to refine power supply, need 5v and 3v3, ideally I would do it all from a USB-C to battery pack.
DEVTODO - need short antenna lead


## Lora Machine Workaround

2023 Apr 27

* The lora machine in TinyGo v27 works but does not have UART2
* The lora machine in TinyFo v28 has UART2 but gives "panic: runtime error at 0x08002157: heap alloc in interrupt"

So I am going to use the `board_lorae5.go` from v28 in v27

```sh
# backup v27 
cp /usr/local/Cellar/tinygo/0.27.0/src/machine/board_lorae5.go /usr/local/Cellar/tinygo/0.27.0/src/machine/board_lorae5.go.orig  

# copy v28 into v27
cp /Users/tgilkerson/github/tinygo/src/machine/board_lorae5.go /usr/local/Cellar/tinygo/0.27.0/src/machine/board_lorae5.go

# Flash with v27
/usr/local/bin/tinygo flash -target=lorae5
```