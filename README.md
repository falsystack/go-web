# Go with net/http package
- [net/http package](https://pkg.go.dev/net/http)
- Go言語でWeb Serverを立てるための勉強
- net packageの配下にhttp packageがある。
- Networkの深い理解が必要。
- Request -> RFC7230（HTTP Protocol） -> HTTP Message  

## TCP Server
- IETF(Internet Engineering Task Force)でHTTP標準を定義
  - 現状主にHTTP1.1を使われている
- ServerはMethod（HTTP Method）とラウターを通してどのコードを実行させるかを決める。
- 主にHTTP1.1が使用されている: [RFC7230](https://www.rfc-editor.org/rfc/rfc7230#section-3.1.2) 

## What is Mux
厳密には違うがmux, servemux, router, server, http mux, multiplexer(電気の経路を決めるのに使う装備)等々同じ意味である。

## net
netパッケージを用いてTCPサーバーを立てる
`net.Listen`
```go
li, err := net.Listen("tcp", ":8080")
conn, err := li.Accept()
scanner := bufio.NewScanner(conn)
for scanner.Scan() {
    ln := scanner.Text()
    fmt.Println(ln)
}
```

`net.Dial`
```go
conn, err := net.Dial("tcp", "localhost:8080")
bs, err := io.ReadAll(conn)
fmt.Println(string(bs))
```

connectionのdead line設定方法
```go
err := conn.SetDeadline(time.Now().Add(10 * time.Second))
```

## net/http
### Handler interface
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
- polymorphismによって`ServeHTTP(http.ResponseWriter, *http.Request)`を具現するメソッドならHandlerの役割を果たせる
### Server
http.ListenAndServe
- ListenAndServe内部的に`net.Listen`が使われている
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
### *http.Request
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
- `req.Form`, `req.PostForm`を使用するためには先に`req.ParseForm()`を呼ぶ必要がある, `req.ParseForm()`を呼ぶと`req.Form`を更新してくれる
- `req.Form` : mapタイプ

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

### ServeMux
**基本使用方法**

```go
type hotdog, hotcat int

func (d hotdog) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, "dog dog dog")
}

func (c hotcat) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, "cat cat cat")
}

func main() {
  mux := http.NewServeMux()
  mux.Handle("/cat", c)
  mux.Handle("/dog/", d)
  
  http.ListenAndServe(":8080", mux)
}
```
**DefaultServeMuxの使用**
- ListenAndServeにnilを渡すとDefaultServeMuxが使用される。
```go
type hotdog, hotcat int

func (d hotdog) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog dog dog")
}

func (c hotcat) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "cat cat cat")
}

func main() {
	var c hotcat
	var d hotdog

	http.Handle("/cat", c)
	http.Handle("/dog", d)

	// nilを入れるとdefault serve muxが使用される。
	http.ListenAndServe(":8080", nil)
}
```

**HandleFuncの使用**
```go
// ServeHTTPではない
func d(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog dog dog")
}

func c(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "cat cat cat")
}

func main() {
	http.HandleFunc("/dog", d)
	http.HandleFunc("/cat", c)

	// use default serve mux
	http.ListenAndServe(":8080", nil)
}
```

**HandlerFuncの使用**
一番多く使われる。
```go
func d(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog dog dog")
}

func c(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "cat cat cat")
}

func main() {
    // handleを使用している
	// http.HandlerFunc()でタイプをconversionしている
	http.Handle("/cat", http.HandlerFunc(c)) 
	http.Handle("/dog", http.HandlerFunc(d))

	http.ListenAndServe(":8080", nil)
}
```

# Serve File
`io.Copy`, `http.ServeContent`, `http.ServeFile`, `http.FileServer`がある。
## io.Copy
単一ファイルを提供
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

## ServeContent
単一ファイルを提供
```go
func dogPic(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("toby.jpg")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}

	http.ServeContent(w, req, file.Name(), info.ModTime(), file)
}
```

## ServeFile
単一ファイルを提供
cacheされたFileがある場合関数が呼ばれない
```go
func dogPic(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[dogPic] serving picture")
	http.ServeFile(w, req, "toby.jpg")
}
```

## FileServer
- directoryを指定して提供できる。
- FileServerはHnadlerを返す。

```go
http.Handle("/", http.FileServer(http.Dir(".")))
```

## StripPrefix
StripPrefixはHandlerを返す。
```go
func main() {
	http.HandleFunc("/", dog)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8080", nil)
}

func dog(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 指定されたprefixを除去して残りを元にファイルを探す, /resources/toby.jpg -> toby.jpg 
	io.WriteString(w, `<img src="/resources/toby.jpg">`)
}
```

## index.html
特別にindex.htmlがある場合index.htmlだけが提供される。
```go
// 現在のディレクトリを全部提供しているがそのディレクトリにindex.htmlがある場合
// index.htmlだけが提供される。
log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("."))))
```

# log.Fatal & http error
## log.Fatal
log出力してプログラムを終了させる。
```go
log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("."))))
```

## http.Error
- エラーは、指定されたエラー メッセージと HTTP コードでリクエストに応答します。 
- それ以外の場合はリクエストは終了しません。 
- 呼び出し元は、w への書き込みがそれ以上行われないようにする必要があります。 
- エラー メッセージはプレーンテキストである必要があります。

```go
// http.StatusNotFound -> 404
http.Error(w, "file not found", http.StatusNotFound)
```

## NotFoundHandler
http.StatusNotFoundを返す単純なハンドラー
```go
// chromeだとfaviconを要求する、firefoxはしない
http.Handle("/favicon.ico", http.NotFoundHandler())
```

## Passing Data
get -> url
post -> body
### Query parameter
`req.FormValue(ket string) string`は
- キーが存在しない場合空文字列を返す。
- URLのクエリパラメタよりHTML FormのPOST / PUTのBodyパラメタが優先される。
- 同じキーに複数の値を持つ場合、ParseFormを呼び出し後Request.Formを調べる。
```go
func foo(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	fmt.Fprintln(w, "Do my Search: "+query)
}


func foo(w http.ResponseWriter, req *http.Request) {
    if err := req.ParseForm(); err != nil {
        log.Println(err)
    }
    
    // 同じキーに値が複数ある場合
    form := req.Form // map[string][]string
    values := form["q"]
    fmt.Println(values[0], values[1])
}

```

### Form
form valueも`FormValue()`関数一つで取れる、JavaのSpringと一緒でURL側とForm側両方一つで取れる
- urlもbodyもkey=value形式だから当然だけど、なぜかnodejsは違った記憶がある。

```go
func foo(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `
		<form method="get">
			<input type="text" name="q">
			<input type="submit" value="send">
		</form>
		<br>
	`+query)
}
```
### Upload
fileもfile headerとbodyが存在する
- `req.FormFile(key string)`でファイルをOpen
- `io.ReadAll(file)`でファイルを読む(Read)
```go
// open
file, header, err := req.FormFile("q")
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
defer file.Close()

log.Println("\n file: ", file, "\n header: ", header, "\n err: ", err)

// read
bs, err := io.ReadAll(file)
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
s = string(bs)
```