# Go with net/http package
- [net/http package](https://pkg.go.dev/net/http)
- Go言語でWeb Serverを立てるための勉強
- Go언어로 Web Server를 구축하기 위한 공부
- net packageの配下にhttp packageがある。

## TCP Server
- HTTPはTCP上で動く
- IETF(Internet Engineering Task Force)でHTTP標準を定義
  - 現状主にHTTP1.1を使われている
- ServerはMethod（HTTP Method）とラウターを通してどのコードを実行させるかを決める。
- 主にHTTP1.1が使用されている: [RFC7230](https://www.rfc-editor.org/rfc/rfc7230#section-3.1.2) 

## What is Mux
厳密には違うがmux, servemux, router, server, http mux等々同じ意味である。

## net/http
### Handler interface
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
- polymorphismによって`ServeHTTP(ResponseWriter, *Request)`を満足するメソッドならHandlerの役割を果たせる
### Server
http.ListenAndServe
```go
func ListenAndServe(addr string, handler Handler) error
```
http.ListenAndServeTLS(https)
- TLS（Transport Layer Security）は、SSLをもとに標準化させたもの
```go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```
### Request
- go docの[request](https://pkg.go.dev/net/http#Request)
### Response
- go docの[response](https://pkg.go.dev/net/http#Response)
### ResponseWriter
- https://pkg.go.dev/net/http#ResponseWriter
