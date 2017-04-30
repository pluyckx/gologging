
/*
   Copyright 2017 Philip Luyckx

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

package logging

import (
	"os"
	"log"
	"io"
	"fmt"
	"bytes"
)

const (
	LevelOff Level = 0
	LevelInfo Level = 100
	LevelError Level = 200
	LevelDebug Level = 300
)

type Level int32

type logger struct {
	out io.Writer
	level Level
	flags int
}

var loggers map[string]logger = make(map[string]logger)
var prefixes map[Level]string = map[Level]string{
	LevelInfo: "Info",
	LevelError: "Error",
	LevelDebug: "Debug",
}

func GetLogger(name string) *logger {
	logHandle, ok := loggers[name]

	if !ok {
		logHandle = logger{
			out:os.Stdout,
			level:LevelInfo,
			flags: log.Ldate|log.Ltime|log.Lshortfile}
		loggers[name] = logHandle
	}

	return &logHandle
}

func (this *logger) SetFormat(format int) {
	this.flags = format
}

func (this *logger) SetOutput(out io.Writer) {
	this.out = out
}

func (this *logger) SetLevel(level Level) {
	this.level = level
}

func (this *logger) Info(out ...interface{}) {
	this.Log(LevelInfo, out...)
}

func (this *logger) Error(out ...interface{}) {
	this.Log(LevelError, out...)
}

func (this *logger) Debug(out ...interface{}) {
	this.Log(LevelDebug, out...)
}

func (this *logger) Log(level Level, out ...interface{}) {
	if this.level >= level {
		log.SetOutput(this.out)
		log.SetFlags(this.flags)

		tmp := bytes.Buffer{}
		fmt.Fprintf(&tmp, "%s: ", prefixes[level])

		for _,o := range out {
			fmt.Fprint(&tmp, o)
		}

		log.Output(3, tmp.String())
	}
}
