# Processes and Signals

In the previous chapter, we talked about many interesting topics including working with Unix system files, dealing with dates and times in Go, finding information about file permissions and users as well as regular expressions and pattern matching.

The central subject of this chapter is developing Go applications that can handle the Unix signals that can be caught and handled. Go offers the os/signal package for dealing with signals, which uses Go channels. Although channels are fully explored in the next chapter, this will not stop you from learning how to work with Unix signals in Go programs.

Furthermore, you will learn how to create Go command-line utilities that can work with Unix pipes, how to draw bar charts in Go, and how to implement a Go version of the cat(1) utility. So, in this chapter you will learn about the following topics:

-   Listing the processes of a Unix machine
-   Signal handling in Go
-   The signals that a Unix machine supports as well as how to use the kill(1) command to send these signals
-   Making signals do the work you want
-   Implementing a simple version of the cat(1) utility in Go
-   Plotting data in Go
-   Using pipes in order to send the output of one program to another
-   Converting a big program into two smaller ones that will cooperate with the help of Unix pipes
-   Creating a client for a Unix socket

Just Imagine

# About Unix processes and signals

Strictly speaking, a **process** is an execution environment that contains instructions, user-data and system-data parts, and other kinds of resources that are obtained during runtime, whereas a **program** is a file that contains instructions and data, which are used for initializing the instruction and user-data parts of a process.

Just Imagine

# Process management

Go is not that good at dealing with processes and process management in general. Nevertheless, this section will present a small Go program that lists all the processes of a Unix machine by executing a Unix command and getting its output. The name of the program will be listProcess.go. It works on both Linux and macOS systems, and will be presented in three parts.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "os/exec" 
   "syscall" 
) 
```

The second part of listProcess.go has the following Go code:

```markup
func main() { 
 
   PS, err := exec.LookPath("ps") 
   if err != nil { 
         fmt.Println(err) 
   } 
fmt.Println(PS) 
 
   command := []string{"ps", "-a", "-x"} 
   env := os.Environ() 
   err = syscall.Exec(PS, command, env) 
```

As you can see, you first need to get the path of the executable file using exec.LookPath() to make sure that you are not going to accidentally execute another binary file and then define the command you want to execute, including the parameters of the command, using a slice. Next, you will have to read the Unix environment using os.Environ(). Also, you execute the desired command using syscall.Exec(), which will automatically print its output, which is not a very elegant way to execute commands because you have no control over the task and because you are calling processes at the lowest level instead of using a higher level library such as os/exec.

The last part of the program is for printing the error message of the previous code, in case there is one:

```markup
   if err != nil { 
         fmt.Println(err) 
   } 
} 
```

Executing listProcess.go will generate the following output: the head(1) utility is used to get a smaller output:

```markup
$ go run listProcess.go | head -3
/bin/ps
  PID TTY           TIME CMD
    1 ??         0:30.72 /sbin/launchd
signal: broken pipe
```

# About Unix signals

Have you ever pressed _Ctrl_ + _C_ in order to stop a program from running? If yes, then you are already familiar with signals because _Ctrl_ + _C_ sends the SIGINT signal to the program.

Strictly speaking, Unix **signals** are software interrupts that can be accessed either by a name or number and offer a way of handling asynchronous events such as when a child process exits or a process is told to pause on a Unix system.

A program cannot handle all signals; some of them are non-catchable and non-ignorable. The SIGKILL and SIGSTOP signals cannot be caught, blocked, or ignored. The reason for this is that they provide the kernel and the root user a way of stopping any process. The SIGKILL signal, which is also known by the number 9, is usually called in extreme conditions where you need to act fast; so, it is the only signal that is usually called by number because it is quicker to do so. The most important thing to remember here is that not all Unix signals can be handled!

Just Imagine

# Unix signals in Go

Go provides the os/signal package to programmers to help them handle incoming signals. However, we will start the discussion about handling by presenting the kill(1) utility.

Just Imagine

# The kill(1) command

The kill(1) command is used for either terminating a process or sending a less cruel signal to it. Keep in mind that the fact that you can send a signal to a process does not mean that the process can or has code to handle this signal.

By default, kill(1) sends the SIGTERM signal. If you want to find out all the supported signals of your Unix machine, you should execute the kill -l command. On a macOS Sierra machine, the output of kill -l is the following:

```markup
$ kill -l
1) SIGHUP   2) SIGINT        3) SIGQUIT   4) SIGILL
5) SIGTRAP  6) SIGABRT       7) SIGEMT    8) SIGFPE
9) SIGKILL 10) SIGBUS        11) SIGSEGV 12) SIGSYS
13) SIGPIPE 14) SIGALRM       15) SIGTERM 16) SIGURG
17) SIGSTOP 18) SIGTSTP       19) SIGCONT 20) SIGCHLD
21) SIGTTIN 22) SIGTTOU       23) SIGIO   24) SIGXCPU
25) SIGXFSZ 26) SIGVTALRM     27) SIGPROF 28) SIGWINCH
29) SIGINFO 30) SIGUSR1       31) SIGUSR2
```

If you execute the same command on a Debian Linux machine, you will get a more enriched output:

```markup
$ kill -l
 1) SIGHUP   2) SIGINT   3) SIGQUIT  4) SIGILL   5) SIGTRAP
 6) SIGABRT  7) SIGBUS   8) SIGFPE   9) SIGKILL 10) SIGUSR1
11) SIGSEGV 12) SIGUSR2 13) SIGPIPE 14) SIGALRM 15) SIGTERM
16) SIGSTKFLT     17) SIGCHLD 
18) SIGCONT       19) SIGSTOP 20) SIGTSTP
21) SIGTTIN       22) SIGTTOU 
23) SIGURG        24) SIGXCPU 25) SIGXFSZ
26) SIGVTALRM     27) SIGPROF 28) SIGWINCH      
29) SIGIO         30) SIGPWR
31) SIGSYS        34) SIGRTMIN      
35) SIGRTMIN+1    36) SIGRTMIN+2    37) SIGRTMIN+3
38) SIGRTMIN+4    39) SIGRTMIN+5    
40) SIGRTMIN+6    41) SIGRTMIN+7    42) SIGRTMIN+8
43) SIGRTMIN+9    44) SIGRTMIN+10   
45) SIGRTMIN+11   46) SIGRTMIN+12   47) SIGRTMIN+13
48) SIGRTMIN+14   49) SIGRTMIN+15   
50) SIGRTMAX-14   51) SIGRTMAX-13   52) SIGRTMAX-12
53) SIGRTMAX-11   54) SIGRTMAX-10   
55) SIGRTMAX-9    56) SIGRTMAX-8    57) SIGRTMAX-7
58) SIGRTMAX-6    59) SIGRTMAX-5    
60) SIGRTMAX-4    61) SIGRTMAX-3    62) SIGRTMAX-2
63) SIGRTMAX-1    64) SIGRTMAX
```

If you try to kill or send another signal to the process of another user without having the required permissions, which most likely will happen if you are not the _root_ user, kill(1) will not do the job and you will get an error message similar to the following:

```markup
$ kill 2908
-bash: kill: (2908) - Operation not permitted
```

# A simple signal handler in Go

This subsection will present a naive Go program that handles only the SIGTERM and SIGINT signals. The Go code of h1s.go will be presented in three parts; the first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "os/signal" 
   "syscall" 
   "time" 
) 
 
func handleSignal(signal os.Signal) { 
   fmt.Println("Got", signal) 
} 
```

