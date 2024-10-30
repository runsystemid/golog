package golog

import "time"

type LogModel struct {
	TraceID       string        `json:"traceId"`
	CorrelationID string        `json:"correlationId"`
	SrcIP         string        `json:"srcIp"`
	IP            string        `json:"ip"`
	Port          string        `json:"port"`
	Path          string        `json:"path"`
	Method        string        `json:"method"`
	Header        interface{}   `json:"header"`
	Request       interface{}   `json:"request"`
	StatusCode    string        `json:"statusCode"`
	HttpStatus    uint64        `json:"httpStatus"`
	Response      interface{}   `json:"response"`
	ResponseTime  time.Duration `json:"rt"`
	Error         interface{}   `json:"error"`
	OtherData     interface{}   `json:"otherData"`
}

type Config struct {
	// App name
	App string `json:"app"`

	// App Version
	AppVer string `json:"appVer"`

	// Log environment (development or production)
	Env string `json:"env"`

	// Location where the system log will be saved
	FileLocation string `json:"fileLocation"`

	// Location where the tdr log will be saved
	FileTDRLocation string `json:"fileTDRLocation"`

	// Maximum size of a single log file.
	// If the capacity reach, file will be saved but it will be renamed
	// with suffix the current date
	FileMaxSize int `json:"fileMaxSize"`

	// Maximum number of backup file that will not be deleted
	FileMaxBackup int `json:"fileMaxBackup"`

	// Number of days where the backup log will not be deleted
	FileMaxAge int `json:"fileMaxAge"`

	// Log will be printed in console if the value is true
	Stdout bool `json:"stdout"`
}
