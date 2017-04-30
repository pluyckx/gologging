
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


/*
A small main package to debug the file rotater. At this moment IntelliJ idea has no native support for debugging test
code.
 */
package main

import (
	"log"
	"os"
	"github.com/pluyckx/gologging"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(os.Stderr)

	logger := logging.GetLogger("test")
	out, err := logging.NewRotateFile("/tmp/logs/test.log", 100, logging.SIZE_MIN, 5)
	if err != nil {
		log.Println("Could not create rotate file:", err)
	} else {
		err = out.Open()

		if err != nil {
			panic(err)
		} else {
			defer out.Close()
		}

		logger.SetOutput(out)
		logger.SetFormat(log.Ldate | log.Ltime | log.Lshortfile)
		logger.SetLevel(logging.LevelInfo)

		logger.Info("Dit is een test1.")
		logger.Info("Dit is een test2.")
		logger.Info("Dit is een test3.")
		logger.Info("Dit is een test4.")
		logger.Info("Dit is een test5.")
		logger.Info("Dit is een test6.")
		logger.Info("Dit is een test7.")
		logger.Info("Dit is een test8.")
		logger.Info("Dit is een test9.")
		logger.Info("Dit is een test10.")
		logger.Info("Dit is een test11.")
		logger.Info("Dit is een test12.")
		logger.Info("Dit is een test13.")
	}
}
