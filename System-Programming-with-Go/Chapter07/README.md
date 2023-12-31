# Writing Web Applications in Go

In the previous chapter, we discussed many advanced topics related to goroutines and channels as well as shared memory and mutexes.

The main subject of this chapter is the development of web applications in Go. However, this chapter will also talk about how to interact with two popular databases in your Go programs. The Go standard library provides packages that can help you develop web applications using higher level functions, which means that you can do complex things such as reading web pages by just calling a couple of Go functions with the right arguments. Although this kind of programming hides the complexity behind a request and offers less control over the details, it allows you to develop difficult applications using fewer lines of code, which also results in having fewer bugs in your programs.

However, as this book is about systems programming, this chapter will not go into too much depth: you might consider the presented information as a good starting point for anyone who wants to learn about web development in Go.

More specifically, this chapter will talk about the following topics:

-   Creating a Go utility for MySQL database administrators
-   Administering a MongoDB database
-   Using the Go MongoDB driver to talk to a MongoDB database
-   Creating a web server in Go
-   Creating a web client in Go
-   The http.ServeMux type
-   Dealing with JSON data in Go
-   The net/http package
-   The html/template Go standard package
-   Developing a command-line utility that searches web pages for a given keyword

Just Imagine

# What is a web application?

A web application is a client-server software application where the client part runs on a web browser. Web applications include webmail, instant messaging services, and online stores.

Just Imagine

# About the net/http Go package

The hero of this chapter will be the net/http package that can help you write web applications in Go. However, if you are interested in dealing with TCP/IP connections at a lower level, then you should go to [](https://subscription.imaginedevops.io/book/programming/9781787125643/12)[Chapter 12](https://subscription.imaginedevops.io/book/programming/9781787125643/12), _Network Programming_, which talks about developing TCP/IP applications using lower level function calls.

The net/http package offers a built-in web server as well as a built-in web client that are both pretty powerful. The http.Get() method can be used for making HTTP and HTTPS requests, whereas the http.ListenAndServe() function can be used for creating naive web servers by specifying the IP address and the TCP port the server will listen to, as well as the functions that will handle incoming requests.

Another very convenient package is html/template, which is part of the Go standard library and allows you to generate an HTML output using Go HTML template files.

Just Imagine

# Developing web clients

In this section, you will learn how to develop web clients in Go and how to time out a web connection that takes too long to finish.

# Fetching a single URL

In this subsection, you will learn how to read a single web page with the help of the http.Get() function, which is going to be demonstrated in the getURL.go program. The utility will be presented in four parts; the first part of the program is the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "net/http" 
   "os" 
   "path/filepath" 
) 
```

Although there is nothing new here, you might find impressive the fact that you will use Go packages that are related to file input and output operations even though you are reading data from the internet. The explanation for this is pretty simple: Go has a uniform interface for reading and writing data regardless of the medium it is in.

The second part of getURL.go has the following Go code:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Printf("Usage: %s URL\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   URL :=os.Args[1] 
   data, err := http.Get(URL) 
```

The URL you want to fetch is given as a command-line argument to the program. Additionally, you can see the call to http.Get(), which does all the dirty work! What http.Get() returns is a Response variable, which in reality is a Go structure with various properties and methods.

The third part is the following:

```markup
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } else { 
```

If there is an error after calling http.Get(), this is the place to check for it.

The fourth part contains the following Go code:

```markup
         defer data.Body.Close() 
         _, err := io.Copy(os.Stdout, data.Body) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
   } 
}
```

As you can see, the data of URL is written in standard output using os.Stdout, which is the preferred way for printing data on the screen. Additionally, the data is saved in the Body property of the return value of the http.Get() call. However, not all HTTP requests are simple. If the response streams a video or something similar, it would make sense to be able to read it one piece at a time instead of getting all of it in a single data piece. You can do that with io.Reader and the Body part of the response.

Executing getURL.go will generate the following raw results, which is what a web browser would have gotten and rendered:

```markup
$ go run getURL.go http://www.mtsoukalos.eu/ | head
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML+RDFa 1.0//EN"
  "http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd">
<htmlxmlns="http://www.w3.org/1999/xhtml" xml:lang="en" version="XHTML+RDFa 1.0" dir="ltr"
xmlns:content=http://purl.org/rss/1.0/modules/content/
. . .
</script>
</body>
</html>
```

Generally speaking, although getURL.go does the desired job, the way it works is not so sophisticated because it gives you no flexibility or a way to be creative.

# Setting a timeout

In this subsection, you will learn how to set a timeout for a http.Get() request. For reasons of simplicity, it will be based on the Go code of getURL.go. The name of the program will be timeoutHTTP.go and will be presented in five parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "net" 
   "net/http" 
   "os" 
   "path/filepath" 
   "time" 
) 
 
var timeout = time.Duration(time.Second) 
```

Here, you declare the desired timeout period, which is 1 second, as a global parameter.

The second part of timeoutHTTP.go has the following Go code:

```markup
func Timeout(network, host string) (net.Conn, error) { 
   conn, err := net.DialTimeout(network, host, timeout) 
   if err != nil { 
         return nil, err 
   } 
   conn.SetDeadline(time.Now().Add(timeout)) 
   return conn, nil 
} 
```

Here, you define two types of timeouts, the first one is defined with net.DialTimeout() and is for the time it will take your client to connect to the server. The second one is the read/write timeout, which has to do with the time you want to wait to get a response from the web server after connecting to it: this is defined with the call to the conn.SetDeadline() function.

The third part of the presented program is the following:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Printf("Usage: %s URL\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   URL :=os.Args[1] 
```

The fourth portion of the program is the following:

```markup
   t := http.Transport{ 
         Dial: Timeout, 
   } 
 
   client := http.Client{ 
         Transport: &t, 
   } 
   data, err := client.Get(URL) 
```

Here, you define the desired parameters of the connection using an http.Transport variable.

The last part of the program contains the following Go code:

```markup
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } else { 
         deferdata.Body.Close() 
         _, err := io.Copy(os.Stdout, data.Body) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
   } 
} 
```

This part of the program is all about error handling!

Executing timeoutHTTP.go will generate the following output in case of a timeout:

```markup
$ go run timeoutHTTP.go http://localhost:8001
Get http://localhost:8001: read tcp [::1]:58018->[::1]:8001: i/o timeout
exit status 100
```

The simplest way to deliberately create a timeout during a web connection is to call the time.Sleep() function in the handler function of a web server.

# Developing better web clients

