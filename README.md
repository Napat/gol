# Gol

Gol, the ananym of log, go log, go logrus or whatever, is logrus wrapper.  
This package is focused on writing log to files purpose.

## Features

- Support auto Sync config.
- Auto Rotate log file by date.
- Minimum logging levels: info/error/debug.  

## Example

``` golang
    // Create log dir
    os.MkdirAll("log", os.ModePerm)

    // Log to log/ingress_YYYYMMDD.log, Rotate file by Bangkok date time.
    x, err := gol.NewWithHookAutoFileDate("log/ingress.log", "Asia/Bangkok")
    if err != nil {
        panic(err)
    }

    // Log to log/egress_YYYYMMDD.log, Rotate file by Bangkok date time.
    y, err := gol.NewJSONWithHookAutoFileDate("log/egress.log", "Asia/Bangkok")
    if err != nil {
        panic(err)
    }

    x.GDebugf("GInfof debug")

    x.WithFields(gol.Fields{
        "isAwesome": true,
        "star":      9999,
    }).GInfof("Thailand has many awesome places to explore")

    x.WithFields(gol.Fields{
        "isBeauty": true,
        "star":      9999,
    }).GErrorf("Music is beautiful in a way that nothing else can be")

    y.WithFields(gol.Fields{
        "service":     "golapp",
        "status code": 0,
    }).GDebugf("Healthy")
```

## References

- https://dave.cheney.net/2015/11/05/lets-talk-about-logging
- https://www.scalyr.com/blog/getting-started-quickly-with-go-logging
- https://github.com/natefinch/lumberjack/issues/54
- https://github.com/sirupsen/logrus/issues/1012
- https://github.com/sirupsen/logrus/issues/784
- https://github.com/snowzach/rotatefilehook
- https://github.com/natefinch/lumberjack
