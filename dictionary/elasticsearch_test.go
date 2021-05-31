package dictionary

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestElasticsearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Elasticsearch Test")
}

var _ = Describe("Running Elasticsearch", func() {
	Context("Test CRUD", func() {
		c := Elastic{}
		tag := Tag{
			Name: "test.md",
			Tags: "## test",
		}
		It("Init", func() {
			Expect(c.Init()).Should(BeNil())
		})
		It("Put", func() {
			Expect(c.Set(tag)).Should(BeNil())
		})
		It("Get", func() {
			Expect(c.Get(tag.Name)).ShouldNot(BeNil())
		})
		It("Delete", func() {
			Expect(c.Delete(tag.Name)).Should(BeNil())
		})
		It("Delete again", func() {
			Expect(c.Delete(tag.Name)).ShouldNot(BeNil())
		})
	})
})