Although getURL.go does the required job pretty quickly and without writing too much Go code, it is in a way not adaptable or informative. It just prints a bunch of raw HTML code without any other information and without the capability of dividing the HTML code into its logical parts. Therefore, getURL.go needs to be improved!

The name of the new utility will be webClient.go and will be presented to you in five segments of Go code.

The first part of the utility is the following:

```markup
package main 

import ( 
   "fmt" 
   "net/http" 
   "net/http/httputil" 
   "net/url" 
   "os" 
   "path/filepath" 
   "strings" 
) 
```

The second part of the Go code from webClient.go is the following:

```markup
func main() { 
   if len(os.Args) != 2 { 
         fmt.Printf("Usage: %s URL\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   URL, err :=url.Parse(os.Args[1]) 
   if err != nil { 
         fmt.Println("Parse:", err) 
         os.Exit(100) 
   } 
```

The only new thing here is the use of the url.Parse() function that creates a URL structure from a URL that is given as a string to it.

The third part of webClient.go has the following Go code:

```markup
   c := &http.Client{} 
 
   request, err := http.NewRequest("GET", URL.String(), nil) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   httpData, err := c.Do(request) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

In this Go code, you first create an http.Client variable. Then, you construct a GET HTTP request using http.NewRequest(). Last, you send the HTTP request using the Do() function, which returns the actual response data saved in the httpData variable.

The fourth code part of the utility is the following:

```markup
   fmt.Println("Status code:", httpData.Status) 
   header, _ := httputil.DumpResponse(httpData, false) 
   fmt.Print(string(header)) 
 
   contentType := httpData.Header.Get("Content-Type") 
   characterSet := strings.SplitAfter(contentType, "charset=") 
   fmt.Println("Character Set:", characterSet[1]) 
 
   if httpData.ContentLength == -1 { 
         fmt.Println("ContentLength in unknown!") 
   } else { 
         fmt.Println("ContentLength:", httpData.ContentLength) 
   } 
```

Here, you find the status code of the HTTP request using the Status property. Then, you do a little digging into the Header part of the response in order to find the character set of the response. Last, you check the value of the ContentLength property, which equals \-1 for dynamic pages: this means that you do not know the page size in advance.

The last part of the program has the following Go code:

```markup
   length := 0 
   var buffer [1024]byte 
   r := httpData.Body 
   for { 
         n, err := r.Read(buffer[0:]) 
         if err != nil { 
               fmt.Println(err) 
               break 
         } 
         length = length + n 
   } 
   fmt.Println("Response data length:", length) 
} 
```

Here, you find the length of the response by reading from the Body reader and counting its data. If you want to print the contents of the response, this is the right place to do it.

Executing webClient.go will create the following output:

```markup
$ go run webClient.go invalid
Get invalid: unsupported protocol scheme ""
exit status 100
$ go run webClient.go https://www.mtsoukalos.eu/
Get https://www.mtsoukalos.eu/: dial tcp 109.74.193.253:443: getsockopt: connection refused
exit status 100
$ go run webClient.go http://www.mtsoukalos.eu/
Status code: 200 OK
HTTP/1.1 200 OK
Accept-Ranges: bytes
Age: 0
Cache-Control: no-cache, must-revalidate
Connection: keep-alive
Content-Language: en
Content-Type: text/html; charset=utf-8
Date: Mon, 10 Jul 2017 07:29:48 GMT
Expires: Sun, 19 Nov 1978 05:00:00 GMT
Server: Apache/2.4.10 (Debian) PHP/5.6.30-0+deb8u1 mod_wsgi/4.3.0 Python/2.7.9
Vary: Accept-Encoding
Via: 1.1 varnish-v4
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Generator: Drupal 7 (http://drupal.org)
X-Powered-By: PHP/5.6.30-0+deb8u1
X-Varnish: 6922264
    
Character Set: utf-8
ContentLength in unknown!
EOF
Response data length: 50176
```

Just Imagine

# A small web server

Enough with the web clients: in this section, you will learn how to develop web servers in Go!

The Go code for the implementation of a naive web server can be found in webServer.go, and this will be presented in four parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net/http" 
   "os" 
) 
```

The second part is where things start to get tricky and strange:

```markup
func myHandler(w http.ResponseWriter, r *http.Request) { 
   fmt.Fprintf(w, "Serving: %s\n", r.URL.Path) 
   fmt.Printf("Served: %s\n", r.Host) 
} 
```

This is a kind of function that handles HTTP requests: the function takes two arguments, a http.ResponseWriter variable and a pointer to an http.Request variable. The first argument will be used for constructing the HTTP response, whereas the http.Request variable holds the details of the HTTP request that was received by the server, including the requested URL and the IP address of the client.

The third part of webServer.go has the following Go code:

```markup
func main() { 
   PORT := ":8001" 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Using default port number: ", PORT) 
   } else { 
         PORT = ":" + arguments[1] 
   } 
```

Here, you just deal with the port number of the web server: the default port number is 8001, unless there is a command-line argument.

The last chunk of Go code for webServer.go is the following:

```markup
   http.HandleFunc("/", myHandler) 
   err := http.ListenAndServe(PORT, nil) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(10) 
   } 
} 
```

The http.HandleFunc() call defines the name of the handler function (myHandler) as well as the URLs that it will support: you can call http.HandleFunc() multiple times. The current handler supports /URL, which in Go matches all URLs!

After you are done with the http.HandleFunc() calls, you are ready to call http.ListenAndServe() and start waiting for incoming connections! If you do not specify an IP address in the http.ListenAndServe() function call, then the web server will listen to all configured network interfaces of the machine.

Executing webServer.go will generate no output, unless you try to fetch some data from it: in this case, it will print logging information on your Terminal, which will show the server name (localhost) and port number (8001) of the request, as shown here:

```markup
$ go run webServer.go
Using default port number:  :8001Served: localhost:8001Served: localhost:8001
Served: localhost:8001
```

The following screenshot shows three outputs of webServer.go on a web browser:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/92cc0e7a-c289-4559-ae6f-6fc6e9d42ad3.png)

Using webServer.go

However, if you use a command-line utility such as wget(1) or getURL.go instead of a web browser, you will get the following output when you try to connect to the Go web server:

```markup
$ go run getURL.go http://localhost:8001/
Serving: /
```

The biggest advantage you get from custom made web servers is security because they are really difficult to hack when developed with security as well as easier customization in mind.

The next subsection will show how to create web servers using http.ServeMux.

# The http.ServeMux type

In this subsection, you will learn how to use the http.ServeMux type in order to improve the way your Go web server will operate. Putting it simply, http.ServeMux is a HTTP request router.

# Using http.ServeMux

