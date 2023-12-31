# Files and Directories

In the previous chapter, we talked about many important topics including developing and using Go packages, Go data structures, algorithms, and GC. However, until now, we have not developed any actual system utility. This will change very soon because starting from this really important chapter, we will begin developing real system utilities in Go by learning how to use Go, to work with the various types of files and directories of a filesystem.

You should always have in mind that Unix considers everything a file including symbolic links, directories, network devices, network sockets, entire hard drives, printers, and plain text files. The purpose of this chapter is to illustrate how the Go standard library allows us to understand if a path exists or not, as well as how to search directory structures to detect the kind of files we want. Additionally, this chapter will prove, using Go code as evidence, that many traditional Unix command-line utilities that work with files and directories do not have a difficult implementation.

In this chapter, you will learn the following topics:

-   The Go packages that will help you manipulate directories and file
-   Processing command-line arguments and options easily using the flag package
-   Developing a version of the which(1) command-line utility in Go
-   Developing a version of the pwd(1) command-line utility in Go
-   Deleting and renaming files and directories
-   Traversing directory trees easily
-   Writing a version of the find(1) utility in Go
-   Duplicating a directory structure in another place

Just Imagine

# Useful Go packages

The single most important package that allows you to manipulate files and directories as entities is the os package, which we will use extensively in this chapter. If you consider files as boxes with contents, the os package allows you to move them, put them into the wastebasket, change their names, visit them, and decide which ones you want to use, whereas the io package, which will be presented in the next chapter, allows you to manipulate the contents of a box without worrying too much about the box itself!

The flag package, which you will see in a while, lets you define and process your own flags and manipulate the command-line arguments of a Go program.

The filepath package is extremely handy as it includes the filepath.Walk() function that allows you to traverse entire directory structures in an easy way.

Just Imagine

# Command-line arguments revisited!

As we saw in [Chapter 2](https://subscription.imaginedevops.io/book/programming/9781787125643/2), _Writing Programs in Go_, you cannot work efficiently with multiple command-line arguments and options using if statements. The solution to this problem is to use the flag package, which will be explained here.

Remembering that the flag package is a standard Go package and that you do not have to search for the functionality of a flag elsewhere is extremely important.

# The flag package

The flag package does the dirty work of parsing command-line arguments and options for us; so, there is no need for writing complicated and perplexing Go code. Additionally, it supports various types of parameters, including strings, integers, and Boolean, which saves you time as you do not have to perform any data type conversions.

The usingFlag.go program illustrates the use of the flag Go package and will be presented in three parts. The first part has the following Go code:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
) 
```

The second part, which has the most important Go code of the program, is as follows:

```markup
func main() { 
   minusO := flag.Bool("o", false, "o") 
   minusC := flag.Bool("c", false, "c") 
   minusK := flag.Int("k", 0, "an int") 
 
   flag.Parse() 
```

In this part, you can see how you can define the flags that interest you. Here, you defined \-o, \-c, and \-k. Although the first two are Boolean flags, the \-k flag requires an integer value, which can be given as \-k=123.

The last part comes with the following Go code:

```markup
   fmt.Println("-o:", *minusO) 
   fmt.Println("-c:", *minusC) 
   fmt.Println("-K:", *minusK) 
 
   for index, val := range flag.Args() { 
         fmt.Println(index, ":", val) 
   } 
} 
```

In this part, you can see how you can read the value of an option, which also allows you to tell whether an option has been set or not. Additionally, flag.Args() allows you to access the unused command-line arguments of the program.

The use and the output of usingFlag.go are showcased in the following output:

```markup
$ go run usingFlag.go
-o: false
-c: false
-K: 0
$ go run usingFlag.go -o a b
-o: true
-c: false
-K: 0
0 : a
1 : b
```

However, if you forget to type the value of a command-line option (\-k) or the provided value is of the wrong type, you will get the following messages and the program will terminate:

```markup
$ ./usingFlag -k
flag needs an argument: -k
Usage of ./usingFlag:
  -c  c
  -k int
      an int
  -o  o$ ./usingFlag -k=abcinvalid value "abc" for flag -k: strconv.ParseInt: parsing "abc": invalid syntax
Usage of ./usingFlag:
  -c  c
  -k int
      an int
  -o  o
```

If you do not want your program to exit when there is a parse error, you can use the ErrorHandling type provided by the flag package, which allows you to change the way flag.Parse() behaves on errors with the help of the NewFlagSet() function. However, in systems programming, you usually want your utility to exit when there is an error in one or more command-line options.

Just Imagine

# Dealing with directories

Directories allow you to create a structure and store your files in a way that is easy for you to organize and search for them. In reality, directories are entries on a filesystem that contain lists of other files and directories. This happens with the help of **inodes**, which are data structures that hold information about files and directories.

As you can see in the following figure, directories are implemented as lists of names assigned to inodes. As a result, a directory contains an entry for itself, its parent directory, and each of its children, which among other things can be regular files or other directories:

What you should remember is that an inode holds metadata about a file, not the actual data of a file.

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/e74853d3-8d25-49c3-a968-dc7713c53a72.png)

A graphical representation of inodes

# About symbolic links

**Symbolic links** are pointers to files or directories, which are resolved at the time of access. Symbolic links, which are also called **soft links**, are not equal to the file or the directory they are pointing to and are allowed to point to nowhere, which can sometimes complicate things.

The following Go code, saved in symbLink.go and presented in two parts, allows you to check whether a path or file is a symbolic link or not. The first part is as follows:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "path/filepath" 
) 
 
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide an argument!") 
         os.Exit(1) 
   } 
   filename := arguments[1] 
```

