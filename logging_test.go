
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
	"testing"
)

type TestWriter struct {
	Data string
}

func (this *TestWriter) Write(p []byte) (n int, e error) {
	this.Data += string(p[:])

	return len(p), nil
}

func TestInfo(t *testing.T) {
	out := TestWriter{}
	logger := GetLogger("test")
	logger.SetOutput(&out)
	logger.SetFormat(0)

	expected := "Info: We expect this text with also the number 36\n"
	logger.Info("We expect ", "this text with also the number ", 36)

	logger.SetLevel(LevelOff)
	logger.Info("This should be ignored!")

	if out.Data != expected {
		t.Errorf("Expected '%s', received '%s'", expected, out.Data)
	}
}

func TestError(t *testing.T) {
	out := TestWriter{}
	logger := GetLogger("test")
	logger.SetOutput(&out)
	logger.SetFormat(0)

	logger.SetLevel(LevelError)
	expected := "Error: We expect this text with also the number 36\n"
	logger.Error("We expect ", "this text with also the number ", 36)

	logger.SetLevel(LevelInfo)
	logger.Error("This should be ignored!")

	if out.Data != expected {
		t.Errorf("Expected '%s', received '%s'", expected, out.Data)
	}
}

func TestDebug(t *testing.T) {
	out := TestWriter{}
	logger := GetLogger("test")
	logger.SetOutput(&out)
	logger.SetFormat(0)

	logger.SetLevel(LevelDebug)
	expected := "Debug: We expect this text with also the number 36\n"
	logger.Debug("We expect ", "this text with also the number ", 36)

	logger.SetLevel(LevelError)
	logger.Debug("This should be ignored!")

	if out.Data != expected {
		t.Errorf("Expected '%s', received '%s'", expected, out.Data)
	}
}