Apart from the preamble of the program, there is also a function named handleSignal() that will be called when the program receives any of the two supported signals.

The second part of h1s.go contains the following Go code:

```markup
func main() { 
   sigs := make(chan os.Signal, 1) 
   signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) 
   go func() { 
         for { 
               sig := <-sigs 
               fmt.Println(sig) 
               handleSignal(sig) 
         } 
   }() 
```

The previous code uses a **goroutine** and a Go **channel**, which are Go features that have not been discussed in this book. Unfortunately, you will have to wait until [Chapter 9](https://subscription.imaginedevops.io/book/programming/9781787125643/9)_,_ _Goroutines - Basic Features_, to learn more about both of them. Note that although os.Interrupt and syscall.SIGTERM belong to different Go packages, they are both signals.

For now, understanding the technique is important; it includes three steps:

1.  The definition of a channel, which acts as a way of passing data around, that is required for the technique (sigs).
2.  Calling signal.Notify() in order to define the list of signals you want to be able to catch.
3.  The definition of an anonymous function that runs in a goroutine (go func()) right after signal.Notify(), which is used for deciding what you are going to do when you get any of the desired signals.

In this case, the handleSignal() function will be called. The for loop inside the anonymous function is used to make the program to keep handling all signals and not stop after receiving its first signal.

The last part of h1s.go is the following:

```markup
   for { 
         fmt.Printf(".") 
         time.Sleep(10 * time.Second) 
   } 
} 
```

This is an endless for loop that delays the ending of the program forever: in its place you would most likely put the actual code of your program. Executing h1s.go and sending signals to it from another Terminal will make h1s.go generate the following output:

```markup
$ ./h1s
......................^CinterruptGot interrupt
^Cinterrupt
Got interrupt
.Hangup: 1
```

The bad thing here is that h1s.go will stop when it receives the SIGHUP signal because the default action for SIGHUP when it is not being specifically handled by a program is to kill the process! The next subsection will show how to handle three signals in a better way, and the subsection after that will teach you how to handle all signals that can be handled.

# Handling three different signals!

This subsection will teach you how to create a Go application that can handle three different signals: the name of the program will be h2s.go, and it will handle the SIGTERM, SIGINT, and SIGHUP signals.

The Go code of h2s.go will be presented in four parts.

The first part of the program contains the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "os/signal" 
   "syscall" 
   "time" 
) 
```

The second part has the following Go code:

```markup
func handleSignal(signal os.Signal) { 
   fmt.Println("* Got:", signal) 
} 
 
func main() { 
   sigs := make(chan os.Signal, 1) 
   signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP) 
```

Here, the last statement tells you that the program will only handle the os.Interrupt, syscall.SIGTERM, and syscall.SIGHUP signals.

The third part of h2s.go is the following:

```markup
   go func() { 
         for { 
               sig := <-sigs 
               switch sig { 
               case os.Interrupt: 
                     handleSignal(sig) 
               case syscall.SIGTERM: 
                     handleSignal(sig) 
               case syscall.SIGHUP: 
                     fmt.Println("Got:", sig) 
                     os.Exit(-1) 
               } 
         } 
   }() 
```

Here, you can see that it is not compulsory to call a separate function when a given signal is caught; it is also allowed to handle it inside the for loop as it happens with syscall.SIGHUP. However, I find the use of a named function better because it makes the Go code easier to read and modify. The good thing is that Go has a central place for handling all signals, which makes it easy to find out what is going on with your program.

Additionally, h2s.go specifically handles the SIGHUP signal, although a SIGHUP signal will still terminate the program; however, this time this is our decision.

Keep in mind that it is considered good practice to make one of the signal handlers to stop the program because otherwise you will have to terminate it by issuing a kill -9 command.

The last part of h2s.go is the following:

```markup
   for { 
         fmt.Printf(".") 
         time.Sleep(10 * time.Second) 
   } 
}
```

Executing h2s.go and sending four signals to it (SIGINT, SIGTERM, SIGHUP, and SIGKILL) from another shell will generate the following output:

```markup
$ go build h2s.go
$ ./h2s
..* Got: interrupt
* Got: terminated
.Got: hangup
.Killed: 9
```

The reason for building h2s.go is that it is easier to find the process ID of an autonomous program: the go run command builds a temporary executable program behind the scenes, which in this case offers less flexibility. If you want to improve h2s.go, you can make it call os.Getpid() in order to print its process ID, which will save you from having to find it on your own.

The program handles three signals before getting a SIGKILL that cannot be handled and therefore terminates it!

# Catching every signal that can be handled

This subsection will present a simple technique that allows you to catch every signal that can be handled: once again, you should keep in mind that you cannot handle all signals! The program will stop once it gets a SIGTERM signal.

The name of the program will be catchAll.go and will be presented in three parts.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "os/signal" 
   "syscall" 
   "time" 
) 
 
func handleSignal(signal os.Signal) { 
   fmt.Println("* Got:", signal) 
} 
```

The second part of the program is the following:

```markup
func main() { 
   sigs := make(chan os.Signal, 1) 
   signal.Notify(sigs) 
   go func() { 
         for { 
               sig := <-sigs 
               switch sig { 
               case os.Interrupt: 
                     handleSignal(sig) 
               case syscall.SIGTERM: 
                     handleSignal(sig) 
                     os.Exit(-1) 
               case syscall.SIGUSR1: 
                     handleSignal(sig) 
               default: 
                     fmt.Println("Ignoring:", sig) 
               } 
         } 
   }() 
```

