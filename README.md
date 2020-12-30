# fluffy

fluffy is a command line application that allows you to monitor your logs locally. 

### Installation

```
* Install go using: https://golang.org/doc/install
* Move this(fluffy) package to $Home
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
* Run *go test ./...* to run all unit test
* Run *go test -run {TestMethod}* to run individual tests

#### Complete monitor testing
* Start server using shell script
* Move to directory generator and run *go run main.go*
* It will generate data to given file

### Improvements
* Use docker conainer to daemonize the process instead of shell script
* Implement Notification service like SNS to send notifications to customers
* Remove heavy dependencies on hpcloud/tail and logparse. Write custom solution instead of using these dependencies directly.
* Use some time-series internal db to serve the queries. Internal db can be used to store the events and query them efficiently. In certain time period this internal db can be flushed to disk, thus avoiding memory issues.
* Better and clean CLI support.

### Scalable Solution

#### Use streaming service
* Pull logs from servers.
* Ingest them to any streaming service like kafka stream or AWS Kinesis Stream.
* Use workers to pull these logs, analyze and persist them on time-series database.
* For alerts, these databases can be polled to check if requests crossed the threshold.


#### Set up your own logging clusters
* It is much better to use prometheus and other similar tools for monitoring rather than parsing logs
* A Prometheus server can be setup which exposes certain metrics api. It supports alerting as well.
* Prometheus stores data in in-memory time-series db, which is regularly flushed to disk.
* Grafana can be used to visualize metrics collection
* AlertManager can be used to send alerts to client.
