Small Go example showcasing a bug with envoy v1.24.0 or earlier:

When trying to send ImmediateResponse on chunked response body with processing mode PARTIAL_BUFFERED, envoy sends the 
immediate response body with the original server response's header.

- This does not happen with processing mode BUFFERED
- This does happen with STREAMED as well, but STREAMED does not guarantee ImmediateResponse to be properly sent back 
  to the client (may append or reset instead)

Usage:

```shell
make docker-build
make run
curl "localhost:10000/chunked?size=1024&chunkedsize=512" -v
curl "localhost:10000/normal?size=1024" -v
make stop
```

Expected result response for both:
```
*   Trying 127.0.0.1:10000...
* Connected to localhost (127.0.0.1) port 10000 (#0)
> GET /normal?size=1024 HTTP/1.1
> Host: localhost:10000
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< content-length: 23
< content-type: text/plain
< date: Tue, 25 Oct 2022 11:53:56 GMT
< server: envoy
< 
* Connection #0 to host localhost left intact
Immediate Response Body%
```
And yet, chunked response results in:
```
*   Trying 127.0.0.1:10000...
* Connected to localhost (127.0.0.1) port 10000 (#0)
> GET /chunked?size=1024&chunkedsize=512 HTTP/1.1
> Host: localhost:10000
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< x-content-type-options: nosniff
< date: Tue, 25 Oct 2022 11:56:59 GMT
< x-envoy-upstream-service-time: 0
< server: envoy
< transfer-encoding: chunked
< 
* transfer closed with outstanding read data remaining
* Closing connection 0
curl: (18) transfer closed with outstanding read data remaining
```