In this case, all the difference is made by the way you call signal.Notify() in your code. As you do not define any particular signals, the program will be able to handle any signal that can be handled. However, the for loop inside the anonymous function only takes care of three signals while ignoring the remaining ones! Note that I believe that this is the best way to handle signals in Go: catch everything while processing only the signals that interest you. However, some people believe that being explicit about what you handle is a better approach. There is no right or wrong here.

The catchAll.go program will not terminate when it gets SIGHUP because the default case of the switch block handles it.

The last part is the expected call to the time.Sleep() function:

```markup
   for { 
         fmt.Printf(".") 
         time.Sleep(10 * time.Second) 
   } 
} 
```

Executing catchAll.go will create the following output:

```markup
$ ./catchAll
.Ignoring: hangup.......................................* Got: interrupt
* Got: user defined signal 1
.Ignoring: user defined signal 2
Ignoring: hangup
.* Got: terminated
$
```

# Rotating log files revisited!

As I told you back in [](https://subscription.imaginedevops.io/book/programming/9781787125643/7)[Chapter 7](https://subscription.imaginedevops.io/book/programming/9781787125643/7)_,_ _Working with System Files_, this chapter will present you with a technique that will allow you to end the program and rotate log files in a more conventional way with the help of signals and signal handling.

The name of the new version of rotateLog.go will be rotateSignals.go and will be presented in four parts. Moreover, when the utility receives os.Interrupt, it will rotate the current log file, whereas when it receives syscall.SIGTERM, it will terminate its execution. Every other signal that can be handled will create a log entry without any other action.

The first part of the rotateSignals.go is the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "log" 
   "os" 
   "os/signal" 
   "strconv" 
   "syscall" 
   "time" 
) 
 
var TOTALWRITES int = 0 
var openLogFile os.File 
```

The second part of rotateSignals.go has the following Go code:

```markup
func rotateLogFile(filename string) error { 
   openLogFile.Close() 
   os.Rename(filename, filename+"."+strconv.Itoa(TOTALWRITES)) 
   err := setUpLogFile(filename) 
   return err 
} 
 
func setUpLogFile(filename string) error { 
   openLogFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644) 
   if err != nil { 
         return err 
   } 
   log.SetOutput(openLogFile) 
   return nil 
} 
```

You have just defined two functions here that perform two tasks. The third part of rotateSignals.go contains the following Go code:

```markup
func main() { 
   filename := "/tmp/myLog.log" 
   err := setUpLogFile(filename) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(-1) 
   } 
 
   sigs := make(chan os.Signal, 1) 
   signal.Notify(sigs) 
