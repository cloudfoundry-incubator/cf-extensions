package bot

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/cf-extensions/models"
)

var _ = Describe("utilities", func() {

	var infos models.Infos

	BeforeEach(func() {
		infos = models.Infos{
			models.Info{
				Name:        "info0",
				Description: "Info0 description",
				GitUrl:      "info0.git.url",
				TrackerUrl:  "info0.tracker.url",

				LeadCompany:  "info0 Inc.",
				ContactEmail: "info@info0.com",
				ProposedDate: "",

				Repo: nil,
			},
			models.Info{
				Name:        "info1",
				Description: "Info1 description",
				GitUrl:      "info1.git.url",
				TrackerUrl:  "info1.tracker.url",

				LeadCompany:  "info1 Inc.",
				ContactEmail: "info@info1.com",
				ProposedDate: "",

				Repo: nil,
			},
		}
	})

	Context("#Length", func() {
		It("returns the length of the Infos", func() {
			Expect(Length(infos)).To(Equal(2))
			Expect(Length(infos)).To(Equal(infos.Len()))
		})
	})

	Context("#CurrentTime", func() {
		It("return current time", func() {
			Expect(CurrentTime()).ToNot(Equal(time.Now()))
		})
	})

	Context("#FormatAsDate", func() {
		It("formats time.Time as date string", func() {
			currentTime := time.Date(2017, 7, 19, 0, 0, 0, 0, time.Local)
			Expect(FormatAsDate(currentTime)).To(Equal("7/19/2017"))
		})
	})

	Context("#FormatAsDateTime", func() {
		It("formats time.Time as date and time string", func() {
			currentTime := time.Date(2017, 7, 19, 11, 30, 0, 0, time.Local)
			Expect(FormatAsDateTime(currentTime)).To(Equal("7/19/2017 @ 11:30:0"))
		})
	})

	Context("#ParseAsDate", func() {
		It("parses date and date/time string as a time.Time", func() {
			//Fail("Implement me!")
		})
	})
})