Nothing special is happening here: you just need to make sure that you get one command-line argument in order to have something to test. The second part is the following Go code:

```markup
   fileinfo, err := os.Lstat(fil /etcename) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
 
   if fileinfo.Mode()&os.ModeSymlink != 0 { 
         fmt.Println(filename, "is a symbolic link") 
         realpath, err := filepath.EvalSymlinks(filename) 
         if err == nil { 
               fmt.Println("Path:", realpath) 
         } 
   } 
 
}
```

The aforementioned code of symbLink.go is more cryptic than usual because it uses lower-level functions. The technique for finding out whether a path is a real path or not involves the use of the os.Lstat() function that gives you information about a file or directory and the use of the Mode() function on the return value of the os.Lstat() call in order to compare the outcome with the os.ModeSymlink constant, which is the symbolic link bit.

Additionally, there exists the filepath.EvalSymlinks() function that allows you to evaluate any symbolic links that exist and return the true path of a file or directory, which is also used in symbLink.go. This might make you think that we are using lots of Go code for such a simple task, which is partially true, but when you are developing systems software, you are obliged to consider all possibilities and be cautious.

Executing symbLink.go, which only takes one command-line argument, generates the following output:

```markup
$ go run symbLink.go /etc
/etc is a symbolic link
Path: /private/etc
```

You will also see some of the aforementioned Go code as a part of bigger programs in the rest of this chapter.

# Implementing the pwd(1) command

When I start thinking about how to implement a program, so many ideas come to my mind that sometimes it becomes too difficult to decide what to do! The key here is to do something instead of waiting because as you write code, you will be able to tell whether the approach you are taking is good or not, and whether you should try another approach or not.

The pwd(1) command-line utility is pretty simplistic, yet it does a pretty good job. If you write lots of shell scripts, you should already know about pwd(1) because it is pretty handy when you want to get the full path of a file or a directory that resides in the same directory as the script that is being executed.

The Go code of pwd.go will be presented in two parts and will only support the \-P command-line option, which resolves all symbolic links and prints the physical current working directory. The first part of pwd.go is as follows:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "path/filepath" 
) 
 
func main() { 
   arguments := os.Args 
 
   pwd, err := os.Getwd() 
   if err == nil { 
         fmt.Println(pwd) 
   } else { 
         fmt.Println("Error:", err) 
   } 
```

The second part is as follows:

```markup
   if len(arguments) == 1 { 
         return 
   } 
 
   if arguments[1] != "-P" { 
         return 
   } 
 
   fileinfo, err := os.Lstat(pwd) 
   if fileinfo.Mode()&os.ModeSymlink != 0 { 
         realpath, err := filepath.EvalSymlinks(pwd) 
         if err == nil { 
               fmt.Println(realpath) 
         } 
   } 
} 
```

Note that if the current directory can be described by multiple paths, which can happen if you are using symbolic links, os.Getwd() can return any one of them. Additionally, you need to reuse some of the Go code found in symbLink.go to discover the physical current working directory in case the \-P option is given and you are dealing with a directory that is a symbolic link. Also, the reason for not using the flag package in pwd.go is that I find the code much simpler the way it is.

Executing pwd.go will generate the following output:

```markup
$ go run pwd.go
/Users/mtsouk/Desktop/goCourse/ch/ch5/code
```

On macOS machines, the /tmp directory is a symbolic link, which can help us verify that pwd.go works as expected:

```markup
$ go run pwd.go
/tmp
$ go run pwd.go -P
/tmp
/private/tmp
```

# Developing the which(1) utility in Go

The which(1) utility searches the value of the PATH environment variable in order to find out if an executable file can be found in one of the directories of the PATH variable. The following output shows the way the which(1) utility works:

```markup
$ echo $PATH
/home/mtsouk/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games
$ which ls
/home/mtsouk/bin/ls
code$ which -a ls
/home/mtsouk/bin/ls
/bin/ls
```

Our implementation of the Unix utility will support the two command-line options supported by the macOS version of which(1), which are \-a and \-s with the help of the flag package: the Linux version of which(1) does not support the \-s option. The \-a option lists all the instances of the executable instead of just the first one while the \-s returns 0 if the executable was found and 1 otherwise: this is not the same as printing 0 or 1 using the fmt package.

In order to check the return value of a Unix command-line utility in the shell, you should do the following:

```markup
$ which -s ls$ echo $?
0
```

Note that go run prints out nonzero exit codes.

The Go code for which(1) will be saved in which.go and will be presented in four parts. The first part of which.go has the following Go code:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "os" 
   "strings" 
) 
```

The strings package is needed in order to split the contents of the PATH variable after you read it. The second part of which.go deals with the use of the flag package:

```markup
func main() { 
   minusA := flag.Bool("a", false, "a") 
   minusS := flag.Bool("s", false, "s") 
 
   flag.Parse() 
   flags := flag.Args() 
   if len(flags) == 0 { 
         fmt.Println("Please provide an argument!") 
         os.Exit(1) 
   } 
   file := flags[0] 
   fountIt := false 
```

One very important part of which.go is the part that reads the PATH shell environment variable in order to split it and use it, which is presented in the third part here:

```markup
   path := os.Getenv("PATH") 
   pathSlice := strings.Split(path, ":") 
   for _, directory := range pathSlice { 
         fullPath := directory + "/" + file 
```

The last statement here constructs the full path of the file we are searching for, as if it existed in each separate directory of the PATH variable because if you have the full path of a file, you do not have to search for it!

The last part of which.go is as follows:

```markup
         fileInfo, err := os.Stat(fullPath) 
         if err == nil { 
               mode := fileInfo.Mode() 
               if mode.IsRegular() { 
                     if mode&0111 != 0 { 
                           fountIt = true 
                           if *minusS == true { 
                                 os.Exit(0) 
                           } 
                           if *minusA == true { 
                                 fmt.Println(fullPath) 
                           } else { 
                                 fmt.Println(fullPath) 
                                 os.Exit(0) 
                           } 
                     } 
               } 
         } 
   } 
   if fountIt == false { 
         os.Exit(1) 
   } 
} 
```

Here, the call to os.Stat() tells whether the file we are looking for actually exists or not. In case of success, the mode.IsRegular() function checks whether the file is a regular file or not because we are not looking for directories or symbolic links. However, we are not done yet! The which.go program performs a test to find out whether the file that was found is indeed an executable file: if it is not an executable file, it will not get printed. So, the if mode&0111 != 0 statement verifies that the file is actually an executable file using a binary operation.

Next, if the \-s flag is set to \*minusS == true, then the \-a flag does not really matter because the program will terminate as soon as it finds a match.

As you can see, there are lots of tests involved in which.go, which is not rare for systems software. Nevertheless, you should always examine all possibilities in order to avoid surprises later. The good thing is that most of these tests will be used later on in the Go implementation of the find(1) utility: it is good practice to test some features by writing small programs before putting them all together into bigger programs because by doing so, you learn the technique better and you can detect silly bugs more easily.

Executing which.go will produce the following output:

```markup
$ go run which.go ls
/home/mtsouk/bin/ls
$ go run which.go -s ls
$ echo $?
0
$ go run which.go -s ls123123
exit status 1
$ echo $?
1
$ go run which.go -a ls
/home/mtsouk/bin/ls
/bin/ls
```

# Printing the permission bits of a file or directory

With the help of the ls(1) command, you can find out the permissions of a file:

```markup
$ ls -l /bin/ls
-rwxr-xr-x  1 root  wheel  38624 Mar 23 01:57 /bin/ls
```

In this subsection, we will look at how to print the permissions of a file or directory using Go: the Go code will be saved in permissions.go and will be presented in two parts. The first part is as follows:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
) 
 
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide an argument!") 
         os.Exit(1) 
   } 
 
   file := arguments[1] 
```

The second part contains the important Go code:

```markup
   info, err := os.Stat(file) 
   if err != nil { 
         fmt.Println("Error:", err) 
         os.Exit(1) 
   } 
   mode := info.Mode() 
   fmt.Print(file, ": ", mode, "\n") 
} 
```

Once again, most of the Go code is for dealing with the command-line argument and making sure that you have one! The Go code that does the actual job is mainly the call to the os.Stat() function, which returns a FileInfo structure that describes the file or directory examined by os.Stat(). From the FileInfo structure, you can discover the permissions of a file by calling the Mode() function.

Executing permissions.go produces the following output:

```markup
$ go run permissions.go /bin/ls
/bin/ls: -rwxr-xr-x
$ go run permissions.go /usr
/usr: drwxr-xr-x
$ go run permissions.go /us
Error: stat /us: no such file or directory
exit status 1
```

Just Imagine

# Dealing with files in Go

An extremely important task of an operating system is working with files because all data is stored in files. In this section, we will show you how to delete and rename files, and in the next section, _Developing find(1) in Go_, we will teach you how to search directory structures in order to find the files you want.

# Deleting a file

In this section, we will illustrate how to delete files and directories using the os.Remove() Go function.

When testing programs that delete files and directories be extra careful and use common sense!

The rm.go file is a Go implementation of the rm(1) tool that illustrates how you can delete files in Go. Although the core functionality of rm(1) is there, the options of rm(1) are missing: it would be a good exercise to try to implement some of them. Just pay extra attention when implementing the \-f and \-R options.

The Go code of rm.go is as follows:

```markup
package main 
import ( 
   "fmt" 
   "os" 
) 
 
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Please provide an argument!") 
         os.Exit(1) 
   } 
 
   file := arguments[1] 
   err := os.Remove(file) 
   if err != nil { 
         fmt.Println(err) 
         return 
   } 
} 
```

If rm.go is executed without any problems, it will create no output according to the Unix philosophy. So, what is interesting here is watching the error messages you can get when the file you are trying to delete does not exist: both when you do not have the necessary permissions to delete it and when a directory is not empty:

```markup
$ go run rm.go 123
remove 123: no such file or directory
$ ls -l /tmp/AlTest1.err-rw-r--r--  1 root  wheel  1278 Apr 17 20:13 /tmp/AlTest1.err
$ go run rm.go /tmp/AlTest1.err
remove /tmp/AlTest1.err: permission denied
$ go run rm.go test
remove test: directory not empty
```

# Renaming and moving files

In this subsection, we will show you how to rename and move a file using Go code: the Go code will be saved as rename.go. Although the same code can be used for renaming or moving directories, rename.go is only allowed to work with files.

When performing things that cannot be easily undone, such as overwriting a file, you should be extra careful and maybe inform the user that the destination file already exists in order to avoid unpleasant surprises. Although the default operation of the traditional mv(1) utility will automatically overwrite the destination file if it exists, I do not think that this is very safe. Therefore, rename.go will not overwrite destination files by default.

When developing systems software, you have to deal with all the details or the details will reveal themselves as bugs when least expected! Extensive testing will allow you to find the details you missed and correct them.

The code of rename.go will be presented in four parts. The first part includes the expected preamble as well as the Go code for dealing with the setup of the flag package:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "os" 
   "path/filepath" 
) 
 
func main() { 
   minusOverwrite := flag.Bool("overwrite", false, "overwrite") 
 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) < 2 { 
         fmt.Println("Please provide two arguments!") 
         os.Exit(1) 
   } 
```

