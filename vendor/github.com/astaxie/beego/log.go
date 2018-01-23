// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"strings"

	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// BeeLogger references the used application logger.
var BeeLogger = logs.GetBeeLogger()

// SetLevel sets the global log level used by the simple logger.
func SetLevel(l int) {
	logs.SetLevel(l)
}

// SetLogFuncCall set the CallDepth, default is 3
func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	return logs.SetLogger(adaptername, config)
}

// Emergency logs a message at emergency level.
func Emergency(v ...interface{}) {
	logs.Emergency(generateFmtStr(len(v)), v...)
}

// Alert logs a message at alert level.
func Alert(v ...interface{}) {
	logs.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	logs.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	logs.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	logs.Warning(generateFmtStr(len(v)), v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	logs.Warn(generateFmtStr(len(v)), v...)
}

// Notice logs a message at notice level.
func Notice(v ...interface{}) {
	logs.Notice(generateFmtStr(len(v)), v...)
}

// Informational logs a message at info level.
func Informational(v ...interface{}) {
	logs.Informational(generateFmtStr(len(v)), v...)
}

// Info compatibility alias for Warning()
func Info(v ...interface{}) {
	logs.Info(generateFmtStr(len(v)), v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	logs.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	logs.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

//add by bai

func Emergencyf(format string,v ...interface{}) {
	logs.Emergency(format,v...)
}

func Alertf(format string,v ...interface{}) {
	logs.Alert(format,v...)
}
// Critical logs a message at critical level.
func Criticalf(format string,v ...interface{}) {
	logs.Critical(format,v...)
}

// Error logs a message at error level.
func Errorf(format string,v ...interface{}) {
	logs.Error(format,v...)
}

// Warning logs a message at warning level.
func Warningf(format string,v ...interface{}) {
	logs.Warning(format,v...)
}
// compatibility alias for Warning()
func Warnf(format string,v ...interface{}) {
	logs.Warn(format,v...)
}

func Noticef(format string,v ...interface{}) {
	logs.Notice(format,v...)
}

// Info logs a message at info level.
func Informationalf(format string,v ...interface{}) {
	logs.Informational(format,v...)
}
// compatibility alias for Warning()
func Infof(format string,v ...interface{}) {
	logs.Info(format,v...)
}
// Debug logs a message at debug level.
func Debugf(format string,v ...interface{}) {
	logs.Debug(format,v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Tracef(format string,v ...interface{}) {
	logs.Trace(format,v...)
}