The web server implementation of this section will support multiple paths with the help of http.ServeMux, which will be illustrated in the serveMux.go program that will be displayed in four parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "net/http" 
   "time" 
) 
```

The second part of serveMux.go has the following Go code:

```markup
func about(w http.ResponseWriter, r *http.Request) { 
   fmt.Fprintf(w, "This is the /about page at %s\n", r.URL.Path) 
   fmt.Printf("Served: %s\n", r.Host) 
} 
 
func cv(w http.ResponseWriter, r *http.Request) { 
   fmt.Fprintf(w, "This is the /CV page at %s\n", r.URL.Path) 
   fmt.Printf("Served: %s\n", r.Host) 
} 
 
func timeHandler(w http.ResponseWriter, r *http.Request) { 
   currentTime := time.Now().Format(time.RFC1123) 
   title := currentTime 
   Body := "The current time is:" 
   fmt.Fprintf(w, "<h1 align=\"center\">%s</h1><h2 align=\"center\">%s</h2>", Body, title) 
   fmt.Printf("Served: %s for %s\n", r.URL.Path, r.Host) 
} 
```

Here, you have the implementation of three HTTP handler functions. The first two display a static page, whereas the third one displays the current time, which is a dynamic text.

The third part of the program is the following:

```markup
func home(w http.ResponseWriter, r *http.Request) { 
   ifr.URL.Path == "/" { 
         fmt.Fprintf(w, "Welcome to my home page!\n") 
   } else { 
         fmt.Fprintf(w, "Unknown page: %s from %s\n", r.URL.Path, r.Host) 
   } 
   fmt.Printf("Served: %s for %s\n", r.URL.Path, r.Host) 
} 
```

The home() handler function will have to make sure that it is actually serving /Path, because /Path catches everything!

The last part of serveMux.go contains the following Go code:

```markup
func main() { 
   m := http.NewServeMux() 
   m.HandleFunc("/about", about) 
   m.HandleFunc("/CV", cv) 
   m.HandleFunc("/time", timeHandler) 
   m.HandleFunc("/", home) 
 
   http.ListenAndServe(":8001", m) 
} 
```

Here, you define the paths that your web server will support. Note that paths are case sensitive and that the last path in the preceding code catches everything. This means that if you put m.HandleFunc("/", home) first, you will not be able to match anything else. Putting it simply, the order of the m.HandleFunc() statements matters. Also, note that if you want to support both /about and /about/, you should have both m.HandleFunc("/about", about) and m.HandleFunc("/about/", about).

Running serveMux.go will generate the following output:

```markup
$ go run serveMux.goServed: / for localhost:8001Served: /123 for localhost:8001
Served: localhost:8001
Served: /cv for localhost:8001
```

The following screenshot shows the various kinds of outputs generated by serveMux.go on a web browser: note that the browser output is not related to the preceding output from the go run serveMux.go command:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/67d1c2e9-8fd4-4134-8ae7-f992fdb8d0d7.png)

Using serveMux.go

If you use wget(1) instead of a web browser, you will get the following output:

```markup
$ wget -qO- http://localhost:8001/CV
This is the /CV page at /CV
$ wget -qO- http://localhost:8001/cv
Unknown page: /cv from localhost:8001
$ wget -qO- http://localhost:8001/time
<h1 align="center">The current time is:</h1><h2 align="center">Mon, 10 Jul 2017 13:13:27 EEST</h2>
$ wget -qO- http://localhost:8001/time/
Unknown page: /time/ from localhost:8001
```

So, http.HandleFunc() is the default call in the library that will be used for first time implementations, whereas the HandleFunc() function of http.NewServeMux() is for everything else. Putting it simply, it is better to use the http.NewServeMux() version instead of the default one except in the simplest of cases.

Just Imagine

# The html/template package

**Templates** are mainly used for separating the formatting and data parts of the output. Note that a Go template can be either a file or string: the general idea is to use strings for smaller templates and files for bigger ones.

In this section, we will talk about the html/template package by showing an example, which can be found in the template.go file and will be presented in six parts. The general idea behind template.go is that you are reading a text file with records that you want to present in HTML format. Given that the name of the package is html/template, a better name for the program would have been genHTML.go or genTemplate.go.

There is also the text/template package, which is more useful for creating plain text output. However, you cannot import both text/template and html/template on the same Go program without taking some extra steps to disambiguate them because the two packages have the same package name (template). The key distinction between the two packages is that html/template does sanitization of the data input for HTML injection, which means that it is more secure.

The first part of the source file is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "html/template" 
   "net/http" 
   "os" 
   "strings" 
) 
 
type Entry struct { 
   WebSite string 
   WebName string 
   Quality string 
} 
 
var filename string 
```

The definition of the structure is really important because this is how your data is going to be passed to the template file.

The second part of template.go has the following Go code:

```markup
func dynamicContent(w http.ResponseWriter, r *http.Request) { 
   var Data []Entry 
   var f *os.File 
   if filename == "" { 
         f = os.Stdin 
   } else { 
         fileHandler, err := os.Open(filename) 
         if err != nil { 
               fmt.Printf("error opening %s: %s", filename, err) 
               os.Exit(1) 
         } 
         f = fileHandler 
   } 
   defer f.Close() 
   scanner := bufio.NewScanner(f) 
   myT := template.Must(template.ParseGlob("template.gohtml")) 
```

The template.ParseGlob() function is used for reading the external template file, which can have any file extension you want. Using the .gohtml extension might make your life simpler when you are looking for Go template files in your projects.

Although I personally prefer the .gohtml extension for Go template files, .tpl is a pretty common extension that is widely used. You can choose whichever you like.

The third chunk of code from template.go is the following:

```markup
       for scanner.Scan() { 
 
         parts := strings.Fields(scanner.Text()) 
         if len(parts) == 3 { 
               temp := Entry{WebSite: parts[0], WebName: parts[1], Quality: parts[2]} 
               Data = append(Data, temp) 
         } 
   } 
 
   fmt.Println("Serving", r.Host, "for", r.URL.Path) 
   myT.ExecuteTemplate(w, "template.gohtml", Data) 
} 
```

The third parameter to the ExecuteTemplate() function is the data you want to process. In this case, you pass a slice of records to it.

The fourth part of the program is the following:

```markup
func staticPage(w http.ResponseWriter, r *http.Request) { 
   fmt.Println("Serving", r.Host, "for", r.URL.Path) 
   myT := template.Must(template.ParseGlob("static.gohtml")) 
   myT.ExecuteTemplate(w, "static.gohtml", nil) 
} 
```