```

Once again, all signals will be caught. The last part of rotateSignals.go is the following:

```markup
   go func() { 
         for { 
               sig := <-sigs 
               switch sig { 
               case os.Interrupt: 
                     rotateLogFile(filename) 
                     TOTALWRITES++ 
               case syscall.SIGTERM: 
                     log.Println("Got:", sig) 
                     openLogFile.Close() 
                     TOTALWRITES++ 
                     fmt.Println("Wrote", TOTALWRITES, "log entries in total!") 
                     os.Exit(-1) 
               default: 
                     log.Println("Got:", sig) 
                     TOTALWRITES++ 
               } 
         } 
   }() 
 
   for { 
         time.Sleep(10 * time.Second) 
   } 
} 
```

As you can see, rotateSignals.go records information about the signals it has received by writing one log entry for each signal. Although presenting the entire code of rotateSignals.go is good, it would be very educational to see the output of the diff(1) utility to show the code differences between rotateLog.go and rotateSignals.go:

```markup
$ diff rotateLog.go rotateSignals.go
6a7
>     "os/signal"
7a9
>     "syscall"
12,13d13
< var ENTRIESPERLOGFILE int = 100
< var WHENTOSTOP int = 230
33d32
<     numberOfLogEntries := 0
41,51c40,59
<     for {
<           log.Println(numberOfLogEntries, "This is a test log entry")
<           numberOfLogEntries++
<           TOTALWRITES++
<           if numberOfLogEntries > ENTRIESPERLOGFILE {
<                 _ = rotateLogFile(filename)
<                 numberOfLogEntries = 0
<           }
<           if TOTALWRITES > WHENTOSTOP {
<                 _ = rotateLogFile(filename)
<                 break
---
>     sigs := make(chan os.Signal, 1)
>     signal.Notify(sigs)
>
>     go func() {
>           for {
>                 sig := <-sigs
>                 switch sig {
>                 case os.Interrupt:
>                       rotateLogFile(filename)
>                       TOTALWRITES++
>                 case syscall.SIGTERM:
>                       log.Println("Got:", sig)
>                       openLogFile.Close()
>                       TOTALWRITES++
>                       fmt.Println("Wrote", TOTALWRITES, "log entries in total!")
>                       os.Exit(-1)
>                 default:
>                       log.Println("Got:", sig)
>                       TOTALWRITES++
>                 }
53c61,64
<           time.Sleep(time.Second)
---
>     }()
>
>     for {
>           time.Sleep(10 * time.Second)
55d65
<     fmt.Println("Wrote", TOTALWRITES, "log entries!")
```

The good thing here is that the use of signals in rotateSignals.go makes most of the global variables used in rotateLog.go unnecessary because you can now control the utility by sending signals. Additionally, the design and the structure of rotateSignals.go are simpler than rotateLog.go because you only have to understand what the anonymous function does.

After executing rotateSignals.go and sending some signals to it, the contents of /tmp/myLog.log will look like the following:

```markup
$ cat /tmp/myLog.log
2017/06/03 14:53:33 Got: user defined signal 1
2017/06/03 14:54:08 Got: user defined signal 1
2017/06/03 14:54:12 Got: user defined signal 2
2017/06/03 14:54:19 Got: terminated
```

Additionally, you will have the following files inside /tmp:

```markup
$ ls -l /tmp/myLog.log*
-rw-r--r--  1 mtsouk  wheel  177 Jun  3 14:54 /tmp/myLog.log
-rw-r--r--  1 mtsouk  wheel  106 Jun  3 13:42 /tmp/myLog.log.0
```

Just Imagine

# Improving file copying

The original cp(1) utility prints useful information when it receives a SIGINFO signal, as shown in the following output:

```markup
$ cp FileToCopy /tmp/copy
FileToCopy -> /tmp/copy  26%
FileToCopy -> /tmp/copy  29%
FileToCopy -> /tmp/copy  31%
```

So, the rest of this section will implement the same functionality to the Go implementation of the cp(1) command. The Go code in this section will be based on the cp.go program because it can be very slow when used with a small buffer size giving us time for testing. The name of the new copy utility will be cpSignal.go and will be presented in four parts.

The fundamental difference between cpSignal.go and cp.go is that cpSignal.go should find the size of the input file and keep the number of bytes that have been written at a given point. Apart from those modifications there is nothing else that you should worry about because the core functionality of the two versions, which is copying a file, is exactly the same.

The first part of the program is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "os" 
   "os/signal" 
   "path/filepath" 
   "strconv" 
   "syscall" 
) 
 
var BUFFERSIZE int64 
var FILESIZE int64 
var BYTESWRITTEN int64 
```

In order to make things simpler for the developer, the program introduces two global variables called FILESIZE and BYTESWRITTEN and these keep the size of the input file and the number of bytes that have been written, respectively. Both variables are used by the function that handles the SIGINFO signal.

The second part is as follows:

```markup
func Copy(src, dst string, BUFFERSIZE int64) error { 
   sourceFileStat, err := os.Stat(src) 
   if err != nil { 
         return err 
   } 
 
   FILESIZE = sourceFileStat.Size() 
 
   if !sourceFileStat.Mode().IsRegular() { 
         return fmt.Errorf("%s is not a regular file.", src) 
   } 
 
   source, err := os.Open(src) 
   if err != nil { 
         return err 
   } 
   defer source.Close() 
 
   _, err = os.Stat(dst) 
   if err == nil { 
         return fmt.Errorf("File %s already exists.", dst) 
   } 
 
   destination, err := os.Create(dst) 
   if err != nil { 
         return err 
   } 
   defer destination.Close() 
 
   if err != nil { 
         panic(err) 
   } 
 
   buf := make([]byte, BUFFERSIZE) 
   for { 
         n, err := source.Read(buf) 
         if err != nil && err != io.EOF { 
               return err 
         } 
         if n == 0 { 
               break 
         } 
         if _, err := destination.Write(buf[:n]); err != nil { 
               return err 
         } 
         BYTESWRITTEN = BYTESWRITTEN + int64(n) 
   } 
   return err 
} 
```

Here, you use the sourceFileStat.Size() function to get the size of the input file and set the value of the FILESIZE global variable.

The third part is where you define the signal handling:

```markup
func progressInfo() { 
   progress := float64(BYTESWRITTEN) / float64(FILESIZE) * 100 
   fmt.Printf("Progress: %.2f%%\n", progress) 
} 
 
func main() { 
   if len(os.Args) != 4 { 
         fmt.Printf("usage: %s source destination BUFFERSIZE\n", filepath.Base(os.Args[0])) 
         os.Exit(1) 
   } 
 
   source := os.Args[1] 
   destination := os.Args[2] 
   BUFFERSIZE, _ = strconv.ParseInt(os.Args[3], 10, 64) 
   BYTESWRITTEN = 0 
 
   sigs := make(chan os.Signal, 1) 
   signal.Notify(sigs) 
```

Here, you choose to catch all signals. However, the Go code of the anonymous function will only call progressInfo() after receiving a syscall.SIGINFO signal.

If you want to have a way of gracefully terminating the program, you might want to use the SIGINT signal because when capturing all signals, gracefully terminating a program is no longer possible: you will need to send a SIGKILL in order to terminate your program, which is a little cruel.

The last part of cpSignal.go is the following:

```markup
   go func() { 
         for {               sig := <-sigs 
               switch sig { 
               case syscall.SIGINFO:
                     progressInfo() 
               default: 
                     fmt.Println("Ignored:", sig) 
               } 
         } 
   }() 
 
   fmt.Printf("Copying %s to %s\n", source, destination) 
   err := Copy(source, destination, BUFFERSIZE) 
   if err != nil { 
         fmt.Printf("File copying failed: %q\n", err) 
   } 
} 
```

Executing cpSignal.go and sending two SIGINFO signals to it will generate the following output:

```markup
$ ./cpSignal FileToCopy /tmp/copy 2
Copying FileToCopy to /tmp/copy
Ignored: user defined signal 1
Progress: 21.83%
^CIgnored: interrupt
Progress: 29.78%
```

Just Imagine

# Plotting data

The utility that will be developed in this section will read multiple log files and will create a graphical image with as many bars as the number of log files read. Each bar will represent the number of times a given IP address has been found in a log file.

However, the Unix philosophy tells us that instead of developing a single utility, we should make two distinct utilities: one for processing the log files and creating a report and another for plotting the data generated by the first utility: the two utilities will communicate using Unix pipes. Although this section will implement the first approach, you will see the implementation of the second approach later in _The_ plotIP.go _utility revisited_ section of this chapter.

The idea for the presented utility came from a tutorial that I wrote for a magazine where I developed a small Go program that did some plotting: even small and naive programs can inspire you to develop bigger things, so do not underestimate their power.

The name of the utility will be plotIP.go, and it will be presented in seven parts: the good thing is that plotIP.go will reuse some of the code of countIP.go and findIP.go. The only thing that plotIP.go does not do is writing text to the image, so you can only plot the bars without knowing the actual values or the corresponding log file of a particular bar: you can try to add text capabilities to the program as an exercise.

Also, plotIP.go will require at least three parameters, which are the width and height of the image and the name of the log file that will be used: in order to make plotIP.go smaller, plotIP.go will not use the flag package and assume that you will give its parameters in the correct order. If you give it more parameters, it will consider them as log files.

The first part of plotIP.go is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "image" 
   "image/color" 
   "image/png" 
   "io" 
   "os" 
   "path/filepath" 
   "regexp" 
   "strconv" 
) 
var m *image.NRGBAvar x int 
var y int 
var barWidth int 
```

These global variables related to the dimensions of the image (x and y), the image as a Go variable (m), and the width of one of its bars (barWidth) that depends on the size of the image and the number of the bars that will be plotted. Note that using x and y as variable names instead of something like IMAGEWIDTH and IMAGEHEIGHT might be a little wrong and dangerous here.

The second part is the following:

```markup
func findIP(input string) string { 
   partIP := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])" 
   grammar := partIP + "\\." + partIP + "\\." + partIP + "\\." + partIP 
   matchMe := regexp.MustCompile(grammar) 
   return matchMe.FindString(input) 
} 
 
func plotBar(width int, height int, color color.RGBA) { 
   xx := 0   for xx < barWidth { 
         yy := 0 
         for yy < height { 
               m.Set(xx+width, y-yy, color) 
               yy = yy + 1 
         } 
         xx = xx + 1 
   } 
} 
```

Here, you implement a Go function named plotBar() that does the plotting of each bar, given its height, its width, and its color of the bar. This function is the most challenging part of plotIP.go.

The third part has the following Go code:

```markup
func getColor(x int) color.RGBA { 
   switch { 
   case x == 0: 
         return color.RGBA{0, 0, 255, 255} 
   case x == 1: 
         return color.RGBA{255, 0, 0, 255} 
   case x == 2: 
         return color.RGBA{0, 255, 0, 255} 
   case x == 3: 
         return color.RGBA{255, 255, 0, 255} 
   case x == 4: 
         return color.RGBA{255, 0, 255, 255} 
   case x == 5: 
         return color.RGBA{0, 255, 255, 255} 
   case x == 6: 
         return color.RGBA{255, 100, 100, 255} 
   case x == 7: 
         return color.RGBA{100, 100, 255, 255} 
   case x == 8: 
         return color.RGBA{100, 255, 255, 255} 
   case x == 9: 
         return color.RGBA{255, 255, 255, 255} 
   } 
   return color.RGBA{0, 0, 0, 255} 
} 
```

This function lets you define the colors that will be present in the output: you can change them if you want.

The fourth part contains the following Go code:

```markup
func main() { 
   var data []int 
   arguments := os.Args 
   if len(arguments) < 4 { 
         fmt.Printf("%s X Y IP input\n", filepath.Base(arguments[0])) 
         os.Exit(0) 
   } 
 
   x, _ = strconv.Atoi(arguments[1]) 
   y, _ = strconv.Atoi(arguments[2]) 
   WANTED := arguments[3] 
   fmt.Println("Image size:", x, y) 
```

Here, you read the desired IP address, which is saved in the WANTED variable and you read the dimensions of the generated PNG image.

The fifth part contains the following Go code:

```markup
   for _, filename := range arguments[4:] { 
         count := 0 
         fmt.Println(filename) 
         f, err := os.Open(filename) 
         if err != nil { 
               fmt.Fprintf(os.Stderr, "Error: %s\n", err) 
               continue 
         } 
         defer f.Close() 
 
         r := bufio.NewReader(f) 
         for { 
               line, err := r.ReadString('\n') 
               if err == io.EOF { 
                     break 
               } 
 
if err != nil { 
                fmt.Fprintf(os.Stderr, "Error in file: %s\n", err) 
                     continue 
               } 
               ip := findIP(line) 
               if ip == WANTED { 
                     count++ 
               } 
         } 
         data = append(data, count) 
   } 
```

Here, you process the input log files one by one and store the values you calculate in the data slice. Error messages are printed to os.Stderr: the main advantage you get from printing error messages to os.Stderr is that you can easily redirect error messages to a file while using data written to os.Stdout in a different way.

The sixth part of plotIP.go contains the following Go code:

```markup
   fmt.Println("Slice length:", len(data)) 
   if len(data)*2 > x { 
         fmt.Println("Image size (x) too small!") 
         os.Exit(-1) 
   } 
 
   maxValue := data[0] 
   for _, temp := range data { 
         if maxValue < temp { 
               maxValue = temp 
         } 
   } 
 
   if maxValue > y { 
         fmt.Println("Image size (y) too small!") 
         os.Exit(-1) 
   } 
   fmt.Println("maxValue:", maxValue) 
   barHeighPerUnit := int(y / maxValue) 
   fmt.Println("barHeighPerUnit:", barHeighPerUnit) 
   PNGfile := WANTED + ".png" 
   OUTPUT, err := os.OpenFile(PNGfile, os.O_CREATE|os.O_WRONLY, 0644) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(-1) 
   } 
   m = image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{x, y}}) 
