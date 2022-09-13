package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sapplications/sb/app"
	"github.com/sapplications/smod/lod"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
)

func (s *CmdSuite) TestCodeEmpty(c *check.C) {
	c.Assert(s.Code(), check.ErrorMatches, fmt.Sprintf(lod.ModuleFilesMissingF, ".*"))
}

func (s *CmdSuite) TestCodeAppMissing(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module use a new cmd
	cmd := CmdSuite{}
	cmd.SetUpTest(nil)
	c.Assert(cmd.Mod("init", app.ModKind.SB), check.IsNil)
	// try to generate the empty module
	c.Assert(s.Code(), check.ErrorMatches, app.AppIsMissing)
}

func (s *CmdSuite) TestCodeSbApp(c *check.C) {
	c.Skip("Needs to fix...")
	return
	// create a temporary folder and change the current working directory
	currFolder, _ := os.Getwd()
	defer os.Chdir(currFolder)
	viper.Set("GOWD", currFolder)
	os.Chdir(c.MkDir())
	// copy sb files
	wd, _ := os.Getwd()
	copyFile(currFolder+"\\..\\app.sb", wd+"\\app.sb")
	// generate application's files
	c.Assert(s.Code("sb"), check.IsNil)
}

func (s *CmdSuite) TestCodeHelloWorldApp(c *check.C) {
	// create a temporary folder and change the current working directory
	currFolder, _ := os.Getwd()
	defer os.Chdir(currFolder)
	viper.Set("GOWD", currFolder)
	os.Chdir(c.MkDir())
	// copy sb files
	wd, _ := os.Getwd()
	copyFile(currFolder+"\\..\\samples\\app.sb", wd+"\\app.sb")
	// generate application's files
	c.Assert(s.Code(), check.IsNil)
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err == nil {
		return ioutil.WriteFile(dst, input, 0644)
	} else {
		return err
	}
}