The second part has the following Go code:

```markup
   source := flags[0] 
   destination := flags[1] 
   fileInfo, err := os.Stat(source) 
   if err == nil { 
         mode := fileInfo.Mode() 
         if mode.IsRegular() == false { 
               fmt.Println("Sorry, we only support regular files as source!") 
               os.Exit(1) 
         } 
   } else { 
         fmt.Println("Error reading:", source) 
         os.Exit(1) 
   } 
```

This part makes sure the source file exists, is a regular file, and is not a directory or something else like a network socket or a pipe. Once again, the trick with os.Stat() you saw in which.go is used here.

The third part of rename.go is as follows:

```markup
   newDestination := destination 
   destInfo, err := os.Stat(destination) 
   if err == nil { 
         mode := destInfo.Mode() 
         if mode.IsDir() { 
               justTheName := filepath.Base(source) 
               newDestination = destination + "/" + justTheName 
         } 
   } 
```

There is another tricky point here; you will need to consider the case where the source is a plain file and the destination is a directory, which is implemented with the help of the newDestination variable.

Another special case that you should consider is when the source file is given in a format that contains an absolute or relative path in it like ./aDir/aFile. In this case, when the destination is a directory, you should get the basename of the path, which is what follows the last / character and in this case is aFile, and add it to the destination directory in order to correctly construct the newDestination variable. This happens with the help of the filepath.Base() function, which returns the last element of a path.

Finally, the last part of rename.go has the following Go code:

```markup
   destination = newDestination 
   destInfo, err = os.Stat(destination) 
   if err == nil { 
         if *minusOverwrite == false { 
               fmt.Println("Destination file already exists!") 
               os.Exit(1) 
         } 
   } 
 
   err = os.Rename(source, destination) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

The most important Go code of rename.go has to do with recognizing whether the destination file exists or not. Once again, this is implemented with the support of the os.Stat() function. If os.Stat() returns an error message, this means that the destination file does not exist; so, you are free to call os.Rename(). If os.Stat() returns nil, this means that the os.Stat() call was successful and that the destination file exists. In this case, you should check the value of the overwrite flag to see if you are allowed to overwrite the destination file or not.

When everything is OK, you are free to call os.Rename() and perform the desired task!

If rename.go is executed correctly, it will create no output. However, if there are problems, rename.go will generate some output:

```markup
$ touch newFILE
$ ./rename newFILE regExpFind.go
Destination file already exists!
$ ./rename -overwrite newFILE regExpFind.go
$
```

Just Imagine

# Developing find(1) in Go

This section will teach you the necessary things that you need to know in order to develop a simplified version of the find(1) command-line utility in Go. The developed version will not support all the command-line options supported by find(1), but it will have enough options to be truly useful.

What you will see in the following subsections is the entire process in small steps. So, the first subsection will show you the Go way for visiting all files and directories in a given directory tree.

# Traversing a directory tree

The most important task that find(1) needs to support is being able to visit all files and sub directories starting from a given directory. So, this section will implement this task in Go. The Go code of traverse.go will be presented in three parts. The first part is the expected preamble:

```markup
package main 
 
import ( 
   "fmt" 
   "os" 
   "path/filepath" 
) 
```

The second part is about implementing a function named walkFunction() that will be used as an argument to a Go function named filepath.Walk():

```markup
func walkFunction(path string, info os.FileInfo, err error) error { 
   _, err = os.Stat(path) 
   if err != nil { 
         return err 
   } 
 
   fmt.Println(path) 
   return nil 
} 
```

Once again, the os.Stat() function is used because a successful os.Stat() function call means that we are dealing with something (file, directory, pipe, and so on) that actually exists!

Do not forget that between the time filepath.Walk() is called and the time walkFunction() is called and executed, many things can happen in an active and busy filesystem, which is the main reason for calling os.Stat().

The last part of the code is as follows:

```markup
func main() { 
   arguments := os.Args 
   if len(arguments) == 1 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(1) 
   } 
 
   Path := arguments[1] 
   err := filepath.Walk(Path, walkFunction) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