This function displays a static HTML page, which we are just going to pass through the template engine with the nil data, which is signified by the third argument of the ExecuteTemplate() function. If you have the same function handling different pieces of data, you may end up with cases where there is nothing to render, but keep it there for common code structure.

The fifth part of template.go contains the following Go code:

```markup
func main() { 
   arguments := os.Args 
 
   if len(arguments) == 1 { 
         filename = "" 
   } else { 
         filename = arguments[1] 
   } 
```

The last chunk of Go code from template.go is where you define the supported paths and start the web server using port number 8001:

```markup
   http.HandleFunc("/static", staticPage) 
   http.HandleFunc("/dynamic", dynamicContent) 
   http.ListenAndServe(":8001", nil) 
} 
```

The contents of the template.gohtml file are as follows:

```markup
<!doctype html> 
<htmllang="en"> 
<head> 
   <meta charset="UTF-8"> 
   <title>Using Go HTML Templates</title> 
   <style> 
         html { 
               font-size: 16px; 
         } 
         table, th, td { 
         border: 3px solid gray; 
         } 
   </style> 
</head> 
<body> 
 
<h2 alight="center">Presenting Dynamic content!</h2> 
 
<table> 
   <thead> 
         <tr> 
               <th>Web Site</th> 
               <th>Quality</th> 
         </tr> 
   </thead> 
   <tbody> 
{{ range . }} 
<tr> 
   <td><a href="{{ .WebSite }}">{{ .WebName }}</a></td> 
   <td> {{ .Quality }} </td> 
</tr> 
{{ end }} 
   </tbody> 
</table> 
 
</body> 
</html> 
```

The dot (.) character represents the current data being processed: to put it simply, the dot (.) character is a variable. The {{ range . }} statement is equivalent to a for loop that visits all the elements of the input slice, which are structures in this case. You can access the fields of each structure as .WebSite, .WebName, and .Quality.

The contents of the static.gohtml file are the following:

```markup
<!doctype html> 
<htmllang="en"> 
<head> 
   <meta charset="UTF-8"> 
   <title>A Static HTML Template</title> 
</head> 
<body> 
 
<H1>Hello there!</H1> 
 
</body> 
</html> 
```

If you execute template.go, you will get the following output on the screen:

```markup
$ go run template.go /tmp/sites.html
Serving localhost:8001 for /dynamic
Serving localhost:8001 for /static
```

The following screenshot shows the two outputs of template.go as displayed on a web browser. The sites.html file has three columns, which are the URL, the name and the quality and can have multiple lines. The good thing here is that if you change the contents of the /tmp/sites.html file and reload the web page, you will see the updated contents!

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/31188a00-51f7-4535-9f6b-2ce30fca63c9.png)

Using template.go

Just Imagine

# About JSON

**JSON** stands for JavaScript Object Notation. This is a text-based format designed as an easy and light way to pass information between JavaScript systems.

A simple JSON document has the following format:

```markup
{ "name":"Mihalis",
"surname":"Tsoukalos",
"country":"Greece" }
```

The preceding JSON document has three fields named name, surname, and country. Each field has a single value.

However, JSON documents can have more complex structures with multiple depth levels.

Before seeing some code, I think that it would be very useful to talk about the encoding/json Go package first. The encoding/json package offers the Encode() and Decode() functions that allow the conversion of a Go object into a JSON document and vice versa. Additionally, the encoding/json package offers the Marshal() and Unmarshal() functions that work similarly to Encode() and Decode() and are based on the Encode() and Decode() methods.

The main difference between Marshal()-Unmarshal() and Encode()-Decode() is that the former functions work on single objects, whereas the latter functions can work on multiple objects as well as streams of bytes.

Last, the encoding/json Go package includes two interfaces named Marshaler and Unmarshaler: each one of them requires the implementation of a single method, named MarshalJSON() and UnmarshalJSON(), respectively. These two interfaces allow you to perform custom JSON **Marshalling** and **Unmarshalling** in Go. Unfortunately, those two interfaces will not be covered in this book.

# Saving JSON data

This subsection will teach you how to convert regular data into JSON format in order to send it over a network connection. The Go code of this subsection will be saved as writeJSON.go and will be presented in four parts.

The first chunk of Go code is the expected preamble of the program as well as the definition of two new struct types named Record and Telephone, respectively:

```markup
package main 
 
import ( 
   "encoding/json" 
   "fmt" 
   "os" 
) 
 
type Record struct { 
   Name    string 
   Surname string 
   Tel     []Telephone 
} 
 
type Telephone struct { 
   Mobile bool 
   Number string 
} 
```

Note that only the members of a structure that begin with an uppercase letter will be in the JSON output because members that begin with a lowercase letter are considered private: in this case, all members of Record and Telephone structures are public and will get exported.

The second part is the definition of a function named saveToJSON():

```markup
funcsaveToJSON(filename string, key interface{}) { 
   out, err := os.Create(filename) 
   if err != nil { 
         fmt.Println(err) 
         return 
   } 
 
   encodeJSON := json.NewEncoder(out) 
   err = encodeJSON.Encode(key) 
   if err != nil { 
         fmt.Println(err) 
         return 
   } 
 
   out.Close() 
} 
```

The saveToJSON() function does all the work for us as it creates a JSON encoder variable named encodeJSON, which is associated with a filename, which is where the data is going to be saved. Then, the call to Encode() saves the data of the record to the associated filename and we are done! As you will see in the next section, a similar process will help you read a JSON file and convert it into a Go variable.

The third part of the program has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide a filename!") 
         os.Exit(100) 
   } 
 
   filename := arguments[1] 
```

There is nothing special here: you just get the first command-line argument of the program.

The last part of the utility is the following:

```markup
   myRecord := Record{ 
         Name:    "Mihalis", 
         Surname: "Tsoukalos", 
         Tel: []Telephone{Telephone{Mobile: true, Number: "1234-567"}, 
               Telephone{Mobile: true, Number: "1234-abcd"}, 
               Telephone{Mobile: false, Number: "abcc-567"}, 
         }} 
 
   saveToJSON(filename, myRecord) 
} 
```

Here, we do two things. The first is defining a new Record variable and filling it with data. The second is the call to saveToJSON() for saving the myRecord variable in the JSON format to the selected file.

Executing writeJSON.go will generate the following output:

```markup
$ go run writeJSON.go /tmp/SavedFile
```

After that, the contents of /tmp/SavedFile will be the following:

```markup
$ cat /tmp/SavedFile
{"Name":"Mihalis","Surname":"Tsoukalos","Tel":[{"Mobile":true,"Number":"1234-567"},{"Mobile":true,"Number":"1234-abcd"},{"Mobile":false,"Number":"abcc-567"}]}
```

Sending JSON data over a network requires the use of the net Go standard package that will be discussed in the next chapter.

# Parsing JSON data

This subsection will illustrate how to read a JSON record and convert it into one Go variable that you can use in your own programs. The name of the presented program will be readJSON.go and will be shown to you in four parts.

The first part of the utility is identical to the first part of the writeJSON.go utility:

```markup
package main 
 
