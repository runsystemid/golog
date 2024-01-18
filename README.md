# Logger

This library is used to standardized log message for Runsystem service.
You can use this library to simplify your service logging setup.

## Installation

You can run this in your terminal

```shell
go get -u github.com/runsystemid/golog
```

Import this library in your main function or bootstrap function.

```golang
import "github.com/runsystemid/golog"
```

## Usage

### Config

There some required configs. Make sure to initialize the values first. Best practice is to put the values in config files to make it easier to update.

| Config | Description |
| --- | --- |
| App | Application name |
| AppVer | Application version |
| Env | Environment (development or production) |
| FileLocation | Location where the system log will be saved |
| FileTDRLocation | Location where the tdr log will be saved |
| FileMaxSize | Maximum size in Megabytes of a single log file. If the capacity reach, file will be saved but it will be renamed with suffix the current date |
| FileMaxBackup | Maximum number of backup file that will not be deleted |
| FileMaxAge | Number of days where the backup log will not be deleted |
| Stdout | Log will be printed in console if the value is true |

### Loader

Initialize the loader by using

```golang
loggerConfig := logger.Config{
    App:             yourConfig.AppName,
    AppVer:          yourConfig.AppVersion,
    Env:             yourConfig.Environment,
    FileLocation:    yourConfig.Logger.FileLocation,
    FileTDRLocation: yourConfig.Logger.FileTDRLocation,
    FileMaxSize:     yourConfig.FileMaxSize
    FileMaxBackup:   yourConfig.FileMaxBackup
    FileMaxAge:      yourConfig.FileMaxAge,
    Stdout:          yourConfig.Stdout,
}

logger.Loader(loggerConfig)
```

Now you can use the logger from anywhere in your project.

### Print Log

This library provide 2 kind of logs. System log and TDR log.
Besides TDR function, it will print to system log.

## Output

### File and Stdout (configurable)

In default, the file will be printed in file. You can decide whether the output will be printed in console or not.
Just put true in Stdout config attribute to print the log in console.

### Rotate

Log file will be auto rotated based on file size. If the file size bigger than the config, will be saved to a new file with additional date in file name.

## Contributing

Contributions are welcome! Please follow the [Contribution Guidelines](CONTRIBUTION.md).

## License

This project is licensed under the MIT License.