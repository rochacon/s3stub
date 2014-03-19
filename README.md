s3stub
======

`s3stub` is a simple (and unsafe) blob storage API.

You can use it as an AWS S3 test server if you only use simple operations such as read, write or delete blobs.

Install
-------

```
go get github.com/rochacon/s3stub
```

Run
---

```
% s3stub -h
Usage of s3stub:
  -b="127.0.0.1:8000": The address to bind to
  -r="": The root path of the server
```

If you have Docker installed:

```bash
% docker pull rochacon/s3stub
% docker run -d -p 8000:80 rochacon/s3stub
```

Usage
-----

To write/update a blob, make a `PUT` request:

```bash
% curl -X PUT -d "new file content" 127.0.0.1:8000/filename
0eb88758c79815e61f7c3304ea43340e34773afb8b8edf561a26a40dc36fec2c
```

The SHA-256 hash of the file is returned for integrity check.


To retrieve a blob, make a `GET` request:

```bash
% curl 127.0.0.1:8000/filename
new file content
```


To delete a file, make a `DELETE` request:

```bash
% curl -i -X DELETE 127.0.0.1:8000/filename
HTTP/1.1 204 No Content
Date: Wed, 19 Mar 2014 06:11:40 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```


For both `GET` and `DELETE`, if a file is not found an `HTTP 404 Not Found` response is returned:


```bash
% curl -i 127.0.0.1:8000/nooo
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
Date: Wed, 19 Mar 2014 06:00:03 GMT
Content-Length: 49

open /tmp/s3stub/nooo: no such file or directory
```
