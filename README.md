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
* The lora machine on the dev branch has the fixes I need. 

So I am going to use the TinyGo dev branch and point to the dev version of the devices as a workaround
