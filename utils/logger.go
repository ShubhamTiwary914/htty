package utils 

import (
	"fmt"
	"os"
	"runtime"
	"time"
	global "htty/globals"
	types "htty/types"
)

/*
	main utility for detailed logging, uses "logger" object to log out to stdout/file/etc...
	
	example usage:  Logf(LOG_INFO, "%s %d", str_, int_) [similar to C printf style]
*/
func Logf(level string, format string, args ...interface{}) {
	if !assertAllowedLogLevel(level) {
		return; 
	}
	ts := time.Now().Format("02-01-2006(15:04:05.000)")
	//source = line that called this method
	pc, _, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			caller = fmt.Sprintf("%s:%d", fn.Name(), line)
		}
	}
	msg := fmt.Sprintf(format, args...)
	global.Logger.Printf("level=%s ts=%s caller=%s msg=%q", level, ts, caller, msg)
}

// Logf with LOGLEVEL=debug 
func Debugf(format string, args ...interface{}) {
	Logf(global.LOG_DEBUG, format, args...)
}

// Logf with LOGLEVEL=info
func Infof(format string, args ...interface{}) {
	Logf(global.LOG_INFO, format, args...)
}

// Logf with LOGLEVEL=error
func Errorf(format string, args ...interface{}) {
	Logf(global.LOG_ERROR, format, args...)
}


/*
	Redirect the logs to some file(default=htty.log) by changing the logger object's sink
	(since TUI blocks stdout, its better to log this way)
*/
func RedirectLogs_toFile(outFile string, overwrite bool) *os.File {
	var flags int
	if overwrite {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}
	logfile, err := os.OpenFile(outFile, flags, 0644)
	if err != nil {
		panic(err)
	}
	global.Logger.SetOutput(logfile)
	return logfile	
}


// internal method for Logf to check if the log is allowed in this "LOGLEVEL" 
// (ex LOGLEVEL=all means allow all, LOGLEVEL=debug means allow only debug) 
func assertAllowedLogLevel(level string) bool {
	if global.LOGLEVEL == global.LOG_ALL {
		return true;	
	}		
	if global.LOGLEVEL == level {
		return true;	
	}
	return false
}

//get path to where the logger was called in the program runtime
func pathBase(fp string) string {
	lastslash := -1
	for id := len(fp) - 1; id >= 0; id-- {
		if fp[id] == '/' {
			lastslash = id
			break
		}
	}
	if lastslash >= 0 && lastslash+1 < len(fp) {
		return fp[lastslash+1:]
	}
	return fp
}


// --------------------------------
// Logging for specific structures
// ex: to print status of HttpType object
// -------------- -------------- --
func HttpObjectLogOut(httpObj types.HttpType) {
	Debugf("HTTP REQUEST ----")
	Debugf("Method: %s", httpObj.Method)
	Debugf("Path: %s", httpObj.Path)
	if len(httpObj.Headers) > 0 {
		Debugf("Headers:")
		for k, v := range httpObj.Headers {
			Debugf(" %s: %s", k, v)
		}
	}
	if httpObj.Body != "" {
		Debugf("Body:")
		Debugf("%s", httpObj.Body)
	}
}

func LogPanelGeometry(label string, dims types.PaneGeometry){
	Debugf("%s {X: %s, Y: %s, Width: %s, Height: %s}", label, dims.X, dims.Y, dims.Width, dims.Height)
}
