# Go Basics

With so many programming languages out there, it’s fair to wonder why anyone would have to invent yet another one. What the background is of the people behind Go and what the problems are they are trying to solve with this new language are some of the items we will address in this chapter.

These topics give us some perspective on the challenges large-scale software development presents to software developers today and why modern technologies such as programming languages are constantly evolving.

By the end of this chapter, you should have a better understanding of where Go comes from and its role in developing distributed systems running on multi-core processors, as well as be familiar with Go’s source code structure as we go through the following areas:

-   What is Go?
-   Go’s guiding principles
-   Go source code file structure
-   Go packages and modules
-   Compiling Go programs
-   Running Go programs online
-   Exploring the Go tool to manage Go source code

Just Imagine

# Technical requirements

We assume that you have basic familiarity with the command line, Git, and GitHub. You can find the code examples for this chapter in the book’s GitHub repository, [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/3), in the `ch02` folder.

To run the examples, follow these steps:

1.  Install Go 1.17 or later for your operating system. You can follow the instructions in [_Chapter 1_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_01.xhtml#_idTextAnchor015), _Introduction_, in the _Installing Go_ on your computer section, or go to [https://go.dev/doc/install](https://packages.cloud.google.com/apt/doc/apt-key.gpg).
2.  Clone the book’s GitHub repository with `git clone` at [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go.git](https://kubernetes.io/docs/reference/ports-and-protocols/).
3.  Change the directory to an example’s folder – `cd Network-Automation-with-Go/ch02/pong`.
4.  Execute `go` `run main.go`.

Just Imagine

# What is Go?

During the second half of 2007, _Robert Griesemer_, _Rob Pike_, and _Ken Thompson_ started discussing the design of a new programming language that would solve some problems they were experiencing when writing software at Google, such as the increased complexity to use some languages, long code compilation times, and not being able to program efficiently on multiprocessor computers.

_Rob Pike_ was trying to take some concurrency and communicating channels ideas into C++, based on his earlier work on the Newsqueak language in 1988, as he describes in _Go: Ten years and climbing_ (_Further reading_) and _Less is exponentially more_ (_Further reading_). This turned out to be too hard to implement. He would work out of the same office with _Robert Griesemer_ and _Ken Thompson_. Ken had worked together with Rob Pike in the past to create the character-encoding UTF-8, while _Ken Thompson_ had designed and implemented the Unix operating system and invented the B programming language (the predecessor to the C programming language).

They chose the name **Go** for this new programming language because it’s short, but the DNS entry for [go.com](https://cloud.google.com/architecture/reduce-carbon-footprint) wasn’t available, so Go’s website ended up at [golang.org](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/what-is-ccft.html). And so, **golang** became a nickname for Go. While golang is convenient for search queries, it’s not the name of the language (which is Go):

![Figure 2.1 – The initial Go discussion email thread
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_01.jpg)

Figure 2.1 – The initial Go discussion email thread

Though they initially thought of C/C++ to be the starting point, they ended up starting from scratch to define a more expressive language, despite a large number of simplifications when compared to its predecessors. Go inherits some things from C, such as, but not limited to, basic data types, expression syntax, pointers, and compilation to machine code, but it doesn’t have things such as the following:

-   Header files
-   Exceptions
-   Pointer arithmetic
-   Subtype inheritance (no subclasses)
-   `this` in methods
-   Promotion to a superclass (it uses embedding instead)
-   Circular dependencies

Pascal, Oberon, and Newsqueak are among the programming languages that have influenced Go. In particular, its concurrency model comes from _Tony Hoare_’s **Communicating Sequential Processes** (**CSPs**) (_Further reading_) white paper, and CSP’s implementations in _Rob Pike_’s interpreted language Newsqueak and, later, Phil Winterbottom’s C-like compiled version, Alef. The next figure shows Go’s family tree:

![Figure 2.2 – The Go ancestors
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_02.jpg)

Figure 2.2 – The Go ancestors

The number of C++ programmers that come to Go is just a few compared to what the Go founders expected. Most Go programmers actually come from languages such as Python and Ruby.

Go became an open source project on November 10, 2009. They host Go’s source code at [https://go.googlesource.com/go](https://www.microsoft.com/en-us/sustainability/emissions-impact-dashboard?activetab=pivot_2%3aprimaryr12) and keep a mirror of the code at [https://github.com/golang/go](https://cloud.google.com/carbon-footprint) where you can submit pull requests. While Go is an open source programming language, it’s actually supported by Google.

They wrote the first Go compiler in C, but they later converted it to Go. Russ Cox describes this in detail in Go 1.3+ Compiler Overhaul (_Further reading_). As mind-blowing as it may sound, the Go source code of today is written in Go.

They released Go 1 on March 28, 2012. We highlight some notable changes to the language since then in the summarized version of Go’s timeline in the next figure:

![Figure 2.3 – Go's brief timeline
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_03.jpg)

Figure 2.3 – Go’s brief timeline

Go is a stable language, and the semantics should not change unless Go 2 happens. The only change that the Go team has confirmed at this point is the addition of generic programming using type parameters in early 2022 (Go 1.18), as described in the Type Parameters Proposal (_Further reading_).

Go is a programming language that attempts to combine the ease of programming of a dynamically typed language with the efficiency and safety of a statically typed language. It builds executable files in seconds, and with Go’s first-class support for concurrency, we can take full advantage of multi-core CPUs.

Before we dive into Go code, we cover some guiding principles that make Go unique through the Go proverbs.

Just Imagine

# Go Proverbs

_Rob Pike_ introduced the Go language proverbs at _Gopherfest_ in _2015_ to explain or teach Go philosophically. These are general guidelines that Go developers tend to adhere to. Most of these proverbs are good practices – but optional – that convey the spirit of the language.

We only include our favorite proverbs here. You can check out the full list at _Go Proverbs_ (_Further reading_):

-   Gofmt’s style is no one’s favorite, yet gofmt is everyone’s favorite. When you write code in Go, you don’t have to worry about the debate of white spaces versus tabs, or where you put braces or curly brackets. Gofmt (`gofmt`) formats your code with a prescriptive style guide, so all Go code looks the same. This way, you don’t have to think about it when you write or read Go code:
-   **Clear is better than clever**: Go favors clear code over clever code that is difficult to analyze or describe. Write code other people can read and with behavior they can understand.
-   **Errors are values**: An error in Go is not an exception. It’s a value you can use in your program logic – as a variable, for example.
-   **Don’t just check errors; handle them gracefully**: Go encourages you to think about whether you should do something with an error, instead of just returning it and forgetting about it. Depending on the error, maybe you can trigger a different execution path, add more info to it, or save it for later.
-   **A little copying is better than a little dependency**: If you only need a few lines from a library, maybe you can just copy those lines instead of importing the entire library to keep your dependency tree under control and make your code more compact. This way, your program not only compiles faster but is also more manageable and simpler to understand.
-   **Don’t communicate by sharing memory; share memory by communicating**: This describes how concurrent processes in Go can coordinate between each other. In other languages, concurrent processes communicate by sharing memory, which you have to protect with locks to prevent a data race condition when these processes try to access a memory location concurrently. Go, in contrast, uses channels instead to pass references to data between processes, so only one process has access to the data at a time.
-   **Concurrency is not parallelism**: Concurrency is structuring the execution of independent processes, whose instructions are not necessarily executed in sequence. Whether these instructions run in parallel depends on the availability of different CPU cores or hardware threads. _Rob Pike_’s _Concurrency is not Parallelism_ (_Further reading_) talk is a must for Go developers.

The Go proverbs cover different aspects of Go, from formatting your Go code to how Go achieves concurrency.

Now, it’s time to roll up our sleeves as we start looking into Go source code files.

Just Imagine

# Go source code files

While there isn’t a filename convention for Go source code files, their filenames are typically one-word, all lowercase, and include an underscore if it has more than one word. It ends with the `.``go` suffix.

Each file has three parts:

-   **Package clause**: This defines the name of the package a file belongs to.
-   **Import declaration**: This is a list of packages that you need to import.
-   **Top-level declaration**: This is constant, variable, type, function, and method declarations with a package scope. Every declaration here starts with a keyword (`const`, `var`, `type`, or `func`):

```markup
// package clause
package main
// import declaration
import "fmt"
// top level declaration
const s = "Hello, 世界"
func main() {
    fmt.Println(s)
}
```

The code example shows the package declaration for the `main` package at the top. It follows the import declaration, where we specify that we use the `fmt` package in this file. Then, we include all declarations in the code – in this case, an `s` constant and the `main` function.

## Packages

A package is one or more `.go` files in the same folder that declares the related constants, types, variables, and functions. These declarations are accessible to every file in the same package, so breaking down the code into different files is optional. It’s more of a personal preference on how to better organize code.

In the standard library, they divide the code into separate files for larger packages. The `encoding/base64` package has one `.go` file (other than the test and example files), such as the following:

```markup
$ ls -1 /usr/local/go/src/encoding/base64/ | grep -v _test.go
base64.go
```

By contrast, the `encoding/json` package has nine `.go` source code files:

```markup
$ ls -1 /usr/local/go/src/encoding/json/ | grep -v _test.go
decode.go
encode.go
fold.go
fuzz.go
indent.go
scanner.go
stream.go
tables.go
tags.go
```

Package names are short and meaningful (no underscore). Users of a package refer to the package name when importing something from it – for example, the `Decode` method exists in the `json` and `xml` packages. Users can call these methods with `json.Decode` and `xml.Decode`, respectively.

One special package is `main`. This is the entry point for any program that imports other packages. This package must have a `main` function that takes no arguments and returns no value, such as the code example at the beginning of this section.

## Go modules

Go modules became the default way to release packages in Go 1.16. They were first introduced in Go 1.11, back in 2018, to improve dependency management in Go. It allows you to define an import path and the dependencies for a package or collection of packages.

Let’s define a small package called `ping`, with a `Send` function that returns a string with the word `pong`:

```markup
package ping
func Send() string {
    return "pong"
}
```

This is the [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch02/ping/code.go](http://www.sustainableitplaybook.com) file in the book’s GitHub repository. You can create a module for this package with the `go mod init` command at the root folder of this example (`ch02/ping`). The argument for this command should be the module location, where users can get access to it. The result is a `go.mod` file with the import path and a list of external package dependencies in it:

```markup
ch02/ping$ go mod init github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping
go: creating new go.mod: module github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping
```

With this, anyone can now import this package. The following program imports this package to the `pong` output:

```markup
package main
import (
    "fmt"
    "github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping"
)
func main() {
    s := ping.Send()
    fmt.Println(s)
}
```

You can run this program from the Go Playground (_Further reading_), which imports the module we just created. This is also a great segue into the next section on packet importing and a sneak peek into the Go Playground section that we will cover in just a few more pages.

## Importing packages

The `import` keyword lists the packages to import in a source file. The import path is the module path, followed by the folder where the package is within the module, unless the package is in the standard library, in which case you only need to reference the directory. Let’s examine an example of each scenario.

To give an example, the `google.golang.org/grpc` module has a package in the `credentials` folder. You would import it with `google.golang.org/grpc/credentials`. The last part of the path is how you prefix the package types and functions, `credentials.TransportCredentials` and `credentials.NewClientTLSFromFile`, respectively, in the next code sample.

Go’s standard library (_Further reading_) at `go/src` is a collection of packages of the `std` module. The `fmt` folder hosts the package that implements functions to format input and output. The path to import this package is just `fmt`:

```markup
package xrgrpc
import (
    "fmt"
    /* ... <omitted for brevity > ... */
    "google.golang.org/grpc/credentials"
)
func newClientTLS(c client) (credentials.TransportCredentials, error) {
    if c.Cert != "" {
                return credentials.NewClientTLSFromFile(...)
    }
    /* ... <omitted for brevity > ... */
    fmt.Printf("%s", 'test')
    /* ... <omitted for brevity > ... */
}
```

Packages do not live in a central repository such as `maven`, `pip`, or `npm`. You can share your code by upstreaming it to a version control system and distribute it by sharing its location. Users can download it with the `go` command (`go install` or `go get`).

For developing and testing purposes, you can reference local packages by pointing to their local path in the `go.mod` file:

```markup
module github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/pong
 
go 1.17
 
require github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping v0.0.0-20220223180011-2e4e63479343
 
replace github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping v1.0.0 => ../ping
```

In the `ch02/pong` example, the Go tool automatically created the first three lines of the `go.mod` file for us, referencing the ping module from the book’s GitHub repository (_Further reading_). We later added a fourth line to replace that module, with the contents of the local version of it (`../ping`).

## Comments

Code comments in Go play a key role, as they become your package documentation. The `go doc` tool takes the comments preceding a type, constant, function, or method that you export in a package as a document string for that declaration, producing an HTML file that the tool presents as a web page.

To give an example, all public Go packages (_Further reading_) display this autogenerated documentation.

Go offers two ways to create comments:

-   C++-style `//` line comments, which is the most common form:
    
    ```markup
    // IsPrivate reports whether ip is a private address, according to
    ```
    
    ```markup
    // RFC 1918 (IPv4 addresses) and RFC 4193 (IPv6 addresses).
    ```
    
    ```markup
    func (ip IP) IsPrivate() bool {
    ```
    
    ```markup
        if ip4 := ip.To4(); ip4 != nil {
    ```
    
    ```markup
            return ip4[0] == 10 ||
    ```
    
    ```markup
                (ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
    ```
    
    ```markup
                (ip4[0] == 192 && ip4[1] == 168)
    ```
    
    ```markup
        }
    ```
    
    ```markup
        return len(ip) == IPv6len && ip[0]&0xfe == 0xfc
    ```
    
    ```markup
    }
    ```
    
-   C-style `/* */` block comments, which are primarily for package descriptions or large blocks of formatted/indented code:
    
    ```markup
    /*
    ```
    
    ```markup
    Copyright 2014 The Kubernetes Authors.
    ```
    
    ```markup
    Licensed under the Apache License, Version 2.0 (the "License");
    ```
    
    ```markup
    ...
    ```
    
    ```markup
    See the License for the specific language governing permissions and
    ```
    
    ```markup
    limitations under the License.
    ```
    
    ```markup
    */
    ```
    
    ```markup
    package kubectl
    ```
    

Dave Cheney in _Practical Go: Real-world advice for writing maintainable Go programs_ (_Further reading_) suggests that a code comment should explain one – and only one – of these three things:

-   What it does
-   How something does what it does
-   Why something is why it is

A good practice is to make comments on variables that describe their contents, rather than their purpose. You could use the name of the variable to describe its purpose. This brings us to the naming style.

## Names

The convention for declaring names in Go is to use camel case (MixedCaps or mixedCaps) instead of, for example, dashes or underscores when you use more than one word for the name of a function or variable. The exception to the rule are acronyms that have a consistent case, such as `ServeHTTP` and not `ServeHttp`:

```markup
package net
// IsMulticast reports whether ip is a multicast address.
func (ip IP) IsMulticast() bool {
     if ip4 := ip.To4(); ip4 != nil {
         return ip4[0]&0xf0 == 0xe0
     }
     return len(ip) == IPv6len && ip[0] == 0xff
}
```

The first letter of the name determines whether the package exports this top-level declaration. Packages export names that start with a capital letter. These names are the only ones an external user of the package can reference when importing the package – for example, you can reference `IsMulticast` in the preceding code sample from another package as `net.IsMulticast`:

```markup
package net
func allFF(b []byte) bool {
     for _, c := range b {
          if c != 0xff {
                 return false
          }
     }
     return true
}
```

If the first letter is lowercase, no other package has access to this resource. Packages can have declarations that are only for internal consumption. The `allFF` function in the last code example comes from the `net` package. This means only functions in the `net` package can call the `allFF` function.

Languages such as Java and C++ have explicit keywords such as `public` and `private` to control access to types and methods. Python follows the convention of naming variables or methods for internal use with a single underscore prefix. In Go, you can access any variable or method that starts with a lowercase letter from any source code file within the package, but not from another package.

## Executing your Go code

The Go compiler translates Go programs into machine code, producing a binary file. Aside from your program, the binary includes the Go runtime, which offers services such as garbage collection and concurrency. Having access to binary files that work for different platforms makes Go programs very portable.

Let’s compile the [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch02/pong/code.go](https://www.cio.gov/assets/files/Application-Rationalization-Playbook.pdf) file of the book’s GitHub repository with the `go build` command. You can also time this operation with the `time` command to see how fast Go builds really are:

```markup
ch02/pong$ time go build
real  0m0.154s
user  0m0.190s
sys   0m0.070s
```

Now, you can execute the binary file. The default filename is the package name, `pong`. You can change the filename with the `-o` option of the `go build` command. There will be more on this in the _Go_ _tool_ section:

```markup
ch02/pong$ ./pong
pong
```

If you don’t want to generate a binary or executable file and only run the code, you can use the `go run` command instead:

```markup
ch02/pong$ go run main.go 
pong
```

Either option is fine, and it probably comes down to a matter of personal preference or whether you intend to share the compiled artifact with others or deploy it to servers.

Go files have three main parts and they are organized into packages and modules.

You can run all the examples on your computer after installing Go, or you can run them online, as we discuss in the next section.

Just Imagine

# Running Go programs online

Sometimes, you need to test some code quickly or just want to share a code example with someone who might not have Go installed on their computer. In those situations, there are at least three websites where you can run and share Go code for free:

-   The Go Playground
-   The Go Play Space
-   The Gotip Playground

They all share the backend infrastructure, but with subtle differences.

## The Go Playground

The Go team runs the Go Playground ([https://play.golang.org/](https://github.com/Green-Software-Foundation/awesome-green-software)) on golang.org’s servers. They shared some insights and its architecture in the article _Inside the Go Playground_ (_Further reading_), but more recently, _Brad Fitzpatrick_ shared the history and the implementation details of the latest incarnation of the Go Playground (_Further reading_).

This service receives your program, runs it on a sandbox, and returns its output. This is very convenient if you are on your mobile phone, for example, and you want to verify the syntax of a function or something else.

![Figure 2.4 – The Go Playground
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_04.jpg)

Figure 2.4 – The Go Playground

If you are curious about how they built this service or you want to run it locally in your environment, make sure you check out the Playground source code (_Further reading_).

## The Go Play Space

If you can’t live without syntax highlighting, go to the Go Play Space (_Further reading_). This is an experimental alternative Go Playground frontend. They proxy the code execution to the official Go Playground so that programs work the same. They also store shared snippets on the [golang.org](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/what-is-ccft.html) servers:

![Figure 2.5 – The Go Play Space
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_05.jpg)

Figure 2.5 – The Go Play Space

_Figure 2__.5_ shows some extra features that the Go Play Space includes besides syntax highlighting, such as auto-closing braces, access to documentation, and different UI themes.

![Figure 2.6 – Building a house in the Go Play Space
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_06.jpg)

Figure 2.6 – Building a house in the Go Play Space

We could not pass over the fact that it also has a Turtle graphics mode to help you visualize algorithms for fun, such as having a gopher build a house, as shown in _Figure 2__.6_.

## A look into the Future

The Gotip Playground runs on golang.org’s servers as well. This instance of the Go playground runs the latest development branch of Go. You can use it to test upcoming features that are in active development, such as the syntax described in the Type Parameters Proposal (_Further reading_) or the new `net/netip` package, without having to install more than one Go version on your system if you don’t want to.

![Figure 2.7 – The Gotip Playground
](https://static.packt-cdn.com/products/9781800560925/graphics/image/Figure_2.7.jpg)

Figure 2.7 – The Gotip Playground

You can access the Gotip Playground via [https://gotipplay.golang.org/](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/what-is-ccft.html) or by selecting the **Go dev branch** dropdown at [https://go.dev/play/](https://www.microsoft.com/en-us/sustainability/emissions-impact-dashboard?activetab=pivot_2:primaryr12).

These are all great options to run Go programs online that are available to you at no cost. In the next section, we go back to working on the command line as we explore the Go tool to manage Go source code.

Just Imagine

# The Go tool

One of the convenient things about Go – as a programming language – is that a single tool handles all interactions with, and operations on, the source code. When installing Go, make sure that the `go` tool is in the searchable OS path so that you can invoke it from any command-line terminal. The user experience, regardless of the OS or platform architecture, is uniform and doesn’t require any customization when moving from one machine to another.

IDEs also use the `go` tool to build and run code, report errors, and automatically format Go source code. The `go` executable accepts a _verb_ as the first argument that determines what `go` tool function to apply to Go source files:

```markup
$ go 
Go is a tool for managing Go source code.
Usage:
     go <command> [arguments]
The commands are:
     bug         start a bug report
     build       compile packages and dependencies
     ...       
     mod         module maintenance
     run         compile and run Go program
     test        test packages
     tool        run specified go tool
     version     print Go version
     vet         report likely mistakes in packages
```

We’re only exploring a subset of the functions of the Go tool in this section. You can find the full list and every detail of each one in the Go `cmd` documentation (_Further reading_). The commands we’re covering are as follows:

-   `build`
-   `run`
-   `mod`
-   `get`
-   `install`
-   `fmt`
-   `test`
-   `env`

These help you build and run your Go programs, manage their dependencies, and format and test your code.

## Build

We use the `go build` command to compile a Go program and generate an executable binary. If you are not using Go modules yet, the command expects a list of Go source files to compile as an argument. It generates a binary file as a result, with the same name as the first source file (without the `.go` suffix). In the `ch02/hello` folder of the book’s GitHub repository (_Further reading_), we have the `main.go` and `vars.go` files.

You can build an executable file for the program in these files with the `go` `build` command:

```markup
ch02/hello$ go build *.go
ch02/hello$ ./main
Hello World
```

Packaging compiled binaries is a common way of distributing Go programs, since it allows users of a program to skip the compilation stage and reduce the installation procedure to just a few commands (`download` and `unzip`). But you can only run this binary file on a machine with the same architecture and OS. To produce binary files for other systems, you can cross-compile to a wide range of OSs and CPU architectures. For example, the following table shows some target CPU instruction sets that are supported:

![Table 2.1 – Some supported CPU architectures
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_Table_2.1.jpg)

Table 2.1 – Some supported CPU architectures

Out of a long list of supported operating systems, the next table shows the most popular options:

![Table 2.2 – Some supported OSs
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_Table_2.2.jpg)

Table 2.2 – Some supported OSs

The `GOOS` and `GOARCH` environment variables allow you to generate cross-compiled binaries for any other supported system. If you are on a Windows machine, you can generate a binary for macOS running on a 64-bit Intel processor with the following command:

```markup
ch02/hello$ GOOS=darwin GOARCH=amd64 go build *.go
```

The `go tool dist list` command shows a complete set of unique combinations of OSs and architectures that the Go compiler supports:

```markup
$ go tool dist list
...
darwin/amd64
darwin/arm64
...
linux/386
linux/amd64
linux/arm
linux/arm64
...
windows/386
windows/amd64
```

The `go build` command supports different flags to change its default behavior. Two of the most popular flags are `-o` and `-ldflags`.

You can use `-o` to override the default binary name with a name of your preference. In the example, we’ve selected `another_name`:

```markup
ch02/hello$ go build -o another_name *.go
ch02/hello$ ./another_name
Hello World
```

To inject environment data at compile time into your program, use `-ldflags` with a reference to a variable and its value. This way, you can have access to build information during the program execution, such as the date you compiled the program or the version of the source code (`git commit`) you compiled it from:

```markup
ch02/hello$ go build -ldflags='-X main.Version=1.0 -X main.GitCommit=600a82c442' *.go
ch02/hello$ ./main
Version: "1.0"
Git Commit: "600a82c442"
Hello World
```

The last example is a very common way of version-tagging a Go binary. The benefit of this approach is that it doesn’t require any changes to the source code, and you can automate the entire process in a continuous delivery pipeline.

## Run

Another way to run a Go program is by using the `go run` command. It accepts the same flags as `go build` with two differences:

-   It doesn’t produce a binary.
-   It runs the program right after compilation.

The most common use case for `go run` is local debugging and troubleshooting, where a single command combines the processes of compilation and execution:

```markup
ch02/hello$ go run {main,vars}.go
Hello World
```

In the example, we run the program in the `main.go` and `vars.go` files, which produces the `Hello` `World` output.

## Mod

With the introduction of Go modules, the `go` tool got an extra command to work with them – `go mod`. To describe its functionally, let’s review a typical Go program development workflow:

1.  You create a new project in a folder and initialize Go modules with the `go mod init` command, with a reference to the module name – `go mod init example.com/my-project`. This creates a pair of files, `go.mod` and `go.sum`, that keep track of your project’s dependencies.

The next output shows the size of these two files of a real-life project. `go.mod` lists all the dependencies and is relatively small in size compared to `go.sum`, which has the checksum for all the dependencies:

```markup
$ ls -1hs go.*
4.0K go.mod
 92K go.sum
```

If you plan to share this project with others, the name of the module should be a path that is reachable on the internet. It normally points to your source code repository – for example, `github.com/username/my-project`. A real-life example is `github.com/gohugoio/hugo/`.

1.  As you develop your code and add more and more dependencies, the `go` tool updates the `go.mod` and `go.sum` files automatically whenever you run the `go build` or `go` `run` commands.
2.  When you add a dependency, the `go` tool locks its version in the `go.mod` file to prevent accidental code breakages. If you decide you want to update to a newer minor version, you can use the `go get -u` `package@version` command.
3.  If you remove a dependency, you can run `go mod tidy` to clean up the `go.mod` file.
4.  The two `go.*` files contain a full list of dependencies, including the ones that are not directly referenced in your code, that are indirect or chained/transitive dependencies. If you want to find out why a particular dependency is present in your `go.mod` file, you can use the `go mod why package` or `go mod graph` commands to print the dependency tree on the screen:
    
    ```markup
    hugo$ go mod why go.opencensus.io/internal
    ```
    
    ```markup
    # go.opencensus.io/internal
    ```
    
    ```markup
    github.com/gohugoio/hugo/deploy
    ```
    
    ```markup
    gocloud.dev/blob
    ```
    
    ```markup
    gocloud.dev/internal/oc
    ```
    
    ```markup
    go.opencensus.io/trace
    ```
    
    ```markup
    go.opencensus.io/internal
    ```
    

The `go list` command can also be of help. It lists all the module dependencies:

```markup
hugo$ go list -m all | grep ^go.opencensus.io
go.opencensus.io v0.23.0
```

It also lists the actual package dependencies:

```markup
hugo$ go list all | grep ^go.opencensus.io
go.opencensus.io
go.opencensus.io/internal
go.opencensus.io/internal/tagencoding
go.opencensus.io/metric/metricdata
go.opencensus.io/metric/metricproducer
go.opencensus.io/plugin/ocgrpc
...
go.opencensus.io/trace/propagation
go.opencensus.io/trace/tracestate
```

If you prefer a visual representation, there are projects such as Spaghetti (_Further reading_), a dependency analysis tool for Go packages, that can present this information with a user-friendly interface, as shown in _Figure 2__.8_:

![Figure 2.8 – Hugo dependency analysis
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_08.jpg)

Figure 2.8 – Hugo dependency analysis

One thing that is important to mention is that Go modules use semantic versioning. If you need to import a package that is part of a module at major version 2 or higher, you need to include that major version suffix in their import path (`github.com/username/my-project/v2 v2.0.0`, for example).

Before we move to the next command, let’s create a `go.mod` file for the example in the `ch02/hello` folder of the book’s GitHub repository (_Further reading_):

```markup
ch02/hello$ go mod init hello
go: creating new go.mod: module hello
go: to add module requirements and sums:
go mod tidy
ch02/hello$ go mod tidy
ch02/hello$ go build
ch02/hello$ ./hello
Hello World
```

Now, you can build a binary file for the program with `go build` without having to reference all the Go files in the folder (`*.go`).

## Get

Before the Go 1.11 release, you could use the `go get` command to download and install Go programs. This legacy behavior is being completely deprecated, starting from Go 1.17, so we won’t cover it here. From now on, the sole role of this command is the management of dependencies in the `go.mod` file to update them to a newer minor version.

## Install

The easiest way to compile and install a Go binary without explicitly downloading the source code is to use the `go install [packages]` command. In the background, the `go` tool still downloads the code if necessary, runs `go build`, and copies the binary into the `GOBIN` directory, but the `go` tool hides all this from the end user:

```markup
$ go install example.com/cmd@v1.2.3
$ go install example.com/cmd@latest
```

The `go install` command accepts an optional version suffix – for example, `@latest` – and falls back to the local `go.mod` file if the version is missing. Thus, when running `go install`, it’s recommended to always specify a version tag to avoid errors if the `go` tool cannot find a local `go.mod` file.

## Fmt

Go takes much of the code formatting out of developers’ hands by shipping an opinionated formatting tool that you can invoke with the `go fmt` command pointing to your Go source code – for example, `go fmt source.go`.

[_Chapter 1_](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/2B16971_01.xhtml#_idTextAnchor015), _Introduction_, covers how this improves code readability by making all Go code look similar. Most IDEs with plugins for Go automatically format your code every time you save it, making it one less problem to worry about for developers.

## Test

Go is also opinionated when it comes to testing. It makes a few decisions on behalf of developers about the best way to organize code testing to unify the user experience and discourage the use of third-party frameworks:

1.  It automatically executes all files with the `_test.go` suffix in their filenames when you run the `go test` command. This command accepts an optional argument that specifies which package, path, or source file to test.
2.  The Go standard library includes a special `testing` package that works with the `go test` command. Aside from unit test support, this package offers comprehensive coverage reports and benchmarks.

To put this into practice, we include a test program for the `ping` package that we described in the Go modules section. The `ping` package has a `Send` function, which returns the `pong` string when called. The test we perform should verify this. In the test program, we start by defining a string with the value we expect (`pong`) and then compare it to the result of the `ping` function. The `code_test.go` file ([https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go/blob/main/ch02/ping/code\_test.go](https://sdg-tracker.org/)) in the same folder as `ping` represents this in Go code:

```markup
package ping_test
import (
    "github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping" 
    "testing"
)
func TestSend(t *testing.T) {
    want := "pong"
    result := ping.Send()
    if result != want {
        t.Fatalf("[%s] is incorrect, we want [%s]", result, want)
    }
}
```

All test functions have a `TestXxx`(`t *testing.T`) signature, and whether they have access to any other functions and variables defined in the same package depends on how you name the package:

-   **ping**: This gives you access to everything in the package.
-   **ping\_test**: This is a package type (the `_test` suffix) that can live in the same folder as the package you are testing, but it does not have access to the original package variables and methods, so you must import it as any other user would do it. It’s an effective way to document how to use the package while testing it. In the example, we use the `ping.Send` function instead of `Send` directly, as we are importing the package.

This is an assurance that the `Send` function always does the same even if they must optimize the code later. Now, every time you change the code, you can run the `go test` command to verify that the code still behaves the way you expect. By default, when you run `go test`, it prints the results of every test function it finds along with the time to execute them:

```markup
ch02/ping$ go test
PASS
ok github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping 0.001s
```

If someone makes a change in the code that modifies the behavior of the program so that it can no longer pass the test cases, we are in the presence of a potential bug. You can proactively identify software issues with the `go test` command. Let’s say they change the return value of the `Send` function to `p1ong`:

```markup
func Send() string {
    return "p1ong"
}
```

The `go test` command then generates an error the next time your continuous integration pipeline runs the test cases:

```markup
ch02/ping$ go test
--- FAIL: TestSend (0.00s)
  code_test.go:12: [p1ong] is incorrect, we want [pong]
FAIL
exit status 1
FAIL github.com/ImagineDevOps DevOps/Network-Automation-with-Go/ch02/ping 0.001s
```

Now, you know you can’t promote this code to production. The benefit of testing is that you reduce the number of software bugs your users might run into, as you can catch them beforehand.

## Env

The `go env` command displays the environment variables that the `go` command uses for configuration. The `go` tool can print these variables as flat text or in the JSON format with the `-``json` flag:

```markup
$ go env -json
{
    ...
    "GOPROXY": "https://proxy.golang.org,direct",
    "GOROOT": "/usr/local/go",
    ...
    "GOVERSION": "go1.17",
    "PKG_CONFIG": "pkg-config"
}
```

You can change the value of a variable with `go env -w <NAME>=<VALUE>`. The next table describes some of these configuration environment variables:

![Table 2.3 – Some configuration environment variables
](https://static.packt-cdn.com/products/9781800560925/graphics/image/B16971_02_Table_2.3.jpg)

Table 2.3 – Some configuration environment variables

When you change a variable, the `go` tool stores its new value in the path specified by the `GOENV` variable, which defaults to `~/.config/go`:

```markup
$ go env -w GOBIN=$(go env GOPATH)/bin
$ cat ~/.config/go/env
GOBIN=/home/username/go/bin
```

The preceding output example shows how to set the `GOBIN` directory explicitly and how to verify it.

Go offers a command-line utility that helps you manage your source code, from formatting your code to performing dependency management.

Just Imagine

# Summary

In this chapter, we reviewed Go’s origins and its guiding principles, and how you should structure Go source code files and work with dependencies to run your Go programs.

In the next chapter, we will drill down into the semantics of the Go language, the variable types, math logic, control flow, functions, and, of course, concurrency.

Just Imagine

# Further reading

-   _Less is exponentially_ _more_: [https://commandcenter.blogspot.com/2012/06/less-is-exponentially-more.html?m=1](https://docs.aws.amazon.com/wellarchitected/latest/sustainability-pillar/sustainability-pillar.html)
-   _Go: Ten years and_ _climbing_: [https://commandcenter.blogspot.com/2017/09/go-ten-years-and-climbing.html](https://github.com/kubernetes/kubernetes/releases)
-   _Communicating Sequential_ _Processes_: [https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf](https://kubernetes.io/docs/concepts/cluster-administration/addons/)
-   _Go 1.3+ Compiler_ _Overhaul_: [https://golang.org/s/go13compiler](https://docs.projectcalico.org/manifests/tigera-operator.yaml)
-   _Type Parameters_ _Proposal_: [https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/3)
-   _Go_ _Proverbs_: [https://go-proverbs.github.io/](https://multipass.run/)
-   _Concurrency is not_ _Parallelism_: [https://www.youtube.com/watch?v=oV9rvDllKEg](https://multipass.run/)
-   The book’s GitHub repository: [https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/3)
-   Go Playground: [https://go.dev/play/p/ndfJcayqaGV](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/3)
-   Playground source code: [https://go.googlesource.com/playground](https://docs.aws.amazon.com/wellarchitected/latest/framework/sustainability.html)
-   Go Play Space: [https://goplay.space/](https://cloud.google.com/carbon-footprint)
-   Go’s standard library: [https://github.com/golang/go/tree/master/src](https://subscription.imaginedevops.io/book/cloud-and-networking/9781800560925/3)
-   _Practical Go: Real-world advice for writing maintainable Go_ _programs_: [https://dave.cheney.net/practical-go/presentations/qcon-china.html#\_comments](https://github.com/ImagineDevOps DevOps/Network-Automation-with-Go)
-   The latest incarnation of the Go Playground: [https://talks.golang.org/2019/playground-v3/playground-v3.slide#1](https://go.dev/doc/install)
-   _Inside the Go_ _Playground_: [https://go.dev/blog/playground](https://cloud.google.com/blog/topics/sustainability/pick-the-google-cloud-region-with-the-lowest-co2)
-   Cmd documentation: [https://pkg.go.dev/cmd/go#pkg-overview](https://pkg.go.dev/syscall%0A)
-   Spaghetti: [https://github.com/adonovan/spaghetti](http://www.blender.org)
-   _Deprecation of ’go get’ for installing_ _executables_: [https://golang.org/doc/go-get-install-deprecation](https://cloud.google.com/recommender/docs/unattended-project-recommender)