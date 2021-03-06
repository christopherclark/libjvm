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

package count_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/sclevine/spec"
)

func testClassCounter(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path string
	)

	it.Before(func() {
		var err error

		path, err = ioutil.TempDir("", "class-counter")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	it("counts files in application", func() {
		Expect(ioutil.WriteFile(filepath.Join(path, "alpha.class"), []byte{}, 0644)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(path, "bravo"), 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(path, "bravo", "charlie.class"), []byte{}, 0644)).To(Succeed())

		Expect(count.ClassCounter{SourcePath: path}.Execute()).To(Equal(2))
	})

	it("counts files in archives", func() {
		Expect(count.ClassCounter{SourcePath: "testdata"}.Execute()).To(Equal(2))
	})
}