import ( 
   "encoding/json" 
   "fmt" 
   "os" 
) 
 
type Record struct { 
   Name    string 
   Surname string 
   Tel     []Telephone 
} 
 
type Telephone struct { 
   Mobile bool 
   Number string 
} 
```

The second part of the Go code is the following:

```markup
funcloadFromJSON(filename string, key interface{}) error { 
   in, err := os.Open(filename) 
   if err != nil { 
         return err 
   } 
 
   decodeJSON := json.NewDecoder(in) 
   err = decodeJSON.Decode(key) 
   if err != nil { 
         return err 
   } 
   in.Close() 
   return nil 
} 
```

Here, you define a new function named loadFromJSON() that is used for decoding a JSON file according to a data structure that is given as the second argument to it. You first call the json.NewDecoder() function to create a new JSON decode variable that is associated with a file, and then you call the Decode() function for actually decoding the contents of the file.

The third part of readJSON.go has the following Go code:

```markup
func main() { 
   arguments := os.Args 
   iflen(arguments) == 1 { 
         fmt.Println("Please provide a filename!") 
         os.Exit(100) 
   } 
 
   filename := arguments[1] 
```

The last part of the program is the following:

```markup
   var myRecord Record 
   err := loadFromJSON(filename, &myRecord) 
   if err == nil { 
         fmt.Println(myRecord) 
   } else { 
         fmt.Println(err) 
   } 
} 
```

If you run readJSON.go, you will get the following output:

```markup
$ go run readJSON.go /tmp/SavedFile
{Mihalis Tsoukalos [{true 1234-567} {true 1234-abcd} {false abcc-567}]}
```

Reading your JSON data from a network will be discussed in the next chapter, as JSON records do not differ from any other kind of data transferred over a network.

# Using Marshal() and Unmarshal()

In this subsection, you will see how to use Marshal() and Unmarshal() in order to implement the functionality of readJSON.go and writeJSON.go. The Go code that illustrates the Marshal() and Unmarshal() functions can be found in marUnmar.go, and this will be presented in four parts.

The first part of marUnmar.go is the expected preamble:

```markup
package main 
 
import ( 
   "encoding/json" 
   "fmt" 
   "os" 
) 
 
type Record struct { 
   Name    string 
   Surname string 
   Tel     []Telephone 
} 
 
type Telephone struct { 
   Mobile bool 
   Number string 
} 
```

The second part of the program contains the following Go code:

```markup
func main() { 
   myRecord := Record{ 
         Name:    "Mihalis", 
         Surname: "Tsoukalos", 
         Tel: []Telephone{Telephone{Mobile: true, Number: "1234-567"}, 
               Telephone{Mobile: true, Number: "1234-abcd"}, 
               Telephone{Mobile: false, Number: "abcc-567"}, 
         }} 
```

This is the same record that is used in the writeJSON.go program. Therefore, so far there is nothing special.

The third part of marUnmar.go is where the marshalling happens:

```markup
   rec, err := json.Marshal(&myRecord) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
   fmt.Println(string(rec)) 
```

Note that json.Marshal() requires a pointer for passing data to it even if the value is a map, array, or slice.

The last part of the program contains the following Go code that performs the unmarshalling operation:

```markup
   var unRec Record 
   err1 := json.Unmarshal(rec, &unRec) 
   if err1 != nil { 
         fmt.Println(err1) 
         os.Exit(100) 
   } 
   fmt.Println(unRec) 
} 
```

As you can see from the code, json.Unmarshal() requires the use of a pointer for saving the data even if the value is a map, array, or slice.

Executing marUnmar.go will create the following output:

```markup
$ go run marUnmar.go
{"Name":"Mihalis","Surname":"Tsoukalos","Tel":[{"Mobile":true,"Number":"1234-567"},{"Mobile":true,"Number":"1234-abcd"},{"Mobile":false,"Number":"abcc-567"}]}
{Mihalis Tsoukalos [{true 1234-567} {true 1234-abcd} {false abcc-567}]}
```

As you can see, the Marshal() and Unmarshal() functions cannot help you store your data into a file: you will need to implement that on your own.

Just Imagine

# Using MongoDB

A relational database is a structured collection of data that is strictly organized into tables. The dominant language for querying databases is SQL. NoSQL databases, such as **MongoDB**, do not use SQL, but various other query languages and do not have a strict structure in their tables, which are called **collections** in the NoSQL terminology.

You can categorize NoSQL databases according to their data model as Document, Key-Value, Graph, and Column-family. MongoDB is the most popular document-oriented NoSQL database that is appropriate for use in web applications.

Document databases were not made for dealing with Microsoft Word documents, but for storing semistructured data.

# Basic MongoDB administration

If you want to use MongoDB on your Go applications, it would be very practical to know how to perform some basic administrative tasks on a MongoDB database.

Most of the tasks presented in this section will be performed from the Mongo shell, which starts by executing the mongo command. If no MongoDB instance is running on your Unix machine, you will get the following output:

```markup
$ mongo
MongoDB shell version v3.4.5
connecting to: mongodb://127.0.0.1:27017
2017-07-06T19:37:38.291+0300 W NETWORK  [thread1] Failed to connect to 127.0.0.1:27017, in(checking socket for error after poll), reason: Connection refused
2017-07-06T19:37:38.291+0300 E QUERY    [thread1] Error: couldn't connect to server 127.0.0.1:27017, connection attempt failed :
connect@src/mongo/shell/mongo.js:237:13
@(connect):1:6
exception: connect failed
```

The previous output tells us two things:

-   The default TCP port number for the MongoDB server process is 27017
-   The mongo executable tries to connect to the 127.0.0.1 IP address, which is the IP address of the local machine

In order to execute the following commands, you should start a MongoDB server instance on your local machine. Once the MongoDB server process is up and running, executing mongo will create the following output:

```markup
$ mongo
MongoDB shell version: 2.4.10
connecting to: test
>
```

The following commands will show you how to create a new MongoDB database and a new MongoDB collection, and how to insert some documents in to that collection:

```markup
>use go;
switched to db go
>db.someData.insert({x:0, y:1})
>db.someData.insert({x:1, y:2})
>db.someData.insert({x:2, y:3})
>db.someData.count()
3
```

Once you try to insert a document into a collection using db.someData.insert(), the collection (someData) will be automatically created if it does not already exist. The last command counts the number of records stored into the someData collection of the current database.

MongoDB will not inform you about any typographical errors you might have. Putting it simply, if you mistype the name of a database or a collection, MongoDB will create a totally new database or a new collection while you are trying to find out what went wrong! Additionally, if you put more, less, or different fields on a document and try to save it, MongoDB will not complain!

You can find the records of a collection using the find() function:

```markup
>db.someData.find()
{ "_id" : ObjectId("595e84cd63883cb3fe7f42f3"), "x" : 0, "y" : 1 }
{ "_id" : ObjectId("595e84d263883cb3fe7f42f4"), "x" : 1, "y" : 2 }
{ "_id" : ObjectId("595e84d663883cb3fe7f42f5"), "x" : 2, "y" : 3 }
```

You can find the list of databases on a running MongoDB instance as follows:

```markup
>show databases;
LXF   0.203125GB
go    0.0625GB
local 0.078125GB
```

Similarly, you can find the names of the collections stored in the current MongoDB database as follows:

```markup
>db.getCollectionNames()
[ "someData", "system.indexes" ]
```

You can delete all the records of a MongoDB collection as follows:

```markup
>db.someData.remove()
>show collections
someData
system.indexes
```

Last, you can delete an entire collection, including its records, as follows:

```markup
>db.someData.drop()
true
>show collections
system.indexes
```

The preceding information will get you going for now, but if you want to learn more about MongoDB, you should visit the documentation site of MongoDB at [https://docs.mongodb.com/](https://docs.mongodb.com/).

# Using the MongoDB Go driver

In order to use MongoDB in your Go programs, you should first have the MongoDB Go driver installed on your Unix machine. The name of the MongoDB Go driver is mgo and you can learn more information about the MongoDB Go driver by visiting [https://github.com/go-mgo/mgo](https://github.com/go-mgo/mgo), [https://labix.org/mgo](https://labix.org/mgo), and [https://docs.mongodb.com/ecosystem/drivers/go/](https://docs.mongodb.com/ecosystem/drivers/go/).

As the driver is not part of the standard Go library, you should first download the required packages using the following two commands:

```markup
$ go get labix.org/v2/mgo
$ go get labix.org/v2/mgo/bson
```

After that, you will be free to use it in your own Go utilities. If you try to execute the program without having the two packages on your Unix system, you will get an error message similar to the following:

```markup
$ go run testMongo.go
testMongo.go:5:2: cannot find package "labix.org/v2/mgo" in any of:
      /usr/local/Cellar/go/1.8.3/libexec/src/labix.org/v2/mgo (from $GOROOT)
      /Users/mtsouk/go/src/labix.org/v2/mgo (from $GOPATH)