```

Here, you calculate things about the plot and create the output image file using os.OpenFile(). The PNG file generated by the plotIP.go utility is named after the given IP address to make things simpler.

The last part of the Go code of plotIP.go is the following:

```markup
   i := 0 
   barWidth = int(x / len(data)) 
   fmt.Println("barWidth:", barWidth) 
   for _, v := range data { 
         c := getColor(v % 10) 
         yy := v * barHeighPerUnit 
         plotBar(barWidth*i, yy, c) 
         fmt.Println("plotBar", barWidth*i, yy) 
         i = i + 1 
   } 
   png.Encode(OUTPUT, m) 
} 
```

Here, you read the values of the data slice and create a bar for each one of them by calling the plotBar() function.

Executing plotIP.go will generate the following output:

```markup
$ go run plotIP.go 1300 1500 127.0.0.1 /tmp/log.*
Image size: 1300 1500
/tmp/log.1
/tmp/log.2
/tmp/log.3
Slice length: 3
maxValue: 1500
barHeighPerUnit: 1
barWidth: 433
plotBar 0 1500
plotBar 433 1228
plotBar 866 532
$  ls -l 127.0.0.1.png
-rw-r--r-- 1 mtsouk mtsouk 11023 Jun  5 18:36 127.0.0.1.png
```

However, apart from the generated text output, what is important is the produced PNG file that can be seen in the following figure:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/0705a55e-044d-4918-bfea-70d6b7d9377e.png)

The output generated by the plotIP.go utility

If you want to save the error messages to a different file, you can use a variation of the following command:

```markup
$ go run plotIP.go 130 150 127.0.0.1 doNOTExist 2> err
Image size: 130 150
doNOTExist
Slice length: 0
$ cat err
Error: open doNOTExist: no such file or directory
panic: runtime error: index out of range
    
goroutine 1 [running]:
main.main()
     /Users/mtsouk/Desktop/goCourse/ch/ch8/code/plotIP.go:112 +0x12de
exit status 2
```

The following command discards all error messages by sending them to /dev/null:

```markup
$ go run plotIP.go 1300 1500 127.0.0.1 doNOTExist 2>/dev/null
Image size: 1300 1500
doNOTExist
Slice length: 0  
```

Just Imagine

# Unix pipes in Go

We first talked about pipes in _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_,_ _File Input and Output_. Pipes have two serious limitations: first, they usually communicate in one direction, and second, they can only be used between processes that have a common ancestor.

The general idea behind pipes is that if you do not have a file to process, you should wait to get your input from standard input. Similarly, if you are not told to save your output to a file, you should write your output to standard output, either for the user to see it or for another program to process it. As a result, pipes can be used for streaming data between two processes without creating any temporary files.

This section will present some simple utilities written in Go that use Unix pipes for clarity.

# Reading from standard input

The first thing that you need to know in order to develop Go applications that support Unix pipes is how to read from standard input.

The developed program is named readSTDIN.go and will be presented in three parts.

The first part of the program is the expected preamble:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "os" 
) 
```

