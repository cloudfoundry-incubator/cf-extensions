package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectStatus", func() {
	var (
		projectStatus1, projectStatus2 ProjectStatus
	)

	BeforeEach(func() {
		projectStatus1 = ProjectStatus{
			Name: "name1",
			Status: Status{
				Status:      "status1",
				ChangedDate: "2015-14-07T12:00:00Z07:00",
			},
		}

		projectStatus2 = ProjectStatus{
			Name: "name2",
			Status: Status{
				Status:      "status2",
				ChangedDate: "2015-14-07T12:00:00Z07:00",
			},
		}
	})

	Context("#Equals", func() {
		It("returns true that projects are equal", func() {
			Expect(projectStatus1.Equals(projectStatus2)).To(BeFalse())
			Expect(projectStatus2.Equals(projectStatus1)).To(BeFalse())

			Expect(projectStatus1.Equals(projectStatus1)).To(BeTrue())

			projectStatus2.Name = "name1"
			projectStatus2.Status.Status = "status1"
			Expect(projectStatus1.Equals(projectStatus2)).To(BeTrue())
		})

		It("returns false when compared to Empty project", func() {
			Expect(projectStatus1.Equals(ProjectStatus{})).To(BeFalse())
		})
	})
})

var _ = Describe("ProjectsStatus", func() {
	var (
		projectsStatus1, projectsStatus2 ProjectsStatus
	)

	BeforeEach(func() {
		projectsStatus1 = ProjectsStatus{
			Org: "org1",
			Array: []ProjectStatus{
				{
					Name: "name1",
					Status: Status{
						Status:      "status1",
						ChangedDate: "2015-14-07T12:00:00Z07:00",
					},
				},
			},
		}

		projectsStatus2 = ProjectsStatus{
			Org: "org2",
			Array: []ProjectStatus{
				{
					Name: "name1",
					Status: Status{
						Status:      "status1",
						ChangedDate: "2015-14-07T12:00:00Z07:00",
					},
				},
				{
					Name: "name2",
					Status: Status{
						Status:      "status2",
						ChangedDate: "2015-14-07T12:00:00Z07:00",
					},
				},
			},
		}
	})

	Context("#Equals", func() {
		It("returns true that projects are equal", func() {
			Expect(projectsStatus1.Equals(projectsStatus2)).To(BeFalse())
			Expect(projectsStatus2.Equals(projectsStatus1)).To(BeFalse())

			Expect(projectsStatus1.Equals(projectsStatus1)).To(BeTrue())
		})

		It("returns false when compared to Empty project", func() {
			Expect(projectsStatus1.Equals(ProjectsStatus{})).To(BeFalse())
		})
	})

	Context("#StatusForName", func() {
		It("finds status for existing project", func() {
			status, err := projectsStatus1.StatusForName("name1")
			Expect(err).ToNot(HaveOccurred())
			Expect(status).To(Equal(Status{
				Status:      "status1",
				ChangedDate: "2015-14-07T12:00:00Z07:00",
			}))

			status, err = projectsStatus1.StatusForName("name2")
			Expect(err).To(HaveOccurred())

			status, err = projectsStatus2.StatusForName("name1")
			Expect(err).ToNot(HaveOccurred())
			Expect(status).To(Equal(Status{
				Status:      "status1",
				ChangedDate: "2015-14-07T12:00:00Z07:00",
			}))

			status, err = projectsStatus2.StatusForName("name2")
			Expect(err).ToNot(HaveOccurred())
			Expect(status).To(Equal(Status{
				Status:      "status2",
				ChangedDate: "2015-14-07T12:00:00Z07:00",
			}))
		})
	})
})