All the dirty jobs here are automatically done by the filepath.Walk() function with the help of the walkFunction() function that was defined previously. The filepath.Walk() function takes two parameters: the path of a directory and the walk function it will use.

Executing traverse.go will generate the following kind of output:

```markup
$ go run traverse.go ~/code/C/cUNL
/home/mtsouk/code/C/cUNL
/home/mtsouk/code/C/cUNL/gpp
/home/mtsouk/code/C/cUNL/gpp.c
/home/mtsouk/code/C/cUNL/sizeofint
/home/mtsouk/code/C/cUNL/sizeofint.c
/home/mtsouk/code/C/cUNL/speed
/home/mtsouk/code/C/cUNL/speed.c
/home/mtsouk/code/C/cUNL/swap
/home/mtsouk/code/C/cUNL/swap.c
```

As you can see, the code of traverse.go is pretty naive, as among other things, it cannot differentiate between directories, files, and symbolic links. However, it does the pretty tedious job of visiting every file and directory under a given directory tree, which is the basic functionality of the find(1) utility.

# Visiting directories only!

Although it is good to be able to visit everything, there are times when you want to visit only directories and not files. So, in this subsection, we will modify traverse.go in order to still visit everything but only print the directory names. The name of the new program will be traverseDir.go. The only part of traverse.go that needs to change is the definition of the walkFunction():

```markup
func walkFunction(path string, info os.FileInfo, err error) error { 
   fileInfo, err := os.Stat(path) 
   if err != nil { 
         return err 
   } 
 
   mode := fileInfo.Mode() 
   if mode.IsDir() { 
         fmt.Println(path) 
   } 
   return nil 
} 
```

As you can see, here you need to use the information returned by the os.Stat() function call in order to check whether you are dealing with a directory or not. If you have a directory, then you print its path and you are done.

Executing traverseDir.go will generate the following output:

```markup
$ go run traverseDir.go ~/code
/home/mtsouk/code
/home/mtsouk/code/C
/home/mtsouk/code/C/cUNL
/home/mtsouk/code/C/example
/home/mtsouk/code/C/sysProg
/home/mtsouk/code/C/system
/home/mtsouk/code/Haskell
/home/mtsouk/code/aLink
/home/mtsouk/code/perl
/home/mtsouk/code/python  
```

Just Imagine

# The first version of find(1)

The Go code in this section is saved as find.go and will be presented in three parts. As you will see, find.go uses a large amount of the code found in traverse.go, which is the main benefit you get when you are developing a program step by step.

The first part of find.go is the expected preamble:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "os" 
   "path/filepath" 
) 
```

As we already know that we will improve find.go in the near future, the flag package is used here even if this is the first version of find.go and it does not have any flags!

The second part of the Go code contains the implementation of the walkFunction():

```markup
func walkFunction(path string, info os.FileInfo, err error) error { 
 
   fileInfo, err := os.Stat(path) 
   if err != nil { 
         return err 
   } 
 
   mode := fileInfo.Mode() 
   if mode.IsDir() || mode.IsRegular() { 
         fmt.Println(path) 
   } 
   return nil 
} 
```

From the implementation of the walkFunction() you can easily understand that find.go only prints regular files and directories, and nothing else. Is this a problem? Not, if this is what you want. Generally speaking, this is not good. Nevertheless, having a first version of something that works despite some restrictions is a good starting point! The next version, which will be named improvedFind.go, will improve find.go by adding various command-line options to it.

The last part of find.go contains the code that implements the main() function:

```markup
func main() { 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) == 0 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(1) 
   } 
 
   Path := flags[0] 
   err := filepath.Walk(Path, walkFunction) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

Executing find.go will create the following output:

```markup
$ go run find.go ~/code/C/cUNL
/home/mtsouk/code/C/cUNL
/home/mtsouk/code/C/cUNL/gpp
/home/mtsouk/code/C/cUNL/gpp.c
/home/mtsouk/code/C/cUNL/sizeofint
/home/mtsouk/code/C/cUNL/sizeofint.c
/home/mtsouk/code/C/cUNL/speed
/home/mtsouk/code/C/cUNL/speed.c
/home/mtsouk/code/C/cUNL/swap
/home/mtsouk/code/C/cUNL/swap.c
```

# Adding some command-line options

This subsection will try to improve the Go version of find(1) that you created earlier. Keep in mind that this is the process used for developing real programs because you do not implement every possible command-line option in the first version of a program.

The Go code of the new version is going to be saved as improvedFind.go. Among other things, the new version will be able to ignore symbolic links: symbolic links will only be printed when improvedFind.go is used with the appropriate command-line option. To do this, we will use some of the Go code of symbLink.go.

The improvedFind.go program is a real system tool that you can use on your own Unix machines.

The supported flags will be the following:

-   **\-s**: This is for printing socket files
-   **\-p**: This is for printing pipes
-   **\-sl**: This is for printing symbolic links
-   **\-d**: This is for printing directories
-   **\-f**: This is for printing files