testMongo.go:6:2: cannot find package "labix.org/v2/mgo/bson" in any of:
      /usr/local/Cellar/go/1.8.3/libexec/src/labix.org/v2/mgo/bson (from $GOROOT)
      /Users/mtsouk/go/src/labix.org/v2/mgo/bson (from $GOPATH)
```

Note that you might need to install Bazaar on your Unix system in order to execute the two go get commands. You can get more information about the Bazaar version control system at [https://bazaar.canonical.com/](https://bazaar.canonical.com/).

So, you should first try to run a simple Go program that connects to a MongoDB database, creates a new database and a new collection, and adds new documents to it in order to make sure that everything works as expected: the name of the program will be testMongo.go and will be presented in four parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "labix.org/v2/mgo" 
   "labix.org/v2/mgo/bson" 
   "os" 
   "time" 
) 
 
type Record struct { 
   Xvalueint 
   Yvalueint 
} 
```

Here, you see the use of the Go MongoDB driver in the import block. Additionally, you see the definition of a new Go structure named Record that will hold the data of each MongoDB document.

The second part of testMongo.go has the following Go code:

```markup
func main() { 
   mongoDBDialInfo := &mgo.DialInfo{ 
         Addrs:   []string{"127.0.0.1:27017"}, 
         Timeout: 20 * time.Second, 
   } 
 
   session, err := mgo.DialWithInfo(mongoDBDialInfo) 
   if err != nil { 
         fmt.Printf("DialWithInfo: %s\n", err) 
         os.Exit(100) 
   } 
   session.SetMode(mgo.Monotonic, true) 
 
   collection := session.DB("goDriver").C("someData") 
```

Now the collection variable will be used for dealing with the someData collection of the goDriver database: a better name for the database would have been myDB. Note that there was not a goDriver database in the MongoDB instance before running the Go program; this also means that neither the someData collection was there.

The third part of the program is the following:

```markup
   err = collection.Insert(&Record{1, 0}) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   err = collection.Insert(&Record{-1, 0}) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

Here, you insert two documents to the MongoDB database using the Insert() function.

The last portion of testMongo.go contains the following Go code:

```markup
   var recs []Record 
   err = collection.Find(bson.M{"yvalue": 0}).All(&recs) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
 
   for x, y := range recs { 
         fmt.Println(x, y) 
   } 
   fmt.Println("Found:", len(recs), "results!") 
} 
```

As you do not know the number of documents that you will get from the Find() query, you should use a slice of records for storing them.

Additionally, note that you should put the yvalue field in lowercase in the Find() function because MongoDB will automatically convert the fields of the Record structure in lowercase when you are storing them!

Now, execute testMongo.go, as shown here:

```markup
$ go run testMongo.go
0 {1 0}
1 {-1 0}
Found: 2 results!
```

Note that if you execute testMongo.go multiple times, you will find the same documents inserted multiple times into the someData collection. However, MongoDB will not have any problems differentiating between all these documents because the key of each document is the \_id field, which is automatically inserted by MongoDB each time you insert a new document to a collection.

After that, connect to your MongoDB instance using the MongoDB shell command to make sure that everything worked as expected:

```markup
$ mongo
MongoDB shell version v3.4.5
connecting to: mongodb://127.0.0.1:27017
MongoDB server version: 3.4.5
>use goDriver
switched to db goDriver
>show collections
someData
>db.someData.find()
{ "_id" : ObjectId("595f88593fb7048f4846e555"), "xvalue" : 1, "yvalue" : 0 }
{ "_id" : ObjectId("595f88593fb7048f4846e557"), "xvalue" : -1, "yvalue" : 0 }
>
```

Here, it is important to understand that MongoDB documents are presented in JSON format, which you already know how to handle in Go.

Also, note that the Go MongoDB driver has many more capabilities than the ones presented here. Unfortunately, talking more about it is beyond the scope of this book, but you can learn more by visiting [https://github.com/go-mgo/mgo](https://github.com/go-mgo/mgo), [https://labix.org/mgo](https://labix.org/mgo), and [https://docs.mongodb.com/ecosystem/drivers/go/](https://docs.mongodb.com/ecosystem/drivers/go/).

# Creating a Go application that displays MongoDB data

The name of the utility will be showMongo.go and it will be presented in three parts. The utility will connect to a MongoDB instance, read a collection, and display the documents of the collection as a web page. Note that showMongo.go is based on the Go code of template.go.

The first part of the web application is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "html/template" 
   "labix.org/v2/mgo" 
   "net/http" 
   "os" 
   "time" 
) 
 
var DatabaseName string 
var collectionName string 
 
type Document struct { 
   P1 int 
   P2 int 
   P3 int 
   P4 int 
   P5 int 
}
```

