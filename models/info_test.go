package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Infos", func() {
	var infos Infos

	BeforeEach(func() {
		infos = Infos{
			Info{
				Name:        "info0",
				Description: "Info0 description",
				GitUrl:      "info0.git.url",
				TrackerUrl:  "info0.tracker.url",

				OwnerCompany: "info0 Inc.",
				ContactEmail: "info@info0.com",
				ProposedDate: "",

				Repo:              nil,
				LatestRepoRelease: nil,
			},
			Info{
				Name:        "info1",
				Description: "Info1 description",
				GitUrl:      "info1.git.url",
				TrackerUrl:  "info1.tracker.url",

				OwnerCompany: "info1 Inc.",
				ContactEmail: "info@info1.com",
				ProposedDate: "",

				Repo:              nil,
				LatestRepoRelease: nil,
			},
		}
	})

	Context("#Len", func() {
		It("return the length of array", func() {
			Expect(infos.Len()).To(Equal(2))
		})
	})

	Context("#Swap", func() {
		It("swap the two Info objects in array", func() {
			infos.Swap(0, 1)
			Expect(infos[0].Name).To(Equal("info1"))
			Expect(infos[1].Name).To(Equal("info0"))
		})
	})

	Context("#Less", func() {
		It("return the Info with ordered name", func() {
			Expect(infos.Less(0, 1)).To(BeTrue())
			Expect(infos.Less(1, 0)).To(BeFalse())
		})
	})
})
