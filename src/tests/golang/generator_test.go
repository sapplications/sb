package golang

import (
	"github.com/sapplications/sbuilder/src/golang"
	helper "github.com/sapplications/sbuilder/src/helper/hashicorp/hclog"
	"gopkg.in/check.v1"
)

func (s *GoSuite) TestGenerate(c *check.C) {
	g := golang.Coder{
		Logger: helper.NewStdOut("sgo", helper.LogLever.Debug),
	}
	g.Init(
		map[string]map[string]string{
			"apps": {"test": "github.com/sapplications/sbuilder/src/tests/golang.Item1"},
			"github.com/sapplications/sbuilder/src/tests/golang.Item1": {
				"Int1":      "5",
				"Float1":    "5.02",
				"Field1":    "github.com/sapplications/sbuilder/src/tests/golang.NewField1()",
				"Field1V2":  "github.com/sapplications/sbuilder/src/tests/golang.NewField1V2(\"Ariana\", \"Noha\")",
				"Field2":    "github.com/sapplications/sbuilder/src/tests/golang.NewField2(\"Vitalii\")",
				"Field3":    "github.com/sapplications/sbuilder/src/tests/golang.NewField3(github.com/sapplications/sbuilder/src/tests/golang.Field1)",
				"Runner":    "*github.com/sapplications/sbuilder/src/tests/golang.RunnerImpl",
				"Logger":    "github.com/sapplications/sbuilder/src/helper/hashicorp/hclog.NewFileOut(\"sgo\", 3)",
				"Hello":     "github.com/sapplications/sbuilder/src/tests/golang.Hello()",
				"EmptyFunc": "github.com/sapplications/sbuilder/src/tests/golang.EmptyFunc()"},
		},
	)
	err := g.Generate("test")
	c.Assert(err, check.IsNil)
}
