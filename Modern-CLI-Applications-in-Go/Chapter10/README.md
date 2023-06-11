# The Empathic Side of Development

Empathy has been a hot topic lately, and its relation to software is no exception. This chapter will discuss how to use empathy to develop a better CLI. Empathy-driven CLI development is done with consideration of the output and errors that are written and the clarity and reassurance it may give the user. Written documentation that takes an empathetic approach also provides users with an effortless way to get started, while help and support are readily available for users when they need it.

This chapter will give examples of how to rewrite errors in a way that users may easily understand, not just by being clearer that an error occurred but also how and where (with debug and traceback information), which can be provided with a `--verbose` flag and detailed logging. It is very important to provide logs for users, and this implementation will be described when discussing debug and traceback information. Users can also feel more reassured with the help of man pages, usage examples of each command, empathically written documentation, and a quick and easy way to submit bugs that are encountered within the application.

Taking an empathetic approach into many different areas of your application, as well as in your life, is a form of not only self-care but care for others as well. Hopefully, these tips will help to create a CLI that meets the user at their perspective and provides them with a feeling of reassurance. Specifically, this chapter will cover the following topics:

-   Rewriting errors to be human-readable
-   Providing debug and traceback information
-   Effortless bug submission
-   Help, documentation, and support

# Rewriting errors to be human-readable

Errors can be a big point of frustration for users as they can set users off their original plans. Users will be grateful, though, if you can make the process as painless as possible. In this section, we will discuss some ways to ease users when an error occurs and provide some guidelines for creating better error messages and avoiding some common mistakes. Creating clear and helpful error messages is often overlooked, yet they are very impactful toward an optimal UX.

Think of some of your subjective experiences while working with CLIs and some of the errors you have encountered. This is an opportunity to think about how experiences can be improved for yourself when working with your own CLI, but also for others.

## Guidelines for writing error messages

Here are some useful guidelines when writing error messages:

-   **Be specific**: Customize messages toward the actual task that has occurred. This error message is critical if the task required inputting credentials or a final command to complete a workflow. The best experience would include specifying the exact problem and then providing a way toward correcting the issue. Specific guidance helps the users stay engaged and willing to make corrections.
-   **Remind the user that there’s a human on the other end**: A generic error message can sound very technical and intimidating to most users. By rewriting the error message, you can make them more useful and less intimidating. Empathize with your users and make sure that you don’t place blame on the user, which can be particularly discouraging. It’s important to encourage the user by being understanding, friendly, and speaking the same language, both literally and figuratively! How do the words you use sound in conversation?
-   **Keep it light-hearted**: Keeping a light-hearted tone can help ease any tension when an error occurs, but be careful! In certain situations, it might make the situation a bit worse—especially if it’s a critical task. Users do not want to feel mocked. Regardless, with humor or not, the error message should still be informational, clear, and polite.
-   **Make it easy**: This will require you to do a bit more of the heavy lifting, but it will certainly be worth it in the end. Provide clear next steps, or commands to run, to resolve the issue and to help the user get back on track to what they had originally intended on doing. With helpful suggestions, the user will at least see the path through the trees and know what to do next.
-   **Consider the best placement**: When outputting error messages, it’s best to place them in an area where users will look first. In the case of the CLI, it’s most likely at the end of the output.
-   **Consolidate errors**: If there are multiple error messages, especially similar ones, group them together. It will look much better than repeating the same error message over and again.
-   **Optimize your error message with icons and text**: Usually, important information is placed at the end of the output, but if there’s any red text on the screen, that is often where the user’s eyes will be drawn to. Given the power of color, use it sparingly and with intention.
-   **Consider capitalization and punctuation**: Don’t write in all caps or with multiple exclamation points. Consider consistency as well—do your errors start with capitalization? If they are output to a log, errors may start all in lowercase letters.

## Decorating errors

Wrapping errors with additional information and context is a very important step. What is the specific task that failed and why? This helps the user understand what happened. Providing actions to take toward resolution will also help the user feel more supported and willing to move forward.

First, there are a few ways to decorate your errors with additional information. You can use the `fmt.Errorf` function:

