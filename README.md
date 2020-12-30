# fluffy

fluffy is a command line application that allows you to monitor your logs locally. 

### Installation

```
* Install go using: https://golang.org/doc/install
```

### Setup
* Change config.yml for configuration
    ```
    alert-path:  Path of the file where all alerts will be logged
    alert-window: Window to consider for alerts.
    filename: Log file to monitor.
    log-path: Path of the file where all logs will be saved.
    pid: This will be used by CLI to keep track of daemon process.
    report-path: Path of the file where reports will be saved.
    report-time: Report will be generated after every given second.
    threshold: If requests/second in given alert window exceeds threshold then alarm will be generated.
    ```
* After changing config.yml, run 
    ```
    ./start-server.sh
    ```
* This will start a daemon process which will be monitoring logs
* Use **fluffy status** to get the process id for the daemon
* Use **fluffy stop** to stop the monitoring daemon

### Usage

```
Fluffy is a CLI library that is helpful in monitoring logs and generating alerts

Usage:
  fluffy [command]

Available Commands:
  help        Help about any command
  start       Fluffy will start monitoring your logs
  status      Return the status of fluffy
  stop        Stop the currently running fluffy process

Flags:
      --config string   config file (default is $HOME/.fluffy.yaml)
  -h, --help            help for fluffy
  -t, --toggle          Help message for toggle

Use "fluffy [command] --help" for more information about a command.
```

### UI

#### REPORT
```
         (__)
         (oo)
   /------\/
  / |    ||
 *  /\---/\
    ~~   ~~
--------------------------------------------------------------------
                Report for DATE: 2020-12-31 02:43:37
--------------------------------------------------------------------
PATH                    COUNT
--------------------------------------------------------------------
/report/user                    4
/login/user                     3
/logout/user                    2

---ERROR RATE: 30.00%----
```

#### Alert
```
Monitor on alert 80 requests/sec at time 2020-12-30 21:15:54.643223 +0000 UTC
Monitor out of alert 1 requests/sec at 2020-12-30 21:16:09.642211 +0000 UTC
```

### Testing

#### Unit tests
* Run *go test ./...* to run all unit tests
* Run *go test -run {TestMethod}* to run individual tests

#### Complete monitor testing
* Start server using shell script
* Move to directory generator and run *go run main.go*
* It will generate data to given file