package cmdler

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cmdler", func() {
	var ()

	BeforeEach(func() {
	})
	Describe("Running Commands", func() {
		Context("one command", func() {
			It("Should save stdout as string", func() {
				input := "hello"
				out := New(fmt.Sprintf("echo -n %s", input))
				out.Run()

				Expect(out.HasRun()).To(BeTrue())
				Expect(out.Stdout).To(Equal(input))
				Expect(out.Stderr).To(BeEmpty())
				Expect(out.Code).To(BeZero())
				Expect(out.HasErrors()).NotTo(BeTrue())
			})
			It("Should capture stdout as a string", func() {
				input := "hello"
				out := New(fmt.Sprintf("echo -n %s 1>&2", input))
				out.Run()

				Expect(out.HasRun()).To(BeTrue())
				Expect(out.Stdout).To(BeEmpty())
				Expect(out.Stderr).To(Equal(input))
				Expect(out.Code).To(BeZero())
				Expect(out.HasErrors()).NotTo(BeTrue())
			})
			It("Should return the exit code of the command", func() {
				code := 7
				out := New(fmt.Sprintf("return %d", code))
				out.Run()

				Expect(out.HasRun()).To(BeTrue())
				Expect(out.Stdout).To(BeEmpty())
				Expect(out.Stderr).To(BeEmpty())
				Expect(out.Code).To(Equal(code))
				Expect(out.HasErrors()).NotTo(BeTrue())
			})
		})
		Context("Pipe", func() {
			It("Should Create a new command with linked to the original one", func() {
				c1 := New("echo 'hi'")
				c2 := c1.Pipe("echo 'hi2'")

				Expect(c2.prev).To(Equal(c1))
			})
			It("Should Run the first command when the second one is run", func() {
				c1 := New("echo 'hi'")
				c2 := c1.Pipe("echo 'hi2'")
				c2.Run()
				Expect(c2.prev).To(Equal(c1))
				Expect(c1.HasRun()).To(BeTrue())
				Expect(c2.HasRun()).To(BeTrue())
			})
			It("Should pass stdout to stdin of the next command", func() {
				input := "hello"
				c1 := New(fmt.Sprintf("echo -n '%s'", input))
				c2 := c1.Pipe("cat")
				c2.Run()
				Expect(c2.prev).To(Equal(c1))
				Expect(c1.HasRun()).To(BeTrue())
				Expect(c2.HasRun()).To(BeTrue())
				Expect(c2.Stdout).To(Equal(input))
			})
			It("Should pass stdout to stdin of the next two commands", func() {
				input := "hello"
				c1 := New(fmt.Sprintf("echo -n '%s'", input))
				c2 := c1.Pipe("cat")
				c3 := c1.Pipe("cat")
				c2.Run()
				c3.Run()
				Expect(c2.prev).To(Equal(c1))
				Expect(c1.HasRun()).To(BeTrue())
				Expect(c2.HasRun()).To(BeTrue())
				Expect(c3.HasRun()).To(BeTrue())
				Expect(c2.Stdout).To(Equal(input))
				Expect(c3.Stdout).To(Equal(input))
			})
			It("Should be able to use static input", func() {
				input := "hello"
				c1 := New(fmt.Sprintf("echo -n '%s'", input))
				c2 := c1.Pipe("cat")
				c3 := c1.Pipe("cat")
				c2.Run()
				c3.Run()
				Expect(c2.prev).To(Equal(c1))
				Expect(c1.HasRun()).To(BeTrue())
				Expect(c2.HasRun()).To(BeTrue())
				Expect(c3.HasRun()).To(BeTrue())
				Expect(c2.Stdout).To(Equal(input))
				Expect(c3.Stdout).To(Equal(input))
			})
      It("Should work with things like curl", func() {
        c := New("curl https://api.openshift.com/api").Pipe("jq -r '.services | .[] | .name'")
        c.Run()
      })
		})
	})
})