```markup
func Errorf(format string, a ...interface{}) error
```

With this function, you can print out the error as a string with any additional context. Here’s an example within the `errors/errors.go` file in the `Chapter-9` repo:

```markup
birthYear := -1981
err := fmt.Errorf("%d is negative\nYear can't be negative", birthYear)
if birthYear < 0 {
    fmt.Println(err)
} else {
    fmt.Printf("Birth year: %d\n", birthYear)
}
```

The next way to decorate your errors is by using the `errors.Wrap` method. This method is fully defined as follows:

```markup
func Wrap(err error, message string) error
```

It returns an error annotating `err` with a message and a stack trace at the point the method is called. If `err` is `nil`, then the `Wrap` function also returns `nil`.

In the `wrapping()` function, we demonstrate this:

```markup
func wrapping() error {
    err := errors.New("error")
    err1 := operation1()
    if err1 != nil {
        err1 = errors.Wrap(err, "operation1")
    }
    err2 := operation2()
    if err != nil {
        err2 = errors.Wrap(err1, "operation2")
    }
    err3 := operation3()
    if err != nil {
        err3 = errors.Wrap(err2, "operation3")
    }
    return err3
}
```

Notice that the previous error gets wrapped into the next error and so on until the final error is returned. The output of the error returned from the `wrapping()` function is shown here. I’ve removed the longer path for clarity:

```markup
error
.../errors.wrapping
        .../errors/errors.go:73
.../errors.Examples
        .../errors/errors.go:39
main.main
        .../main.go:6
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1594
operation1
.../errors.wrapping
        .../errors/errors.go:76
.../errors.Examples
        .../errors/errors.go:39
main.main
        .../main.go:6
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1594
operation2
.../errors.wrapping
        .../errors/errors.go:80
.../errors.Examples
        .../errors/errors.go:39
main.main
        .../main.go:6
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1594
operation3
.../errors.wrapping
        .../errors/errors.go:84
.../errors.Examples
        .../errors/errors.go:39
main.main
        .../main.go:6
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1594
```

Notice that the errors from `operation1`, `operation2`, and `operation3` are wrapped under the original `error` instance.

Because wrapping annotates the error with the stack trace and message, the line that calls the `wrapping()` function prints the error message followed by the stack trace at the call of the `New()` or `Wrap()` method.

## Customizing errors

Creating custom errors allows you to store whatever information you think is valuable to your users with the error so that when it’s time to print out, all the information is available within a single struct. First, you need to think about the error structure:

```markup
type error interface {
    Error() string
}
```

Simply create any type that implements the `Error() string` method. Think about the data you’d want stored on the custom error structure that might be useful for your users, or even for yourself as the developer, for debugging purposes. This could include the method name where the error occurred, the severity of the error, or the kind of error. In the `Chapter-9` repo, in the `errors.go` file, I provide some examples. To keep things simple, only one additional field, `Task`, is added to the `customError` structure:

```markup
type customError struct {
    Task string
    Err error
}
```

The `Error()` method that satisfies the previous interface is defined here. For fun, we use the `github.com/fatih/color` color page used in the previous chapter and an emoji (a red cross mark) alongside the error message:

```markup
func (e *customError) Error() string {
    var errorColor = color.New(color.BgRed,
        color.FgWhite).SprintFunc()
    return fmt.Sprintf("%s: %s %s", errorColor(e.Task),
        crossMark, e.Err)
}
```

Now, we can demonstrate how this custom error can be used within the `eligibleToVote` function:

```markup
func eligibleToVote(age int) error {
    fmt.Printf("%s Attempting to vote at %d years
        old...\n", votingBallot, age)
    minimumAge := 18
    err := &customError{
        Task: " eligibleToVote",
    }
    if age < minimumAge && age > 0 {
        years := minimumAge - age
        err.Err = fmt.Errorf("too young to vote, at %d,
            wait %d more years", age, years)
        return err
    }
    if age < 0 {
        err.Err = fmt.Errorf("age cannot be negative: %d",
            age)
        return err
    }
    fmt.Println("Voted.", checkMark)
    return nil
}
```

