# Overview
glogrotate (glog rotate) is a log clear and rotate for glog written in Go. 
- This library support for deleting old logs.
- This library support for rotating logs.

# Installation

```go get github.com/cleey/glogrotate```

# Usage
```
import(
    "github.com/cleey/glogrotate"
)
```
## Example 1 · Recommend
this example will rotate every hour(default) and clean the log files that creation time older than 24 hours.

```
glogrotate.Start(glogrotate.RotateOption{
    Remain: time.Duration(time.Hour * 24),
})
```

## Example 2 · With all args
this example will rotate every 3 seconds and clean the log files that creation time older than 10 seconds.

```
glogrotate.Start(glogrotate.RotateOption{
    Dir:    logDir,      // default is glog log_dir
    Prefix: programName, // default is current programName
    RotateInterval: time.Duration(time.Second * 3), // default is 1 hour
    CleanInterval:  time.Duration(time.Second * 3), // default is 1 minute
    Remain: time.Duration(time.Second * 10),
})
```

# Contributing
Contributions are welcome! Open a pull request to fix a bug, or open an issue to discuss a new feature or change.

# Licenses
This program is under the terms of the MIT License. See LICENSE for the full license text.
