# Goweb
--
### Web I/O package in pure golang.

### NOTE: Backwards compatibility to Golang 1.4 is discontinued post 1.0.0-alpha.1

Introduction
--
Wraps net/http into higer level abstractions that handle the following:

* Redirects
* Proxies
* Cookies
* Relative URLs
* JSON structures
* DOM parsing

The goals of Goweb are to be accurate, minimize overhead, and be insanely fast.

Dependencies
--

* [golog](http://github.com/mlavergn/golog)

Installation
--
```bash
	go get github/mlavergn/goweb
```


Examples
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

License
--
The [MIT License](http://choosealicense.com/licenses/mit/)
