package tui

import (
	"fmt"
	"strings"

	"github.com/GuPoroca/HexTest/pkg/typeDefines"
)

func renderTree(p typeDefines.Project) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📂 %s\n", p.Name))
	for _, s := range p.Suites {
		b.WriteString(fmt.Sprintf(" ├─📁 %s\n", s.Name))
		for _, t := range s.Tests {
			b.WriteString(fmt.Sprintf(" │  ├─🧪 %s\n", t.Name))
			for _, a := range t.Asserts {
				b.WriteString(fmt.Sprintf(" │  │  ├─🔎 %s\n", a.Field))
				for _, c := range a.Checks {
					b.WriteString(fmt.Sprintf(" │  │  │  ├─🔲 %s\n", c.Operand))
					for i := range c.Passed {
						b.WriteString(fmt.Sprintf(" │  │  │  │  ├─✔ %s %v\n", c.Expected, c.Passed[i]))
					}
				}
			}
		}
	}
	return b.String()
}
