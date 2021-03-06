/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"math"
	"os"

	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/spf13/pflag"
)

func main() {
	sherpa.Execute(func() error {
		var cc count.ClassCounter

		flagSet := pflag.NewFlagSet("Class Counter", pflag.ExitOnError)
		flagSet.IntVar(&cc.JVMClassCount, "jvm-class-count", 0, "the number of classes in the JVM")
		flagSet.StringVar(&cc.SourcePath, "source", "", "path to application to count classes in")

		if err := flagSet.Parse(os.Args[1:]); err != nil {
			return fmt.Errorf("unable to parse flags: %w", err)
		}

		c, err := cc.Execute()
		if err != nil {
			return err
		}

		fmt.Println(math.Ceil(float64(c) * count.LoadFactor))
		return nil
	})
}