You should know in advance the structure of the MongoDB documents that you will retrieve because the field names are hard coded in the struct type and need to match.

The second part of the program is the following:

```markup
func content(w http.ResponseWriter, r *http.Request) { 
   var Data []Document 
   myT := template.Must(template.ParseGlob("mongoDB.gohtml")) 
 
   mongoDBDialInfo := &mgo.DialInfo{ 
         Addrs:   []string{"127.0.0.1:27017"}, 
         Timeout: 20 * time.Second, 
   } 
 
   session, err := mgo.DialWithInfo(mongoDBDialInfo) 
   if err != nil { 
         fmt.Printf("DialWithInfo: %s\n", err) 
         return 
   } 
   session.SetMode(mgo.Monotonic, true) 
   c := session.DB(DatabaseName).C(collectionName) 
 
   err = c.Find(nil).All(&Data) 
   if err != nil { 
         fmt.Println(err) 
         return 
   } 
 
   fmt.Println("Found:", len(Data), "results!") 
   myT.ExecuteTemplate(w, "mongoDB.gohtml", Data) 
} 
```

As before, you connect to MongoDB using mgo.DialWithInfo() with the parameters that were defined in the mgo.DialInfo structure.

The last part of the web application is the following:

```markup
func main() { 
   arguments := os.Args 
 
   iflen(arguments) <= 2 { 
         fmt.Println("Please provide a Database and a Collection!") 
         os.Exit(100) 
   } else { 
         DatabaseName = arguments[1] 
         collectionName = arguments[2] 
   } 
 
   http.HandleFunc("/", content) 
   http.ListenAndServe(":8001", nil) 
} 
```

The contents of MongoDB.gohtml are similar to the contents of template.gohtml and will not be presented here. You can refer to _The html/template package_ section for the contents of template.gohtml.

The execution of showMongo.go will not display the actual data on the screen: you will need to use a web browser for that:

```markup
$ go run showMongo.go goDriver Numbers
Found: 0 results!
Found: 10 results!
Found: 14 results!
```

The good thing is that if the data of the collections is changed, you will not need to recompile your Go code in order to see the changes: you will just need to reload the web page.

The following screenshot shows the output of showMongo.go as displayed on a web browser:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/a11f9a8c-3fc3-414b-83a7-ecffd7aca8de.png)

Using showMongo.go

Note that the Numbers collection contains the following documents:

```markup
>db.Numbers.findOne()
{
      "_id" : ObjectId("596530aeaab5252f5c1ab100"),
      "p1" : -10,
      "p2" : -20,
      "p3" : 100,
      "p4" : -1000,
      "p5" : 10000
}
```

Have in mind that extra data in the MongoDB structure that does not have corresponding fields in the Go structure is ignored.

# Creating an application that displays MySQL data

In this subsection, we will present a Go utility that executes a query on a MySQL table. The name of the new command-line utility will be showMySQL.go and will be presented in five parts.

Note that showMySQL.go will use the database/sql package that provides a generic SQL interface to relational databases for querying the MySQL database.

The presented utility requires two parameters: a username with administrative privileges and its password.

The first part of showMySQL.go is the following:

```markup
package main 
 
import ( 
   "database/sql"  
   "fmt" 
   _ "github.com/go-sql-driver/mysql" 
   "os" 
   "text/template" 
)
```

There is a small change here, as showMySQL.go uses text/template instead of html/template. Note that the drivers that conform to the database/sql interface are never really referenced directly in your code, but they still need to be initialized and imported. The \_ character in front of "github.com/go-sql-driver/mysql" does this by telling Go to ignore the fact that the "github.com/go-sql-driver/mysql" package is not actually used in the code.

You will also need to download the MySQL Go driver:

```markup
$ go get github.com/go-sql-driver/mysql
```

The second part of the utility has the following Go code:

```markup
func main() { 
   var username string 
   var password string 
 
   arguments := os.Args 
   if len(arguments) == 3 { 
         username = arguments[1] 
         password = arguments[2] 
   } else { 
         fmt.Println("programName Username Password!") 
         os.Exit(100) 
   } 
```

The third chunk of Go code from showMySQL.go is the following:

```markup
   connectString := username + ":" + password + "@unix(/tmp/mysql.sock)/information_schema" 
   db, err := sql.Open("mysql", connectString) 
 
   rows, err := db.Query("SELECT DISTINCT(TABLE_SCHEMA) FROM TABLES;") 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(100) 
   } 
```

Here, you manually construct the connection string to MySQL. For reasons of security, a default MySQL installation works with a socket (/tmp/mysql.sock) instead of a network connection. The name of the database that will be used is the last part of the connection string (information\_schema).

You will most likely have to adjust these parameters for your own database.

The fourth part of showMySQL.go is the following:

```markup
   var DATABASES []string 
   for rows.Next() { 
         var databaseName string 
         err := rows.Scan(&databaseName) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(100) 
         } 
         DATABASES = append(DATABASES, databaseName) 
   } 
   db.Close()
```

The Next() function iterates over all the records returned from the select query and returns them one by one with the help of the for loop.

The last part of the program is the following:

```markup
   t := template.Must(template.New("t1").Parse(` 
   {{range $k := .}} {{ printf "\tDatabase Name: %s" $k}} 
   {{end}} 
   `)) 
   t.Execute(os.Stdout, DATABASES) 
   fmt.Println() 
} 
```

This time, instead of presenting the data as a web page, you will receive it as plain text. Additionally, as the text template is small, it is defined in line with the help of the t variable.