The second part of readSTDIN.go has the following Go code:

```markup
func main() { 
   filename := "" 
   var f *os.File 
   arguments := os.Args 
   if len(arguments) == 1 { 
         f = os.Stdin 
   } else { 
         filename = arguments[1] 
         fileHandler, err := os.Open(filename) 
         if err != nil { 
               fmt.Printf("error opening %s: %s", filename, err) 
               os.Exit(1) 
         } 
         f = fileHandler 
   } 
   defer f.Close() 
```

Here, you resolve whether you have an actual file to process, which can be determined by the number of the command-line arguments of your program. If you do not have a file to process, you will try to read data from os.Stdin. Make sure that you understand the presented technique because it will be used many times in this chapter.

The last part of readSTDIN.go is the following:

```markup
   scanner := bufio.NewScanner(f) 
   for scanner.Scan() { 
         fmt.Println(">", scanner.Text()) 
   } 
} 
```

This code is the same whether you are processing an actual file or os.Stdin, which happens because everything in Unix is a file. Note that the program output begins with the \> character.

Executing readSTDIN.go will generate the following output:

```markup
$ cat /tmp/testfile
1
2
$ go run readSTDIN.go /tmp/testFile
> 1
> 2
$ cat /tmp/testFile | go run readSTDIN.go
> 1
> 2
$ go run readSTDIN.go
3
> 3
2
> 2
1
> 1
```

In the last case, readSTDIN.go echoes each line it reads because the input is read line by line: the cat(1) utility works the same way.

# Sending data to standard output

This subsection will show you how to send data to standard output in a better way than just using fmt.Println() or any other function from the fmt standard Go package. The Go program will be named writeSTDOUT.go and will be presented to you in three parts.

The first part is the following:

```markup
package main 
 
import ( 
   "io" 
   "os" 
) 
```

The second part of writeSTDOUT.go has the following Go code:

```markup
func main() { 
   myString := "" 
   arguments := os.Args 
   if len(arguments) == 1 { 
         myString = "You did not give an argument!" 
   } else { 
         myString = arguments[1] 
   } 
```

The last part of writeSTDOUT.go is the following:

```markup
   io.WriteString(os.Stdout, myString) 
   io.WriteString(os.Stdout, "\n") 
} 
```

The only subtle thing is that you need to put your text into a slice before using io.WriteString() to write data to os.Stdout.

Executing writeSTDOUT.go will generate the following output:

```markup
$ go run writeSTDOUT.go 123456
123456
$ go run writeSTDOUT.go
You do not give an argument!
```

# Implementing cat(1) in Go

This subsection will present a Go version of the cat(1) command-line utility. If you give one or more command-line arguments to cat(1), then cat(1) will print their contents on the screen. However, if you just type cat(1) on your Unix shell, then cat(1) will wait for your input, which will be terminated when you type _Ctrl_ + _D_.

The name of the Go implementation will be cat.go and will be presented in three parts.

The first part of cat.go is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io" 
   "os" 
) 
```

The second part is the following:

```markup
func catFile(filename string) error { 
   f, err := os.Open(filename) 
   if err != nil { 
         return err 
   } 
   defer f.Close() 
   scanner := bufio.NewScanner(f) 
   for scanner.Scan() { 
         fmt.Println(scanner.Text()) 
   } 
   return nil 
} 
```

The catFile() function is called when the cat.go utility has to process real files. Having a function to do your job makes the design of the program better.

The last part has the following Go code:

```markup
func main() { 
   filename := "" 
   arguments := os.Args 
   if len(arguments) == 1 { 
         io.Copy(os.Stdout, os.Stdin) 
         os.Exit(0) 
   } 
 
   filename = arguments[1] 
   err := catFile(filename) 
   if err != nil { 
         fmt.Println(err) 
   } 
} 
```

So, if the program has no arguments, then it assumes that it has to read os.Stdin. In that case, it just echoes each line you give to it. If the program has arguments, then it processes the first argument as a file using the catFile() function.

Executing cat.go will generate the following output:

```markup
$ go run cat.go /tmp/testFile  |  go run cat.go
1
2
$ go run cat.go
Mihalis
Mihalis
Tsoukalos
Tsoukalos$ echo "Mihalis Tsoukalos" | go run cat.go
Mihalis Tsoukalos
```

# The plotIP.go utility revisited

As promised in a previous section of this chapter, this section will create two separate utilities, which when combined will implement the functionality of plotIP.go. Personally, I prefer to have two separate utilities and combine them when needed than having just one utility that does two or more tasks.

The names of the two utilities will be extractData.go and plotData.go. As you can easily understand, only the second utility would have to be able to get input from standard input as long as the first utility prints its output on standard output either using os.Stdout, which is the correct way, or using fmt.Println(), which usually does the job.

I think that I should now tell you my little secret: I created extractData.go and plotData.go first and then developed plotIP.go because it is easier to develop two separate utilities than a bigger one that does everything! Additionally, the use of two different utilities allows you to filter the output of extractData.go using standard Unix utilities such as tail(1), sort(1), and head(1), which means that you can modify your data in different ways without the need for writing any extra Go code.

Taking two command-line utilities and creating one utility that implements the functionality of both utilities is easier than taking one big utility and dividing its functionality into two or more distinct utilities because the latter usually requires more variables and more error checking.

The extractData.go utility will be presented in four parts; the first part is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "io" 
   "os" 
   "path/filepath" 
   "regexp" 
) 
```

The second part of extractData.go has the following Go code:

```markup
func findIP(input string) string { 
   partIP := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])" 
   grammar := partIP + "\\." + partIP + "\\." + partIP + "\\." + partIP 
   matchMe := regexp.MustCompile(grammar) 
   return matchMe.FindString(input) 
} 
```

