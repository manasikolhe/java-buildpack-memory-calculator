/*
 * Copyright 2015-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calculator_test

import (
	"testing"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/calculator"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestCalculator(t *testing.T) {
	spec.Run(t, "Calculator", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var c calculator.Calculator

		it.Before(func() {
			h := flags.HeadRoom(0)
			j := flags.JVMOptions{}
			l := flags.LoadedClassCount(1000)
			t := flags.ThreadCount(10)
			m := flags.TotalMemory(500 * memory.Mibi)

			c = calculator.Calculator{HeadRoom: &h, JvmOptions: &j, LoadedClassCount: &l, ThreadCount: &t, TotalMemory: &m}
		})

		it("uses default and calculated values", func() {
			g.Expect(c.Calculate()).To(ConsistOf(
				memory.DefaultMaxDirectMemory,
				memory.MaxMetaspace(19800000),
				memory.DefaultReservedCodeCache,
				memory.DefaultStack,
				memory.MaxHeap(231858240),
			))
		})

		it("uses configured direct memory", func() {
			d := memory.MaxDirectMemory(memory.Mibi)
			c.JvmOptions.MaxDirectMemory = &d

			g.Expect(c.Calculate()).To(ConsistOf(
				memory.MaxMetaspace(19800000),
				memory.DefaultReservedCodeCache,
				memory.DefaultStack,
				memory.MaxHeap(241295424),
			))
		})

		it("uses configured metaspace", func() {
			m := memory.MaxMetaspace(memory.Mibi)
			c.JvmOptions.MaxMetaspace = &m

			g.Expect(c.Calculate()).To(ConsistOf(
				memory.DefaultMaxDirectMemory,
				memory.DefaultReservedCodeCache,
				memory.DefaultStack,
				memory.MaxHeap(250609664),
			))
		})

		it("uses configured reserved code cache", func() {
			r := memory.ReservedCodeCache(memory.Mibi)
			c.JvmOptions.ReservedCodeCache = &r

			g.Expect(c.Calculate()).To(ConsistOf(
				memory.DefaultMaxDirectMemory,
				memory.MaxMetaspace(19800000),
				memory.DefaultStack,
				memory.MaxHeap(482467904),
			))
		})

		it("uses configured stack", func() {
			s := memory.Stack(memory.Mibi)
			c.JvmOptions.Stack = &s

			g.Expect(c.Calculate()).To(ConsistOf(
				memory.DefaultMaxDirectMemory,
				memory.MaxMetaspace(19800000),
				memory.DefaultReservedCodeCache,
				memory.MaxHeap(231858240),
			))
		})

		it("uses configured heap", func() {
			h := memory.MaxHeap(memory.Mibi)
			c.JvmOptions.MaxHeap = &h

			g.Expect(c.Calculate()).To(ConsistOf(
				memory.DefaultMaxDirectMemory,
				memory.MaxMetaspace(19800000),
				memory.DefaultReservedCodeCache,
				memory.DefaultStack,
			))
		})

		it("returns error if overhead is too large", func() {
			m := memory.MaxMetaspace(500 * memory.Mibi)
			c.JvmOptions.MaxMetaspace = &m

			_, err := c.Calculate()
			g.Expect(err).To(HaveOccurred())
		})

		it("returns error if configured heap is too large", func() {
			h := memory.MaxHeap(500 * memory.Mibi)
			c.JvmOptions.MaxHeap = &h

			_, err := c.Calculate()
			g.Expect(err).To(HaveOccurred())
		})
	})
}