As you will see, most of the new Go code is for supporting the flags added to the program. Additionally, by default, improvedFind.go prints every type of file or directory, and you are allowed to combine any of the preceding flags in order to print the types of files you want.

Apart from the various changes in the implementation of the main() function in order to support all these flags, most of the remaining changes will take place in the code of the walkFunction() function. Additionally, the walkFunction() function will be defined inside the main() function, which happens in order to avoid the use of global variables.

The first part of improvedFind.go is as follows:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "os" 
   "path/filepath" 
) 
 
func main() { 
 
   minusS := flag.Bool("s", false, "Sockets") 
   minusP := flag.Bool("p", false, "Pipes") 
   minusSL := flag.Bool("sl", false, "Symbolic Links") 
   minusD := flag.Bool("d", false, "Directories") 
   minusF := flag.Bool("f", false, "Files") 
 
   flag.Parse() 
   flags := flag.Args() 
 
   printAll := false 
   if *minusS && *minusP && *minusSL && *minusD && *minusF { 
         printAll = true 
   } 
 
   if !(*minusS || *minusP || *minusSL || *minusD || *minusF) { 
         printAll = true 
   } 
 
   if len(flags) == 0 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(1) 
   } 
 
   Path := flags[0] 
```

So, if all the flags are unset, the program will print everything, which is handled by the first if statement. Similarly, if all the flags are set, the program will also print everything. So, a new Boolean variable named printAll is needed.

The second part of improvedFind.go has the following Go code, which is mainly the definition of the walkFunction variable, which in reality is a function:

```markup
   walkFunction := func(path string, info os.FileInfo, err error) error { 
         fileInfo, err := os.Stat(path) 
         if err != nil { 
               return err 
         } 
 
         if printAll { 
               fmt.Println(path) 
               return nil 
         } 
 
         mode := fileInfo.Mode() 
         if mode.IsRegular() && *minusF { 
               fmt.Println(path) 
               return nil 
         } 
 
         if mode.IsDir() && *minusD { 
               fmt.Println(path) 
               return nil 
         } 
 
         fileInfo, _ = os.Lstat(path) 
         if fileInfo.Mode()&os.ModeSymlink != 0 { 
               if *minusSL { 
                     fmt.Println(path) 
                     return nil 
               } 
         } 
 
         if fileInfo.Mode()&os.ModeNamedPipe != 0 { 
               if *minusP { 
                     fmt.Println(path) 
                     return nil 
               } 
         } 
 
         if fileInfo.Mode()&os.ModeSocket != 0 { 
               if *minusS { 
                     fmt.Println(path) 
                     return nil 
               } 
         } 
 
         return nil 
   } 
```

Here, the good thing is that once you find a match and print a file, you do not have to visit the rest of the if statements, which is the main reason for putting the minusF check first and the minusD check second. The call to os.Lstat() is used to find out whether we are dealing with a symbolic link or not. This happens because os.Stat() follows symbolic links and returns information about the file the link references, whereas os.Lstat() does not do so: the same occurs with stat(2) and lstat(2).

You should be pretty familiar with the last part of improvedFind.go:

```markup
   err := filepath.Walk(Path, walkFunction) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

Executing improvedFind.go generates the following output, which is an enriched version of the output of find.go:

```markup
$ go run improvedFind.go -d ~/code/C
/home/mtsouk/code/C
/home/mtsouk/code/C/cUNL
/home/mtsouk/code/C/example
/home/mtsouk/code/C/sysProg
/home/mtsouk/code/C/system
$ go run improvedFind.go -sl ~/code
/home/mtsouk/code/aLink
```

# Excluding filenames from the find output

There are times when you do not need to display everything from the output of find(1). So, in this subsection, you will learn a technique that allows you to manually exclude files from the output of improvedFind.go based on their filenames.

Note that this version of the program will not support regular expressions and will only exclude filenames that are an exact match.

So, the improved version of improvedFind.go will be named excludeFind.go. The output of the diff(1) utility can reveal the code differences between improvedFind.go and excludeFind.go:

```markup
$ diff excludeFind.go improvedFind.go
10,19d9
< func excludeNames(name string, exclude string) bool {`
<     if exclude == "" {
<           return false
<     }
<     if filepath.Base(name) == exclude {
<           return true
<     }
<     return false
< }
<
27d16
<     minusX := flag.String("x", "", "Files")
54,57d42
<           if excludeNames(path, *minusX) {
<                 return nil
<           }
<
```

The most significant change is the introduction of a new Go function, named excludeNames(), that deals with filename exclusion and the addition of the \-x flag, which is used for setting the filename you want to exclude from the output. All the job is done by the file path. The Base() function finds the last part of a path, even if the path is not a file but a directory, and compares it against the value of the \-x flag.

Note that a more appropriate name for the excludeNames() function could have been isExcluded() or something similar because the \-x option accepts a single value.

Executing excludeFind.go with and without the \-x flag will prove that the new Go code actually works:

```markup
$ go run excludeFind.go -x=dT.py ~/code/python
/home/mtsouk/code/python
/home/mtsouk/code/python/dataFile.txt
/home/mtsouk/code/python/python
$ go run excludeFind.go ~/code/python
/home/mtsouk/code/python
/home/mtsouk/code/python/dT.py
/home/mtsouk/code/python/dataFile.txt
/home/mtsouk/code/python/python
```

# Excluding a file extension from the find output

A file extension is the part of a filename after the last dot (.) character. So, the file extension of the image.png file is png, which applies to both files and directories.

Once again, you will need a separate command-line option followed by the file extension you want to exclude in order to implement this functionality: the new flag will be named \-ext. This version of the find(1) utility will be based on the code of excludeFind.go and will be named finalFind.go. Some of you might say that a more appropriate name for this option would have been \-xext and you would be right about that!

Once again, the diff(1) utility will help us spot the code differences between excludeFind.go and finalFind.go: the new functionality is implemented in a Go function named excludeExtensions(), which makes things easier to understand:

```markup
$ diff finalFind.go excludeFind.go
8d7
<     "strings"
21,34d19
< func excludeExtensions(name string, extension string) bool {
<     if extension == "" {
<           return false
<     }
<     basename := filepath.Base(name)
<     s := strings.Split(basename, ".")
<     length := len(s)
<     basenameExtension := s[length-1]
<     if basenameExtension == extension {
<           return true
<     }
<     return false
< }
<
43d27
<     minusEXT := flag.String("ext", "", "Extensions")
74,77d57
<           if excludeExtensions(path, *minusEXT) {
<                 return nil
<           }
< 
```

As we are looking for the string after the last dot in the path, we use strings.Split() to split the path based on the dot characters it contains. Then, we take the last part of the return value of strings.Split() and we compare it against the extension that was given with the \-ext flag. Therefore, nothing special here, just some string manipulation code. Once again, a more appropriate name for excludeExtensions() would have been isExcludedExtension().

Executing finalFind.go will generate the following output:

```markup
$ go run finalFind.go -ext=py ~/code/python
/home/mtsouk/code/python
/home/mtsouk/code/python/dataFile.txt
/home/mtsouk/code/python/python
$ go run finalFind.go ~/code/python
/home/mtsouk/code/python
/home/mtsouk/code/python/dT.py
/home/mtsouk/code/python/dataFile.txt
/home/mtsouk/code/python/python
```

Just Imagine

# Using regular expressions

This section will illustrate how to add support for regular expressions in finalFind.go: the name of the last version of the tool will be regExpFind.go. The new flag will be called \-re and it will require a string value: anything that matches this string value will be included in the output unless it is excluded by another command-line option. Additionally, due to the flexibility that flags offer, we do not need to delete any of the previous options in order to add another one!

Once again, the diff(1) command will tell us the code differences between regExpFind.go and finalFind.go:

```markup
$ diff regExpFind.go finalFind.go
8d7
<     "regexp"
36,44d34
< func regularExpression(path, regExp string) bool {
<     if regExp == "" {
<           return true
<     }
<     r, _ := regexp.Compile(regExp)
<     matched := r.MatchString(path)
<     return matched
< }
<
54d43
<     minusRE := flag.String("re", "", "Regular Expression")
71a61
>
75,78d64
<           if regularExpression(path, *minusRE) == false {
<                 return nil
<           }
< 
```

In [Chapter 7](https://subscription.imaginedevops.io/book/programming/9781787125643/7)_,_ _Working with System Files_, we ;will talk more about pattern matching and regular expressions in Go: for now, it is enough to understand that regexp.Compile() creates a regular expression and MatchString() tries to do the matching in the regularExpression() function.

Executing regExpFind.go will generate the following output:

```markup
$ go run regExpFind.go -re=anotherPackage /Users/mtsouk/go
/Users/mtsouk/go/pkg/darwin_amd64/anotherPackage.a
/Users/mtsouk/go/src/anotherPackage
/Users/mtsouk/go/src/anotherPackage/anotherPackage.go
$ go run regExpFind.go -ext=go -re=anotherPackage /Users/mtsouk/go
/Users/mtsouk/go/pkg/darwin_amd64/anotherPackage.a
/Users/mtsouk/go/src/anotherPackage 
```

The previous output can be verified by using the following command:

```markup
$ go run regExpFind.go /Users/mtsouk/go | grep anotherPackage
/Users/mtsouk/go/pkg/darwin_amd64/anotherPackage.a
/Users/mtsouk/go/src/anotherPackage
/Users/mtsouk/go/src/anotherPackage/anotherPackage.go
```

# Creating a copy of a directory structure

Armed with the knowledge you gained in the previous sections, we will now develop a Go program that creates a copy of a directory structure in another directory: this means that any files in the input directory will not be copied to the destination directory, only the directories will be copied. This can be handy when you want to save useful files from a directory structure somewhere else while keeping the same directory structure or when you want to take a backup of a filesystem manually.

As you are only interested in directories, the code of cpStructure.go is based on the code of traverseDir.go you saw earlier in this chapter: once again, a small program that was developed for learning purposes helps you implement a bigger program! Additionally, the test option will show what the program will do without actually creating any directories.

The code of cpStructure.go will be presented in four parts. The first one is as follows:

```markup
package main 
 
