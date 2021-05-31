package dictionary

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelloWorld(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dictionary Test")
}

var _ = Describe("Parsing file", func() {
	Context("Test Write File", func() {
		file := "./synonyms.txt"
		text := "love is text"

		It("write first line", func() {
			Expect(writeToFile(file, text)).Should(BeNil())
			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(text))
		})

		It("write to specific line", func() {
			keyword := "love"
			text = "text"
			expected := "love is text,text"

			err := appendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to new keyword", func() {
			keyword := "newkeyword"
			text = "text"
			expected := "love is text,text\nnewkeyword,text"

			err := appendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific new line", func() {
			keyword := "newkeyword"
			text = "new_text"
			expected := "love is text,text\nnewkeyword,text,new_text"

			err := appendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific original line", func() {
			keyword := "love"
			text = "new_text"
			expected := "love is text,text,new_text\nnewkeyword,text,new_text"

			err := appendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})

		It("write to specific multiple line", func() {
			keyword := "newkeyword2"
			text = "new_text"
			expected := "love is text,text\nnewkeyword,text,new_text\nnewkeyword2,new_text"

			err := appendToDictionary(file, keyword, text)

			data, err := ioutil.ReadFile(file)
			Expect(err).Should(BeNil())
			Expect(string(data)).Should(Equal(expected))
		})
	})
})
