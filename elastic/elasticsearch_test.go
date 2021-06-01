package es

import (
	"strings"
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
		c, _ := CreateElasticsearch()
		tag := map[string]interface{}{
			"Name": "test.md",
			"Tags": "## test",
		}
		name := tag["Name"].(string)
		It("Init", func() {
			Expect(c.Init()).Should(BeNil())
		})
		It("Put", func() {
			Expect(c.Set(tag)).Should(BeNil())
		})
		It("Get", func() {
			Expect(c.Get(name)).ShouldNot(BeNil())
		})
		It("Check Update", func() {
			Expect(c.Update(name, "new!")).Should(BeNil())
			name, _ := c.Get(name)
			result := strings.Contains(name, "new!")
			Expect(result).Should(BeTrue())
		})
		It("Delete", func() {
			Expect(c.Delete(name)).Should(BeNil())
		})
		It("Delete again", func() {
			Expect(c.Delete(name)).ShouldNot(BeNil())
		})
	})
})
