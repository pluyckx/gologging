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
	"sync"
	"fmt"
	"bytes"
	"errors"
)

const (
	SIZE_STRICT SizeMode = iota
	SIZE_MIN SizeMode = iota
	SIZE_MAX SizeMode = iota
)

const filenameTemplate string = "%s.%d"

var ErrorUnkownSizeMode = errors.New("Unkown size mode")

type SizeMode int

type RotateFile struct {
	path       string
	file       *os.File
	lock       sync.Mutex
	maxSize    int64
	sizeMode   SizeMode
	maxRotates uint32
}

func NewRotateFile(path string, maxSize int64, sizeMode SizeMode, maxRotates uint32) (*RotateFile, error) {
	rf := RotateFile{path:path, file:nil, maxSize:maxSize, sizeMode:sizeMode, maxRotates:maxRotates}

	err := rf.Open()

	if err == nil {
		rf.Close()
		return &rf, nil
	} else {
		return nil, err
	}
}

func (this *RotateFile) GetPath() string {
	return this.path
}

func (this *RotateFile) Close() error {
	return this.file.Close()
}

func (this *RotateFile) Open() error {
	var err error
	this.file, err = os.OpenFile(this.path, os.O_APPEND | os.O_WRONLY, 0664)

	if err != nil {
		if os.IsNotExist(err) {
			this.file, err = os.OpenFile(this.path, os.O_CREATE | os.O_WRONLY, 0664)
		}

		if err != nil {
			this.file = nil
		}
	}

	return err
}

func (this *RotateFile) Write(data []byte) (int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.write(data)
}

func (this *RotateFile) Rotate() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.rotate()
}

func (this *RotateFile) write(data []byte) (int, error) {
	info, err := this.file.Stat()

	if err != nil {
		return 0, err
	}

	currentSize := info.Size()
	newSize := currentSize + int64(len(data))

	switch(this.sizeMode) {
	case SIZE_MAX:
		if newSize > this.maxSize {
			this.rotate()
		}

		return this.file.Write(data)

	case SIZE_MIN:
		if currentSize > this.maxSize {
			this.rotate()
		}

		return this.file.Write(data)

	case SIZE_STRICT:
		if newSize > this.maxSize {
			free := this.maxSize - info.Size()

			var nPart1 int
			if free > 0 {
				nPart1, err = this.file.Write(data[:free])

				if err != nil {
					return nPart1, err
				}
			} else {
				free = 0
			}

			err = this.rotate()
			if err != nil {
				return nPart1, err
			}

			var nPart2 int
			nPart2, err = this.write(data[free:])

			return nPart1 + nPart2, err
		} else {
			return this.file.Write(data)
		}

	default:
		return 0, ErrorUnkownSizeMode
	}
}

func (this *RotateFile) rotate() error {
	this.file.Close()
	var buff bytes.Buffer

	rotatesFound := uint32(0)

	_, err := fmt.Fprintf(&buff, filenameTemplate, this.path, rotatesFound + 1)

	var currentFilename string
	if err == nil {
		currentFilename = buff.String()
	} else {
		return err
	}

	for _, err := os.Stat(currentFilename); err == nil; _, err = os.Stat(currentFilename) {
		rotatesFound += 1
		buff.Reset()
		_, err = fmt.Fprintf(&buff, filenameTemplate, this.path, rotatesFound + 1)
		if err != nil {
			return err
		} else {
			currentFilename = buff.String()
		}
	}

	for i := rotatesFound; i >= this.maxRotates; i -= 1 {
		buff.Reset()
		_, err = fmt.Fprintf(&buff, filenameTemplate, this.path, i)
		if err != nil {
			return err
		}

		err = os.Remove(buff.String())
		if err != nil {
			return err
		}

		rotatesFound -= 1
	}

	buff.Reset()
	_, err = fmt.Fprintf(&buff, filenameTemplate, this.path, rotatesFound + 1)
	if err != nil {
		return err
	}
	newName := buff.String()

	for i := rotatesFound; i > 0; i -= 1 {
		buff.Reset()
		_, err = fmt.Fprintf(&buff, filenameTemplate, this.path, i)
		if err != nil {
			return err
		}
		oldName := buff.String()

		err = os.Rename(oldName, newName)
		if err != nil {
			return err
		}

		newName = oldName
	}

	err = os.Rename(this.path, newName)

	if err != nil {
		return err
	}

	err = this.Open()

	return err
}

