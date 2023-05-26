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
- nilを入れるとDefaultServeMuxが使用される
```go
var c hotcat
var d hotdog

http.Handle("/cat", c)
http.Handle("/dog", d)

// nilを入れるとdefault serve muxが使用される。
http.ListenAndServe(":8080", nil)
```
http.ListenAndServeTLS(https)
- TLS（Transport Layer Security）は、SSLをもとに標準化させたもの
```go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```
### Request
- go docの[request](https://pkg.go.dev/net/http#Request)
```go
type Request struct {
Method string // http methods
URL *url.URL
//	Header = map[string][]string{
//		"Accept-Encoding": {"gzip, deflate"},
//		"Accept-Language": {"en-us"},
//		"Foo": {"Bar", "two"},
//	}
Header Header
Body io.ReadCloser
ContentLength int64
Host string
// 先に`req.ParseForm()`を呼ぶ必要がある
Form url.Values
// 先に`req.ParseForm()`を呼ぶ必要がある
PostForm url.Values
}
```
- `req.Form`, `req.PostForm`を使用するためには先に`req.ParseForm()`を呼ぶ必要がある

### Response
- go docの[response](https://pkg.go.dev/net/http#Response)
```go
type ResponseWriter interface {
    // HeaderはWriteHeaderで送るHeader Mapを返す。
    Header() Header
	
    // Write は、HTTP 応答の一部としてコネクションにデータを書き込みます。
	// WriteHeader がまだコールされていない場合は、Writeはデータを書き込む前に WriteHeader(http.StatusOK)を呼びます。
    // ヘッダーに Content-Type 行が含まれていない場合、Write は自動的にContent-Typeを入れる
    Write([]byte) (int, error)
	
	// WriteHeader は、ステータス コードを含む HTTP 応答ヘッダーを送信します。
	// WriteHeader への明示的なコールは主にエラーコードを送る時に使用されます。
    WriteHeader(int)
}
```

### ResponseWriter
- https://pkg.go.dev/net/http#ResponseWriter

# Serve File
`io.Copy`, `http.ServeContent`, `http.ServeFile`, `http.FileServer`がある。
## io.Copy
```go
func dogPic(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("toby.jpg")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	defer file.Close()

	io.Copy(w, file)
}
```