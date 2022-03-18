/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package entities_test

import (
	. "github.com/geaaru/entities/pkg/entities"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Context("Loading entities via yaml", func() {
		p := &Parser{}
		It("understands the user kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/simple/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(UserPasswd).Username).Should(Equal("foo"))
		})
		It("understands the shadow kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/shadow/user.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Shadow).Username).Should(Equal("foo"))
		})
		It("understands the group kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/group/group.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(Group).Name).Should(Equal("foo"))
		})
		It("understands the gshadow kind", func() {
			entity, err := p.ReadEntity("../../testing/fixtures/gshadow/gshadow.yaml")
			Expect(err).Should(BeNil())
			Expect(entity.(GShadow).Name).Should(Equal("test"))
		})
	})
})
