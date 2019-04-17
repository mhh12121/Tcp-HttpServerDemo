## SAMPLE DEMO

#### Description
This demo is for Learning Basic RPC communication concepts.




#### How to USE
1. go run tcp_server.go
2. go run http_server.go
3. Enter localhost:8089/login in browser

#### Testing Environment
1. 
### Problems
1. I set unix file open limit up to 4096(default is 1024) so that I can test under 2000 concurrent request. But what if there comes to more requests?
2. Maybe it's the same question as above: it's mother fucking slow!!! Under 2000 concurrent(Login without redis), it gave me around 600 QPS