Notice there are multiple errors, and the error is initially defined at the top of the function, setting only the `Task` field. For each error that occurs, the `Err` field is then set and returned. Within the `Examples` method, we call the function with the following lines:

```markup
birthYear = 2010
currentYear := 2022
age := currentYear - birthYear
err = eligibleToVote(age)
if err != nil {
    fmt.Println("error occurred: ", err)
}
```

The following error is output when the preceding code runs:

![Figure 9.1 – Screenshot of voting error](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_9.1._B18883.jpg)

Figure 9.1 – Screenshot of voting error

There are plenty of other ways to create custom errors, but here are a few things to consider adding to your custom errors:

-   The severity of the error for logging purposes
-   Any data that may be valuable for metrics
-   The kind of error so that you may easily filter out any unexpected errors when they occur

## Writing better error messages

Now that we know how to add more detail to error messages, let’s revisit the `audiofile` CLI and rewrite our error messages to be more human-friendly using the guidelines mentioned earlier in this section. In the repo, for this particular branch, I’ve decorated the errors with extra information so that the user or developer can better understand where the error occurred and why.

Since the `audiofile` CLI interacts with the `audiofile` API, there are HTTP responses that can be handled and rewritten. A `CheckResponse` function exists in the `utils/http.go` file and does this:

```markup
func CheckResponse(resp *http.Response) error {
    if resp != nil {
        if resp.StatusCode != http.StatusOK {
            switch resp.StatusCode {
            case http.StatusInternalServerError:
                return fmt.Errorf(errorColor("retry the command 
                  later"))
            case http.StatusNotFound:
                return fmt.Errorf(errorColor("the id cannot be 
                  found"))
            default:
                return fmt.Errorf(errorColor(fmt.
                  Sprintf("unexpected response: %v", resp.
                  Status)))
            }
        }
        return nil
    } else {
        return fmt.Errorf(errorColor("response body is nil"))
    }
}
```

Consider expanding on this within your own CLI, which might also interact with a REST API. You may check as many responses as you like and rewrite them as errors to be returned by the command.

In previous versions of the `audiofile` CLI, if an `id` parameter was passed into the `get` or `delete` command, nothing would be returned if the ID was not found. However by passing back the `http.StatusNotFound` response and adding additional error decorations, the command that would previously error silently and return no data can now return some useful information:

```markup
mmontagnino@Marians-MacCourse-Pro audiofile % ./bin/audiofile get --id 1234
Sending request: GET http://localhost:8000/request?id=1234 ...
Error:
  checking response: the id cannot be found
Usage:
  audiofile get [flags]
Flags:
  -h, --help        help for get
      --id string   audiofile id
      --json        return json format
```

We can even level up by additionally suggesting how to find an ID. Potentially, ask the user to run the `list` command to confirm the ID. Another thing that can be done, similarly to how we handled the status codes from an HTTP API request, is to check the errors coming back from a local command being called. Whether the command is not found or the command is missing executable permissions, you can similarly use a switch to handle potential errors that can occur when a command is started or run. These potential errors can be rewritten similarly using more user-friendly language.

Just Imagine

# Providing debug and traceback information

Debug and traceback information is mostly useful for you or other developers, but it can also help your end users share valuable information with you to help debug a potential issue found in your code. There are several diverse ways to provide this information. Debug and traceback information is primarily output to a log file, and often, the addition of a `verbose` flag will print this output, which is usually hidden.

## Logging data

Since debug data is usually found in log files, let us discuss how to include logging in your command-line application and determine the levels associated with logging—`info`, `error`, and `debug` levels of severity. In this example, let us use a simple log package to demonstrate this example. There are several different popular structured log packages, including the following:

