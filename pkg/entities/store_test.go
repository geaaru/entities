/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package entities_test

import (
	//"fmt"
	. "github.com/geaaru/entities/pkg/entities"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store Tests", func() {
	Context("Loading entities via yaml", func() {

		It("Check Store Load", func() {
			store1 := NewEntitiesStore()
			err := store1.Load("../../testing/fixtures")
			Expect(err).Should(BeNil())
			Expect(len(store1.Users)).Should(Equal(2))
			Expect(len(store1.Groups)).Should(Equal(2))
			Expect(len(store1.Shadows)).Should(Equal(2))
			Expect(len(store1.GShadows)).Should(Equal(2))
			Expect(len(store1.Groups["foo"].GetUsers())).Should(Equal(4))
		})

		It("Check Store Merge", func() {
			store2 := NewEntitiesStore()
			err := store2.Load("../../testing/fixtures")
			Expect(err).Should(BeNil())
			Expect(len(store2.Users)).Should(Equal(2))
			Expect(len(store2.Shadows)).Should(Equal(2))
			Expect(len(store2.GShadows)).Should(Equal(2))

			// Check merge
			gid := 1
			err = store2.AddEntity(Group{
				Name:     "foo",
				Password: "yy",
				Gid:      &gid,
				Users:    "one,five",
			})
			Expect(err).Should(BeNil())
			Expect(len(store2.Groups)).Should(Equal(2))
			Expect(len(store2.Groups["foo"].GetUsers())).Should(Equal(5))
			Expect(store2.Groups["foo"].Password).Should(Equal("xx"))
		})

	})
})
