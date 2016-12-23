# goweb
--
### Web I/O package in pure golang.

### NOTE: Backwards compatibility to Golang 1.4 is discontinued post 1.0.0-alpha.1

Introduction
--
Wraps net/http into a higer level class that handles the following:

* Redirects
* Cookies
* Relative URLs
* JSON structures
* DOM parsing

Dependencies
--

* This package depends on golog

Installation
--
```bash
	go get github/mlavergn/goweb
```


A minimal example
--
```go

	// HTTP
	c := goweb.NewHTTP()
	c.Get("http://www.google.com")
	if c.Status() == 200 {
		print(c.Contents())
	}

	// DOM
	d := goweb.NewDOM()
	d.SetContents(c.Contents())
	nodes := d.Find("form", map[string]string{"id": "tsf"})
	if len(nodes) > 0 {
		url := nodes[0].Attr("action")
		print(url)
	}

```