-   Zap ([https://github.com/uber-go/zap](https://github.com/uber-go/zap))—Fast structured logger developed by Uber
-   ZeroLog ([https://github.com/rs/zerolog](https://github.com/rs/zerolog))—Fast and simple logger dedicated to JSON format
-   Logrus ([https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus))—Structured logger for Go with the option for JSON-formatted output (currently in maintenance mode)

Although `logrus` is an extremely popular logger, it has not been updated in a while, so let us choose to use `zap` instead. In general, it’s a promising idea to choose an open source project that is actively maintained.

## Initiating a logger

Back to the `audiofile` project, let us add logging for debugging purposes. The very first thing we run under our `audiofile` repo is this:

```markup
go get -u go.uber.org/zap
```

It will get the updated Zap logger dependencies. After that, we can start referencing the import within the project’s Go files. Under the `utils` directory, we add a `utils/logger.go` file to define some code to initiate the Zap logger, which is called within the `main` function:

```markup
package utils
import (
    "go.uber.org/zap"
)
var Logger *zap.Logger
var Verbose *zap.Logger
func InitCLILogger() {
    var err error
    var cfg zap.Config
    config := viper.GetStringMap("cli.logging")
    configBytes, _ := json.Marshal(config)
    if err := json.Unmarshal(configBytes, &cfg); err != nil {
        panic(err)
    }
    cfg.EncoderConfig = encoderConfig()
    err = createFilesIfNotExists(cfg.OutputPaths)
    if err != nil {
        panic(err)
    }
    cfg.Encoding = "json"
    cfg.Level = zap.NewAtomicLevel()
    Logger, err = cfg.Build()
    if err != nil {
        panic(err)
    }
    cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
    Verbose, err = cfg.Build()
    if err != nil {
        panic(err)
    }
    defer Logger.Sync()
}
```

It isn’t necessary, but we define two loggers here. One is a logger, `Logger`, which outputs to an output path defined within the config file, and the other is the verbose logger, `Verbose`, which outputs to standard output and the previously defined output path. Both use the `*zap.Logger` type, which is used when type safety and performance are critical. Zap also provides a sugared logger, which is used when performance is nice to have but not critical. `SugarLogger` also allows for structured logging, but in addition, supports `printf`\-style APIs.

Within the `Chapter-9` branch version of this repo, we replace some of the general `fmt.Println` or `fmt.Printf` output with the logs that can be shown in `verbose` mode. Also, we differentiate when printing out information with the `Info` level versus the `Error` level.

The following code uses Viper to read from the configuration file, which has been modified to hold a few extra configurations for the logger:

```markup
{
    "cli": {
        "hostname": "localhost",
        "port": 8000,
        "logging": {
            "level": "debug",
            "encoding": "json",
            "outputPaths": [
                "/tmp/log/audiofile.json"
            ]
        }
    }
}
```

In the preceding configuration, we set the `level` and `encoding` fields. We choose the `debug` level so that debug and error statements are output to the log file. For the `encoding` value, we chose `json` because it provides a standard structure that can make it easier to understand the error message as each field is labeled. The encoder config is also defined within the same `utils/logger.go` file:

```markup
func encoderConfig() zapcore.EncoderConfig {
    return zapcore.EncoderConfig{
        MessageKey: "message",
        LevelKey: "level",
        TimeKey: "time",
        NameKey: "name",
        CallerKey: "file",
        StacktraceKey: "stacktrace",
        EncodeName: zapcore.FullNameEncoder,
        EncodeTime: timeEncoder,
        EncodeLevel: zapcore.LowercaseLevelEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller: zapcore.ShortCallerEncoder,
    }
}
```

Since the `InitCLILogger()` function is called within the `main` function, the two `Logger` and `Verbose` loggers will be available within any of the commands for use.

## Implementing a logger

Let us look at how we can start using this logger in an effective way. First, we know that we are going to log all the data and output to the user when in verbose mode. We define the `verbose` flag as a persistent flag in the `cmd/root.go` file. This means that the `verbose` flag will be available not only at the root level but also for every subcommand added to it. In that file’s `init()` function, we add this line:

```markup
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
```

Now, rather than checking for every error if the `verbose` flag is called and printing out the error before it is returned, we create a simple function that can be repeated for checking but also returning the error value. Within the `utils/errors.go` file, we define the following function for reuse:

```markup
func Error(errString string, err error, verbose bool) error {
    errString = cleanup(errString, err)
    if err != nil {
        if verbose {
            // prints to stdout also
            Verbose.Error(errString)
        } else {
            Logger.Error(errString)
        }
        return fmt.Errorf(errString)
    }
    return nil
}
```

Let’s take one command as an example, the `delete` command, which shows how this function is called:

```markup
var deleteCmd = &cobra.Command{
    Use: "delete",
    Short: "Delete audiofile by id",
    Long: `Delete audiofile by id. This command removes the
        entire folder containing all stored metadata`,
```

The bulk of the code for the command is usually found within the `Run` or `RunE` method, which receives the `cmd` variable, a `*cobra.Command` instance, and the `args` variable, which holds arguments within a slice of `strings`. Very early on, in each method, we create the client and extract any flags we might need—in this case, the `verbose`, `silence`, and `id` flags:

```markup
    RunE: func(cmd *cobra.Command, args []string) error {
        client := &http.Client{
            Timeout: 15 * time.Second,
        }
        var err error
          silence, _ := cmd.Flags().GetBool("silence")
        verbose, _ := cmd.Flags().GetBool("verbose")
        id, _ := cmd.Flags().GetString("id")
        if id == "" {
            id, err = utils.AskForID()
            if err != nil {
                return utils.Error("\n %v\n try again and
                    enter an id", err, verbose)
            }
        }
```

Next, we construct the request we are sending to the `HTTP` client, which uses the `id` value:

```markup
        params := "id=" + url.QueryEscape(id)
        path := fmt.Sprintf("http://%s:%d/delete?%s",
            viper.Get("cli.hostname"),
            viper.GetInt("cli.port"), params)
        payload := &bytes.Buffer{}
        req, err := http.NewRequest(http.MethodGet,
            path, payload)
        if err != nil {
            return utils.Error("\n %v\n check configuration
                to ensure properly configured hostname and
                port", err, verbose)
        }
```

We check whether there’s any error when creating the request, which is most likely a result of a configuration error. Next, we log the request so that we are aware of any communication to external servers:

```markup
        utils.LogRequest(verbose, http.MethodGet, path,
            payload.String())
```

We’ll execute the request through the client’s `Do` method and return an error if the request was unsuccessful:

```markup
        resp, err := client.Do(req)
        if err != nil {
            return utils.Error("\n %v\n check configuration
                to ensure properly configured hostname and
                port\n or check that api is running", err,
                verbose)
        }
        defer resp.Body.Close()
```

Following the request, we check the response and read the `resp.Body` , or the body of the response, if the response was successful. If not, an error message will be returned and logged:

```markup
        err = utils.CheckResponse(resp)
        if err != nil {
            return utils.Error("\n checking response: %v",
            err, verbose)
        }
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return utils.Error("\n reading response: %v
                \n ", err, verbose)
        }
        utils.LogHTTPResponse(verbose, resp, b)
```

Finally, we check whether the response returns the `success` string, which shows a successful deletion. The result is then printed out to the user:

```markup
        if strings.Contains(string(b), "success") && !silence {
            fmt.Printf("\U00002705 Successfully deleted
                audiofile (%s)!\n", id)
        } else {
            fmt.Printf("\U0000274C Unsuccessful delete of
                audiofile (%s): %s\n", id, string(b))
        }
        return nil
    },
}
```

You’ll see that the `utils.Error` function is called every time an error is encountered. You’ll also see a few other logging functions: `utils.LogRequest` and `utils.LogHTTPResponse`. The first, `utils.LogRequest`, is defined to log the request to either standard output, the log file, or both:

```markup
func LogRequest(verbose bool, method, path, payload string) {
    if verbose {
        Verbose.Info(fmt.Sprintf("sending request: %s %s
            %s...\n", method, path, payload))
    } else {
        Logger.Info(fmt.Sprintf("sending request: %s %s
            %s...\n", path, path, payload))
    }
}
```

The second, `utils.LogHTTPResponse`, similarly logs the response from the previous request to either standard output, the log file, or both:

```markup
func LogHTTPResponse(verbose bool, resp *http.Response, body []byte) {
    if verbose && resp != nil {
        Verbose.Info(fmt.Sprintf("response status: %s,
            body: %s", resp.Status, string(body)))
    } else if resp != nil {
        Logger.Info(fmt.Sprintf("response status: %s, body:
            %s", resp.Status, string(body)))
    }
}
```

Now that this logger has been implemented for all the `audiofile` commands, let’s give it a try and see what the output looks like now that the command has a `verbose` flag to output debug data when necessary.

## Trying out verbose mode to view stack traces

After recompiling the project, we run the `delete` command with an invalid ID and pass the `verbose` command:

```markup
./bin/audiofile delete --id invalidID --verbose
{"level":"info","time":"2022-11-06 21:21:44","file":"utils/logger.go:112","message":"sending request: GET http://localhost:8000/delete?id=invalidID ...\n"}
{"level":"error","time":"2022-11-06 21:21:44","file":"utils/errors.go:17","message":"checking response: \u001b[41;37mthe id cannot be found\u001b[0m","stacktrace":"github.com/marianina8/audiofile/utils.Error\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/utils/errors.go:17\ngithub.com/marianina8/audiofile/cmd.glob..func2\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/cmd/delete.go:54\ngithub.com/spf13/cobra.(*Command).execute\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/vendor/github.com/spf13/cobra/command.go:872\ngithub.com/spf13/cobra.(*Command).ExecuteC\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/vendor/github.com/spf13/cobra/command.go:990\ngithub.com/spf13/cobra.(*Command).Execute\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/vendor/github.com/spf13/cobra/command.go:918\ngithub.com/marianina8/audiofile/cmd.Execute\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/cmd/root.go:21\nmain.main\n\t/Users/mmontagnino/Code/src/github.com/marianina8/audiofile/main.go:11\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:250"}
Error: checking response: the id cannot be found
Usage:
  audiofile delete [flags]
Flags:
  -h, --help        help for delete
      --id string   audiofile id
Global Flags:
  -v, --verbose   verbose
```

Using the `verbose` flag, the debug statements are printed out, and when an error occurs, the stack trace is also output. This is important data for the user to share with the developer to debug what went wrong. Now, let us learn how to give the option to the user to submit a bug.

Just Imagine

# Effortless bug submission

Let us create a `bug` command using the Cobra generator for users to submit issues to the developers of the `audiofile` CLI:

```markup
cobra-cli add bug
bug created at /Users/mmontagnino/Code/src/github.com/marianina8/audiofile
```

Now that we have the `bug` command created, the `Run` field is changed to extract details of the application and launch a web browser with the data already added and ready for the user to just finish off the submission with some extra details:

```markup
var bugCmd = &cobra.Command{
    Use: "bug",
    Short: "Submit a bug",
    Long: "Bug opens the default browser to start a bug
        report which will include useful system
        information.",
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) > 0 {
            return fmt.Errorf("too many arguments")
        }
        var buf bytes.Buffer
        buf.WriteString(fmt.Sprintf("**Audiofile
            version**\n%s\n\n", utils.Version()))
        buf.WriteString(description)
        buf.WriteString(toReproduce)
        buf.WriteString(expectedBehavior)
        buf.WriteString(additionalDetails)
        body := buf.String()
        url := "https://github.com/marianina8/audiofile/issues/new?title=Bug Report&body=" + url.QueryEscape(body)
        // we print if the browser fails to open
        if !openBrowser(url) {
            fmt.Print("Please file a new issue at https://github.com/marianina8/audiofile/issues/new using this template:\n\n")
        fmt.Print(body)
        }
        return nil
    },
}
```

The strings passed into the `buf.WriteString` method are defined outside the command within the same file, `cmd/bug.go`, but once the command is run, the complete template body is defined as follows:

```markup
**Audiofile version**
1.0.0
**Description**
A clear description of the bug encountered.
**To reproduce**
Steps to reproduce the bug.
**Expected behavior**
Expected behavior.
**Additional details**
Any other useful data to share.
```

Calling the `./bin/audiofile bug` command launches the browser to open a new issue on the GitHub repo:

![Figure 9.2 – Screenshot of browser open to a new issue](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_9.2._B18883.jpg)

Figure 9.2 – Screenshot of browser open to a new issue

From the browser window, open the new issue page; the version of the CLI is populated, and then the user can replace the default text for the description, reproduction steps, expected behavior, and other steps with their own.

Just Imagine

# Help, documentation, and support

Part of creating a CLI that empathizes with its users is to supply sufficient help and documentation, as well as support users of all kinds. Luckily, the Cobra CLI framework supports the generation of help from the short and long fields of the Cobra command and the generation of man pages as well. However, bringing empathy into the extended documentation of your CLI may require several techniques.

## Generating help text

By now, there have been many examples of creating commands, but just to reiterate, the command structure and the fields that show up in help are fields within the Cobra commands. Let’s go over a good example:

```markup
var playCmd = &cobra.Command{
    Use: "play",
    Short: "Play audio file by id",
    Long: `Play audio file by id using the default audio
        player for your current system`,
    Example: `./bin/audiofile play –id
        45705eba-9342-4952-8cd4-baa2acc25188`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // code
    },
}
```

Making sure you simply supply a short and long description of the command and one or several examples, you are supplying some help text that can at least get users started using the command. Running this will show the following output:

```markup
audiofile % ./bin/audiofile play --help
Play audio file by id using the default audio player for your current system
Usage:
  audiofile play [flags]
Examples:
  ./bin/audiofile play –id 45705eba-9342-4952-8cd4-baa2acc25188
Flags:
  -h, --help        help for play
      --id string   audiofile id
Global Flags:
  -v, --verbose   verbose
```

A simple command doesn’t need a ton of explanation, so this is enough to help guide the user with usage.

## Generating man pages

In the `audiofile` repo, we’ve added some additional code to generate the man pages for the existing commands and commands in the `Makefile` to run to quickly run the code to do so. There exists a new program within the repo defined under `documentation/main.go`:

```markup
import (
    "log"
    "github.com/marianina8/audiofile/cmd"
    "github.com/spf13/cobra/doc"
)
func main() {
    header := &doc.GenManHeader{
        Title: "Audiofile",
        Source: "Auto generated by marianina8",
    }
    err := doc.GenManTree(cmd.RootCMD(), header, "./pages")
    if err != nil {
        log.Fatal(err)
    }
}
```

We pass in the `root` command and generate the pages in the `./pages` directory. The addition of the `make pages` command within the `Makefile` creates the man pages when called:

```markup
manpages:
    mkdir -p pages
    go run documentation/main.go
```

Within the terminal, if you run `make manpages` and then check to see whether the new pages exist by running `man pages/audiofile.1`, you will see the generated man page for the `audiofile` CLI:

![Figure 9.3 – Screenshot of audiofile man pages in the terminal](https://static.packt-cdn.com/products/9781804611654/graphics/image/Figure_9.3_B18883.jpg)

Figure 9.3 – Screenshot of audiofile man pages in the terminal

You can also see that within the `pages` directory, there’s an individual man page created for all the commands that have been added to the `root` command.

## Embedding empathy into your documentation

By the time a user reaches your documentation, it is likely that they may already have encountered an issue and are frustrated or confused. It’s important that your documentation takes in that perspective and portrays an understanding of the user’s situation.

Although it may feel like documentation takes a lot of time and energy from other areas of development, it is essential for the future of your command-line application.

Within the last few years, there’s been a recent term, _empathy advocacy_, that has come up in regard to technical documentation. It was coined by Ryan Macklin, a technical and UX writer, as well as an empathy advocate. The term is used to describe a subfield of technical communication centered on empathy and realistic respect for human emotion. It can be considered a framework for the way you communicate with your users. Because many people come to your documentation, we know that there’s a varied assortment of brain chemistry, life experience, and recent events playing in mind. Empathy advocacy is one solution to this beautiful challenge.

Macklin has proposed seven philosophical documentation techniques rooted in empathy advocacy. These principles have been informed by disciplines such as UX, trauma psychotherapy, neurobiology, gameplay design, and cultural and language differences. Let’s discuss each of these tenets and why they work:

-   **Employ visual storytelling**—The human brain easily grabs onto stories, and sighted users can benefit from visuals. However, this forces developers to think about different types of accessibility: visual, cognitive, motor, and so on. Telling a story forces the writer to think about structure. On the other hand, dense and long-winded text is **accessibility-hostile**. As a note, this idea doesn’t work for everyone.
-   **Use synopses**—Using a **tl;dr** (short for **too long, don’t read**), a summary line, or a banner provides a shortened explanation for tired and stressed-out readers who benefit from a lower cognitive cost option. Cognitive glue is required for running a collection of cognitive tasks to complete a high level of intelligence. Cognitive glue requires energy, so providing a synopsis will provide a low-cost option for users who are already running on low.
-   **Give time frames**—In general, uncertainty creates **vicious voids**, and dwelling within the unknown time frame can create heightened emotional responses. Providing time frames can help stabilize the void. Time frames can be given if there’s an outage on the server side, an upload to the server, or just a general time to complete a certain task.
-   **Include short videos**—This is a great alternative for some users who struggle with reading comprehension. Typically, younger audiences are used to video, and when you split videos up into a single topic at max, the shorter playtime can be reassuring. Reassurance is a powerful way to regulate emotion. However, there are some pitfalls to video—mainly, that video costs more time and energy to create.
-   **Reduce screenshots**—Providing screenshots can be helpful, but only when the UI can be confusing. Also, providing just enough for the user to figure some things out themselves helps to foster cognitive glue. Otherwise, being bombarded by visuals hurts everyone.
-   **Rethink FAQs**—Instead of a traditional question and answer, break up documentation into single-scoped documents. Provide specific titles and avoid over-promising.
-   **Pick your battles**—It’s difficult to fight every fight; do the best you can, and choose your battles. Not everything you do will work for everyone—learn along the way. After all, advocating for empathy is another means of self-care.

Hopefully, these tenets that describe the philosophy of empathy advocacy help you to think twice about the words you use in your documentation. A few things to consider when you are writing your documentation include how your words may come across to someone who is in a panicked or frustrated state. Also, consider how you can help those about to give up or lacking the energy to complete their task to be successful.

Just Imagine

# Summary

In this chapter, you have learned specific steps to make your command-line application more empathetic. From error handling, debug and traceback information, effortless bug submission, and empathic advocacy in technical communication, you have learned the technical and empathic skills to apply within your application.

Errors can now be rewritten in color to jump out of the screen and decorated with additional information that provides the user information on exactly where an error has occurred and potentially what they may need to do to reach a resolution. When an error seems unresolvable, the user can then run the same command using the `--verbose` flag and view the detail logs, which might contain server requests and responses necessary to trace more specifically where an error may be happening, down to the line of code.

If a bug is encountered, the addition of a new `bug` command allows the user to spawn a new browser straight from their terminal, opening straight to a new template in GitHub’s new issue submission form.

Finally, bridging the gap between technical documentation and the user’s perspective is done by taking an empathetic approach. Several philosophical tenets when using an empathic framework when writing your documentation were discussed.

Just Imagine

# Questions

1.  Which two common methods can you use for decorating your errors?
2.  Between Zap and Logrus loggers, why would you choose Zap?
3.  What is empathy advocacy?

Just Imagine

# Further reading

-   _Empathy_ _Advocacy_: [https://empathyadvocacy.org](https://empathyadvocacy.org)
-   _Write the_ _Docs_: [https://www.writethedocs.org](https://www.writethedocs.org)

Just Imagine

# Answers

1.  `fmt.Errorf(format string, a ...any) error or errors.Wrap(err error, message` `string) error`.
2.  Zap is faster and is actively maintained.
3.  Empathy advocacy is a sub-field of technical communication centered on empathy and realistic respect for human emotion. It can be considered a framework for the way you write your technical documentation and a solution for writing for many types of people with varied backgrounds and accessibilities.