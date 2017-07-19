package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Projects", func() {
	var (
		projects1, projects2 Projects
	)

	BeforeEach(func() {
		projects1 = Projects{
			Org: "org1",
			Infos: Infos{
				Info{
					Name:        "info0",
					Description: "Info0 description",
					GitUrl:      "info0.git.url",
					TrackerUrl:  "info0.tracker.url",

					OwnerCompany: "info0 Inc.",
					ContactEmail: "info@info0.com",
					ProposedDate: "",

					Repo: nil,
				},
				Info{
					Name:        "info1",
					Description: "Info1 description",
					GitUrl:      "info1.git.url",
					TrackerUrl:  "info1.tracker.url",

					OwnerCompany: "info1 Inc.",
					ContactEmail: "info@info1.com",
					ProposedDate: "",

					Repo: nil,
				},
			},
		}

		projects2 = Projects{
			Org: "org2",
			Infos: Infos{
				Info{
					Name:        "info0",
					Description: "Info0 description",
					GitUrl:      "info0.git.url",
					TrackerUrl:  "info0.tracker.url",

					OwnerCompany: "info0 Inc.",
					ContactEmail: "info@info0.com",
					ProposedDate: "",

					Repo: nil,
				},
				Info{
					Name:        "info1",
					Description: "Info1 description",
					GitUrl:      "info1.git.url",
					TrackerUrl:  "info1.tracker.url",

					OwnerCompany: "info1 Inc.",
					ContactEmail: "info@info1.com",
					ProposedDate: "",

					Repo: nil,
				},
			},
		}
	})

	Context("#Equal", func() {
		It("returns true that projects are equal", func() {
			Expect(projects1.Equal(projects2)).To(BeFalse())
			Expect(projects2.Equal(projects1)).To(BeFalse())

			Expect(projects1.Equal(projects1)).To(BeTrue())

			projects2.Org = "org1"
			Expect(projects1.Equal(projects2)).To(BeTrue())
		})

		It("returns flase when compared to Empty project", func() {
			Expect(projects1.Equal(Projects{})).To(BeFalse())
		})
	})
})