You should be familiar with the findIP() function, which you saw in findIP.go in _[](https://subscription.imaginedevops.io/book/programming/9781787125643/7)_[Chapter 7](https://subscription.imaginedevops.io/book/programming/9781787125643/7)_,_ _Working with System files_.

The third part of extractData.go is the following:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) < 3 { 
         fmt.Printf("%s IP <files>\n", filepath.Base(os.Args[0])) 
         os.Exit(-1) 
   } 
 
   WANTED := arguments[1] 
   for _, filename := range arguments[2:] { 
         count := 0 
         buf := []byte(filename)         io.WriteString(os.Stdout, string(buf)) 
         f, err := os.Open(filename) 
         if err != nil { 
               fmt.Fprintf(os.Stderr, "Error: %s\n", err) 
               continue 
         } 
         defer f.Close() 
```

The use of the buf variable is redundant here because filename is a string and io.WriteString() expects a string: it is just my habit to put the value of filename into a byte slice. You can remove it if you want.

Once again, most of the Go code is from the plotIP.go utility. The last part of extractData.go is the following:

```markup
         r := bufio.NewReader(f) 
         for { 
               line, err := r.ReadString('\n') 
               if err == io.EOF { 
                     break 
               } else if err != nil { 
                     fmt.Fprintf(os.Stderr, "Error in file: %s\n", err) 
                     continue 
               } 
 
               ip := findIP(line) 
               if ip == WANTED { 
                     count = count + 1 
               } 
         } 
         buf = []byte(strconv.Itoa(count))         io.WriteString(os.Stdout, " ") 
         io.WriteString(os.Stdout, string(buf)) 
         io.WriteString(os.Stdout, "\n") 
   } 
} 
```

Here, extractData.go writes its output to standard output (os.Stdout) instead of using the functions of the fmt package in order to be more compatible with pipes. The extractData.go utility requires at least two parameters: an IP address and a log file, but it can process as many log files as you wish.

You might want to move the printing of the filename value from the third part here in order to have all printing commands at the same place.

Executing extractData.go will generate the following output:

```markup
$ ./extractData 127.0.0.1 access.log{,.1}
access.log 3099
access.log.1 6333
```

Although extractData.go prints two values in each line, only the second field will be used by plotData.go. The best way to do that is filter the output of extractData.go using awk(1):

```markup
$ ./extractData 127.0.0.1 access.log{,.1} | awk '{print $2}'
3099
6333
```

As you can understand, awk(1) allows you to do many more things with the generated values.

The plotData.go utility will also be presented in six parts; its first part is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "image" 
   "image/color" 
   "image/png" 
   "os" 
   "path/filepath" 
   "strconv" 
) 
 
var m *image.NRGBA 
var x int 
var y int 
var barWidth int 
```

Once again, the use of global variables is for avoiding the passing of too many arguments to some of the functions of the utility.

The second part of plotData.go contains the following Go code:

```markup
func plotBar(width int, height int, color color.RGBA) { 
   xx := 0   for xx < barWidth { 
         yy := 0 
         for yy < height { 
               m.Set(xx+width, y-yy, color) 
               yy = yy + 1 
         } 
         xx = xx + 1 
   } 
} 
```

The third part of plotData.go has the following Go code:

```markup
func getColor(x int) color.RGBA { 
   switch {   case x == 0: 
         return color.RGBA{0, 0, 255, 255} 
   case x == 1: 
         return color.RGBA{255, 0, 0, 255} 
   case x == 2: 
         return color.RGBA{0, 255, 0, 255} 
   case x == 3: 
         return color.RGBA{255, 255, 0, 255} 
   case x == 4: 
         return color.RGBA{255, 0, 255, 255} 
   case x == 5: 
         return color.RGBA{0, 255, 255, 255} 
   case x == 6: 
         return color.RGBA{255, 100, 100, 255} 
   case x == 7: 
         return color.RGBA{100, 100, 255, 255} 
   case x == 8: 
         return color.RGBA{100, 255, 255, 255} 
   case x == 9: 
         return color.RGBA{255, 255, 255, 255} 
   } 
   return color.RGBA{0, 0, 0, 255} 
} 
```

The fourth part of plotData.go contains the following Go code:

```markup
func main() { 
   var data []int 
   var f *os.File 
   arguments := os.Args 
   if len(arguments) < 3 { 
         fmt.Printf("%s X Y input\n", filepath.Base(arguments[0])) 
         os.Exit(0) 
   } 
 
   if len(arguments) == 3 { 
         f = os.Stdin 
   } else { 
         filename := arguments[3] 
         fTemp, err := os.Open(filename) 
         if err != nil { 
               fmt.Println(err) 
               os.Exit(0) 
         } 
         f = fTemp 
   } 
   defer f.Close() 
 
   x, _ = strconv.Atoi(arguments[1]) 
   y, _ = strconv.Atoi(arguments[2]) 
   fmt.Println("Image size:", x, y) 
```

The fifth part of plotData.go is the following:

```markup
   scanner := bufio.NewScanner(f) 
   for scanner.Scan() { 
         value, err := strconv.Atoi(scanner.Text()) 
         if err == nil { 
               data = append(data, value) 
         } else { 
               fmt.Println("Error:", value) 
         } 
   } 
 
   fmt.Println("Slice length:", len(data)) 
   if len(data)*2 > x { 
         fmt.Println("Image size (x) too small!") 
         os.Exit(-1) 
   } 
 
   maxValue := data[0] 
   for _, temp := range data { 
         if maxValue < temp { 
               maxValue = temp 
         } 
   } 
 
   if maxValue > y { 
         fmt.Println("Image size (y) too small!") 
         os.Exit(-1) 
   } 
   fmt.Println("maxValue:", maxValue) 
   barHeighPerUnit := int(y / maxValue) 
   fmt.Println("barHeighPerUnit:", barHeighPerUnit) 
```

The last part of plotData.go is the following:

```markup
   PNGfile := arguments[1] + "x" + arguments[2] + ".png" 
   OUTPUT, err := os.OpenFile(PNGfile, os.O_CREATE|os.O_WRONLY, 0644) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(-1) 
   } 
   m = image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{x, y}}) 
 
   i := 0 
   barWidth = int(x / len(data)) 
   fmt.Println("barWidth:", barWidth) 
   for _, v := range data { 
         c := getColor(v % 10) 
         yy := v * barHeighPerUnit 
         plotBar(barWidth*i, yy, c) 
         fmt.Println("plotBar", barWidth*i, yy) 
         i = i + 1 
   } 
 
   png.Encode(OUTPUT, m) 
} 
```

Although you can use plotData.go on its own, using the output of extractData.go as the input to plotData.go is as easy as executing the following command:

```markup
$ ./extractData.go 127.0.0.1 access.log{,.1} | awk '{print $2}' | ./plotData 6000 6500
Image size: 6000 6500
Slice length: 2
maxValue: 6333
barHeighPerUnit: 1
barWidth: 3000
plotBar 0 3129
plotBar 3000 6333
$ ls -l 6000x6500.png
-rw-r--r-- 1 mtsouk mtsouk 164915 Jun  5 18:25 6000x6500.png
```

