// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd

import (
	"io"
	"os"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var out io.Writer = os.Stdout

func Test(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Root Cmd", func() {
		g.It("Should run command", func() {
			var rootCmd = &cobra.Command{
				Use:   "khan",
				Short: "khan handles clans",
				Long:  `Use khan to handle clans for your game.`,
			}
			Execute(rootCmd)
		})
	})
}