import ( 
   "flag" 
   "fmt" 
   "os" 
   "path/filepath" 
   "strings" 
) 
```

There is nothing special here, just the expected preamble. The second part is as follows:

```markup
func main() { 
   minusTEST := flag.Bool("test", false, "Test run!") 
 
   flag.Parse() 
   flags := flag.Args() 
 
   if len(flags) == 0 || len(flags) == 1 { 
         fmt.Println("Not enough arguments!") 
         os.Exit(1) 
   } 
 
   Path := flags[0] 
   NewPath := flags[1] 
 
   permissions := os.ModePerm 
   _, err := os.Stat(NewPath) 
   if os.IsNotExist(err) { 
         os.MkdirAll(NewPath, permissions) 
   } else { 
         fmt.Println(NewPath, "already exists - quitting...") 
         os.Exit(1) 
   } 
```

The cpStructure.go program demands that the destination directory does not exist in advance in order to avoid unnecessary surprises and errors afterwards.

The third part contains the code of the walkFunction variable:

```markup
   walkFunction := func(currentPath string, info os.FileInfo, err error) error { 
         fileInfo, _ := os.Lstat(currentPath) 
         if fileInfo.Mode()&os.ModeSymlink != 0 { 
               fmt.Println("Skipping", currentPath) 
               return nil 
         } 
 
         fileInfo, err = os.Stat(currentPath) 
         if err != nil { 
               fmt.Println("*", err) 
               return err 
         } 
 
         mode := fileInfo.Mode() 
         if mode.IsDir() { 
               tempPath := strings.Replace(currentPath, Path, "", 1) 
               pathToCreate := NewPath + "/" + filepath.Base(Path) + tempPath 
 
               if *minusTEST { 
                     fmt.Println(":", pathToCreate) 
                     return nil 
               } 
 
               _, err := os.Stat(pathToCreate) 
               if os.IsNotExist(err) { 
                     os.MkdirAll(pathToCreate, permissions) 
               } else { 
                     fmt.Println("Did not create", pathToCreate, ":", err) 
               } 
         } 
         return nil 
   } 
```

Here, the first if statement makes sure that we will deal with symbolic links because symbolic links can be dangerous and create problems: always try to treat special situations in order to avoid problems and nasty bugs.

The os.IsNotExist() function allows you to make sure that the directory you are trying to create is not already there. So, if the directory is not there, you create it using ;os.MkdirAll(). The os.MkdirAll() function creates a directory path including all the necessary parents, which makes things simpler for the developer.

Nevertheless, the trickiest part that the code of the walkFunction variable has to deal with is removing the unnecessary parts of the source path and constructing the new path correctly. The strings.Replace() function used in the program replaces the occurrences of its second argument (Path) that can be found in the first argument (currentPath) with its third argument ("") as many times as its last argument (1). If the last argument is a negative number, which is not the case here, then there will be no limit to the number of replacements. In this case, it removes the value of the Path variable, which is the source directory, from the currentPath variable, which is the directory that is being examined.

The last part of the program is as follows:

```markup
   err = filepath.Walk(Path, walkFunction) 
   if err != nil { 
         fmt.Println(err) 
         os.Exit(1) 
   } 
} 
```

Executing cpStructure.go will generate the following output:

```markup
$ go run cpStructure.go ~/code /tmp/newCode
Skipping /home/mtsouk/code/aLink
$ ls -l /home/mtsouk/code/aLink
lrwxrwxrwx 1 mtsouk mtsouk 14 Apr 21 18:10 /home/mtsouk/code/aLink -> /usr/local/bin 
```

The following figure shows a graphical representation of the source and destination directory structures used in the aforementioned example:

![](https://static.packt-cdn.com/products/9781787125643/graphics/assets/ecf2e299-f496-48f4-9622-d1225bd52ad4.png)

A graphical representation of two directory structures with their files

Just Imagine

# Exercises

1.  Read the documentation page of the os package at [https://golang.org/pkg/os/](https://golang.org/pkg/os/).
2.  Visit [https://golang.org/pkg/path/filepath/](https://golang.org/pkg/path/filepath/) to learn more about the filepath.Walk() function.
3.  Change the code of rm.go in order to support multiple command-line arguments, and then try to implement the \-v command-line option of the rm(1) utility.
4.  Make the necessary changes to the Go code of which.go in order to support multiple command-line arguments.
5.  Start implementing a version of the ls(1) utility in Go. Do not try to support every ls(1) option at once.
6.  Change the code of traverseDir.go in order to print regular files only.
7.  Check the manual page of find(1) and try to add support for some of its options in regExpFind.go.

Just Imagine

# Summary

In this chapter, we talked about many things including the use of the flag standard package, Go functions that allow you to work with directories and files, and traverse directory structures, and we developed Go versions of various Unix command-line utilities including pwd(1), which(1), rm(1), and find(1).

In the next chapter, we will continue talking about file operations, but this time you will learn how to read files and write to files in Go: as you will see there are many ways to do this. Although this gives you versatility, it also demands that you should be able to choose the right technique to do your job as efficiently as possible! So, you will start by learning more about the io package as well as the bufio package and by the end of the chapter, you will have Go versions of the wc(1) and dd(1) utilities!