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

package dns_test

import (
	"io/ioutil"
	"os"
	"testing"

	ddns "github.com/miekg/dns"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libjvm/dns"
	"github.com/sclevine/spec"
)

func testLinkLocalDNS(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path string
	)

	it.Before(func() {
		f, err := ioutil.TempFile("", "link-local-dns")
		Expect(err).NotTo(HaveOccurred())

		_, err = f.WriteString("test")
		Expect(err).NotTo(HaveOccurred())

		Expect(f.Close()).To(Succeed())
		path = f.Name()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	it("does not modify file if not link local", func() {
		config := &ddns.ClientConfig{Servers: []string{"1.1.1.1"}}
		l := dns.LinkLocalDNS{Config: config, SecurityPath: path}

		Expect(l.Execute()).To(Succeed())
		Expect(ioutil.ReadFile(path)).To(Equal([]byte("test")))
	})

	it("modifies file if link local", func() {
		config := &ddns.ClientConfig{Servers: []string{"169.254.0.1"}}
		l := dns.LinkLocalDNS{Config: config, SecurityPath: path}

		Expect(l.Execute()).To(Succeed())
		Expect(ioutil.ReadFile(path)).To(Equal([]byte(`test
networkaddress.cache.ttl=0
networkaddress.cache.negative.ttl=0
`)))
	})
}