The graphical output from the previous command can be an image like the one you can see in the following figure:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/ee09e9bd-e219-47d1-98f4-47de7bc75848.png)

The output generated by the plotData.go utility

Just Imagine

# Unix sockets in Go

There exist two kinds of sockets: Unix sockets and network sockets. Network sockets will be explained in _[](https://subscription.imaginedevops.io/book/programming/9781787125643/12)_[Chapter 12](https://subscription.imaginedevops.io/book/programming/9781787125643/12)_,_ _Network Programming_, whereas Unix sockets will be briefly explained in this section. However, as the presented Go functions also work with TCP/IP sockets, you will still have to wait till [](https://subscription.imaginedevops.io/book/programming/9781787125643/12)[Chapter 12](https://subscription.imaginedevops.io/book/programming/9781787125643/12), _Network Programming_, in order to fully understand them as they will not be explained here. So, this section will just present the Go code of a Unix socket client, which is a program that uses a Unix socket, which is a special Unix file, to read and write data. The name of the program will be readUNIX.go and will be presented in three parts.

The first part is the following:

```markup
package main 
 
import ( 
   "fmt" 
   "io" 
   "net" 
   "strconv" 
   "time" 
) 
```

The second part of readUNIX.go is the following:

```markup
func readSocket(r io.Reader) { 
   buf := make([]byte, 1024) 
   for { 
         n, _ := r.Read(buf[:]) 
         fmt.Print("Read: ", string(buf[0:n])) 
   } 
} 
```

The last part contains the following Go code:

```markup
func main() { 
   c, _ := net.Dial("unix", "/tmp/aSocket.sock") 
   defer c.Close() 
 
   go readSocket(c) 
   n := 0 
   for { 
         message := []byte("Hi there: " + strconv.Itoa(n) + "\n") 
         _, _ = c.Write(message) 
         time.Sleep(5 * time.Second) 
         n = n + 1 
   } 
} 
```

The use of readUNIX.go requires the presence of another process that also reads and writes to the same socket file (/tmp/aSocket.sock).

The generated output depends on the implementation of the other part: in this case, that output was the following:

```markup
$ go run readUNIX.go
Read: Hi there: 0
Read: Hi there: 1
```

If the socket file cannot be found or if no program is watching it, you will get the following error message:

```markup
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10cfe77]
    
goroutine 1 [running]:
main.main()
      /Users/mtsouk/Desktop/goCourse/ch/ch8/code/readUNIX.go:21 +0x67
exit status 2
```

Just Imagine

# RPC in Go

RPC stands for **Remote Procedure Call** and is a way of executing function calls to a remote server and getting the answer back in your clients. Once again, you will have to wait until _[](https://subscription.imaginedevops.io/book/programming/9781787125643/12)_[Chapter 12](https://subscription.imaginedevops.io/book/programming/9781787125643/12)_,_ _Network Programming_, in order to learn how to develop an RPC server and an RPC client in Go.

Just Imagine

# Programming a Unix shell in Go

This section will briefly and naively present Go code that can be used as the foundation for the development of a Unix shell. Apart from the exit command, the only other command that the program can recognize is the version command that just prints the version of the program. All other user input will be echoed on the screen.

The Go code of UNIXshell.go will be presented in three parts. However, before that I will present to you the first version of the shell, which mainly contains comments in order to better understand how I usually start the implementation of a relatively challenging program:

```markup
package main 
 
import ( 
   "fmt" 
) 
 
func main() { 
 
   // Present prompt 
 
   // Read a line 
 
   // Get the first word of the line 
 
   // If it is a built-in shell command, execute the command 
 
   // otherwise, echo the command 
 
} 
```

This is more or less the algorithm that I would use as a starting point: the good thing is that the comments briefly show how the program will operate. Keep in mind that the algorithm does not depend on the programming language. After that, it is easier to start implementing things because you know what you want to do.

So, the first part of the final version of the shell is the following:

```markup
package main 
 
import ( 
   "bufio" 
   "fmt" 
   "os" 
   "strings" 
) 
 
var VERSION string = "0.2" 
```

The second part is the following:

```markup
func main() { 
   scanner := bufio.NewScanner(os.Stdin) 
   fmt.Print("> ") 
   for scanner.Scan() { 
 
         line := scanner.Text() 
         words := strings.Split(line, " ") 
         command := words[0] 
```

Here, you just read the input from the user line by line and find out the first word of the input.

The last part of UNIXshell.go is the following:

```markup
         switch command { 
         case "exit": 
               fmt.Println("Exiting...") 
               os.Exit(0) 
         case "version": 
               fmt.Println(VERSION) 
         default: 
               fmt.Println(line) 
         } 
 
         fmt.Print("> ") 
   } 
} 
```

The aforementioned Go code checks the command that the user gave and acts accordingly.

Executing UNIXshell.go and interacting with it will generate the following output:

```markup
$ go run UNIXshell.go
> version
0.2
> ls -l
ls -l
> exit
Exiting...
```

Should you wish to learn more about creating your own Unix shell in Go, you can visit [https://github.com/elves/elvish](https://github.com/elves/elvish).

Just Imagine

# Yet another minor Go update

While I was writing this chapter, Go was updated: this is a minor update, which mainly fixes bugs:

```markup
$ date
Thu May 25 06:30:53 EEST 2017
$ go version
go version go1.8.3 darwin/amd64
```

Just Imagine

# Exercises

1.  Put the plotting functionality of plotIP.go into a Go package and use that package to rewrite both plotIP.go and plotData.go.
2.  Review the Go code of ddGo.go from _[](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_[Chapter 6](https://subscription.imaginedevops.io/book/programming/9781787125643/6)_,_ _File Input and Output_, in order to print information about its progress when receiving a SIGINFO signal.
3.  Change the Go code of cat.go to add support for multiple input files.
4.  Change the code of plotData.go in order to print gridlines to the generated image.
5.  Change the code of plotData.go in order to leave a little space between the bars of the plot.
6.  Try to make the UNIXshell.go program a little better by adding new features to it.

Just Imagine

# Summary

In this chapter, we talked about many interesting and handy topics, including signal handling and creating graphical images in Go. Additionally, we taught you how to add support for Unix pipes in your Go programs.

In the next chapter, we will talk about the most unique feature of Go, which is goroutines. You will learn what a goroutine is, how to create and synchronize them as well as how to create channels and pipelines. Have in mind that many people come to Go in order to learn a modern and safe programming language, but stay for its goroutines!