Is the use of the template necessary here? Of course not! But it is good to learn how to define Go templates without using an external template file.

Therefore, the output of showMySQL.go will be similar to the following:

```markup
$ go run showMySQL.go root 12345
    
    Database Name: information_schema
    Database Name: mysql
    Database Name: performance_schema
    Database Name: sys
```

The preceding output shows information about the available databases for the current MySQL instance, which is a great way to get information about a MySQL database without having to connect using the MySQL client.

Just Imagine

# A handy command-line utility

In this section, we will develop a handy command-line utility that reads a number of web pages, which can be found in a text file or read from standard input, and returns the number of times a given keyword was found in these web pages. In order to be faster, the utility will use goroutines to get the desired data and a monitoring process to gather the data and present it on the screen. The name of the utility will be findKeyword.go and will be presented in five parts.

The first part of the utility is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "net/http" 
   "net/url" 
   "os" 
   "regexp" 
) 
 
type Data struct { 
   URL     string 
   Keyword string 
   Times   int 
   Error   error 
} 
```

The Data struct type will be used for passing information between channels.

The second part of findKeyword.go has the following Go code:

```markup
func monitor(values <-chan Data, count int) { 
   fori := 0; i< count; i++ { 
         x := <-values 
         if x.Error == nil { 
               fmt.Printf("\t%s\t", x.Keyword) 
               fmt.Printf("\t%d\t in\t%s\n", x.Times, x.URL) 
         } else { 
               fmt.Printf("\t%s\n", x.Error) 
         } 
   } 
} 
```

The monitor() function is where all the information is collected and printed on the screen.

The third part is the following:

```markup
func processPage(myUrl, keyword string, out chan<- Data) { 
   var err error 
   times := 0 
 
   URL, err :=url.Parse(myUrl) 
   if err != nil { 
         out<- Data{URL: myUrl, Keyword: keyword, Times: 0, Error: err} 
         return 
   } 
 
   c := &http.Client{} 
   request, err := http.NewRequest("GET", URL.String(), nil) 
   if err != nil { 
         out<- Data{URL: myUrl, Keyword: keyword, Times: 0, Error: err} 
         return 
   } 
 
   httpData, err := c.Do(request) 
   if err != nil { 
         out<- Data{URL: myUrl, Keyword: keyword, Times: 0, Error: err} 
         return 
   } 
 
   bodyHTML := "" 
   var buffer [1024]byte 
   reader := httpData.Body 
   for { 
         n, err := reader.Read(buffer[0:]) 
         if err != nil { 
               break 
         } 
         bodyHTML = bodyHTML + string(buffer[0:n]) 
   } 
 
   regExpr := keyword 
   r := regexp.MustCompile(regExpr) 
   matches := r.FindAllString(bodyHTML, -1) 
   times = times + len(matches) 
 
   newValue := Data{URL: myUrl, Keyword: keyword, Times: times, Error: nil} 
   out<- newValue 
} 
```

Here, you can see the implementation of the processPage() function that is executed in a goroutine. If the Error field of the Data structure is not nil, then there was an error somewhere.

The reason for using the bodyHTML variable to save the entire contents of a URL is for not having a keyword split between two consecutive calls to reader.Read(). After that, a regular expression (r) is used for searching the bodyHTML variable for the desired keyword.

The fourth part contains the following Go code:

```markup
func main() { 
   filename := "" 
   var f *os.File 
   var keyword string 
 
   arguments := os.Args 
   iflen(arguments) == 1 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(-1) 
   } 
 
   iflen(arguments) == 2 { 
         f = os.Stdin 
         keyword = arguments[1] 
   } else { 
         keyword = arguments[1] 
         filename = arguments[2] 
         fileHandler, err := os.Open(filename) 
         if err != nil { 
               fmt.Printf("error opening %s: %s", filename, err) 
               os.Exit(1) 
         } 
         f = fileHandler 
   } 
 
   deferf.Close() 
```

As you can see, findKeyword.go expects its input from a text file or from standard input, which is the common Unix practice: this technique was first illustrated back in [Chapter 8](https://subscription.imaginedevops.io/book/programming/9781787125643/8)_,_ _Processes and Signals_, in the _Reading from standard input_ section.

The last chunk of Go code for findKeyword.go is the following:

```markup
   values := make(chan Data, len(os.Args[1:])) 
 
   scanner := bufio.NewScanner(f) 
   count := 0 
   forscanner.Scan() { 
         count = count + 1 
         gofunc(URL string) { 
               processPage(URL, keyword, values) 
         }(scanner.Text()) 
   } 
 
   monitor(values, count) 
} 
```

There is nothing special here: you just start the desired goroutines and the monitor() function to take care of them.

Executing findKeyword.go will create the following output:

```markup
$ go run findKeyword.go Tsoukalos /tmp/sites.html
  Get http://really.doesnotexist.com: dial tcp: lookup really.doesnotexist.com: no such host
  Tsoukalos         8      in   http://www.highiso.net/
  Tsoukalos         4      in   http://www.mtsoukalos.eu/
  Tsoukalos         3      in   https://www.imaginedevops.io/networking-and-servers/go-systems-programming
  Tsoukalos         0      in   http://cnn.com/
  Tsoukalos         0      in   http://doesnotexist.com
```

The funny thing here is that the doesnotexist.com domain does actually exist!

Just Imagine

# Exercises

1.  Download and install MongoDB on your Unix machine.
2.  Visit the documentation page of the net/http Go standard package at [https://golang.org/pkg/net/http/](https://golang.org/pkg/net/http/).
3.  Visit the documentation page of the html/template Go standard package at [https://golang.org/pkg/html/template/](https://golang.org/pkg/html/template/).
4.  Change the Go code of getURL.go in order to make it able to fetch multiple web pages.
5.  Read the documentation of the encoding/json package that can be found at [https://golang.org/pkg/encoding/json/](https://golang.org/pkg/encoding/json/).
6.  Visit the MongoDB site at [https://www.mongodb.org/](https://www.mongodb.org/).
7.  Learn how to use text/template by developing your own example.
8.  Change the Go code of findKeyword.go in order to be able to search multiple keywords.

Just Imagine

# Summary

In this chapter, we talked about web development in Go including parsing, marshalling and unmarshalling JSON data, interacting with a MongoDB database; reading data from a MySQL database; creating web servers in Go; creating web clients in Go; and using the http.ServeMux type.

In the next chapter, we will talk about network programming in Go, which includes creating TCP and UDP clients and servers using low level commands. We will also teach you how to develop an RCP client and an RCP server in Go. If you love developing TCP/IP applications, then the last chapter of this book is for you!