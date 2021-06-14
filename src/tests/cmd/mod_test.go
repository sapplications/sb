package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sapplications/sbuilder/src/app"
	src "github.com/sapplications/sbuilder/src/cmd"
	"github.com/sapplications/sbuilder/src/smodule"
	"gopkg.in/check.v1"
)

func (s *CmdSuite) TestModSubcmdMissing(c *check.C) {
	c.Assert(s.Mod(), check.ErrorMatches, src.SubcmdMissing)
}

func (s *CmdSuite) TestModUnknownSubcmd(c *check.C) {
	c.Assert(s.Mod("test"), check.ErrorMatches, fmt.Sprintf(src.UnknownSubcmdF, "test"))
}

// test the init subcommand

func (s *CmdSuite) TestModInitLanguageMissing(c *check.C) {
	c.Assert(s.Mod("init"), check.ErrorMatches, src.LanguageMissing)
}

func (s *CmdSuite) TestModInitLanguageIsNotSupported(c *check.C) {
	c.Assert(s.Mod("init", "delphi"), check.ErrorMatches, fmt.Sprintf(app.LanguageIsNotSupportedF, "delphi"))
}

func (s *CmdSuite) TestModInit(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module
	c.Assert(s.Mod("init", lang()), check.IsNil)
	// read the created module
	m := smodule.Manager{Lang: lang}
	_, err := m.ReadAll(lang())
	if err != nil {
		t, _ := ioutil.ReadFile(smodule.GetModuleFileName(app.DefaultModuleName))
		fmt.Print(string(t))
		c.Error(err)
	}
}

// test the add subcommand

func (s *CmdSuite) TestModAddItemMissing(c *check.C) {
	c.Assert(s.Mod("add"), check.ErrorMatches, src.ItemMissing)
}

func (s *CmdSuite) TestModAddModuleMissing(c *check.C) {
	c.Assert(s.Mod("add", "helloItem"), check.ErrorMatches, src.ModOrDepMissing)
}

func (s *CmdSuite) TestModAddEmpty(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// add a new item
	modName := "new"
	err := s.Mod("add", "helloItem", modName)
	c.Assert(err, check.IsNil)
	c.Assert(smodule.IsModuleExists(modName), check.Equals, true)
	// read the created module
	m := smodule.Manager{Lang: lang}
	_, err = m.ReadAll(lang())
	if err != nil {
		t, _ := ioutil.ReadFile(smodule.GetModuleFileName(modName))
		fmt.Print(string(t))
		c.Error(err)
	}
}

func (s *CmdSuite) TestModAddItem(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module
	c.Assert(s.Mod("init", lang()), check.IsNil)
	// add a new item use a new cmd
	cmd := CmdSuite{}
	cmd.SetUpTest(nil)
	name := "helloItem"
	err := cmd.Mod("add", name, app.DefaultModuleName)
	c.Assert(err, check.IsNil)
	// read the created module
	mod := smodule.Manager{Lang: lang}
	r, err := mod.ReadAll(lang())
	if err != nil {
		t, _ := ioutil.ReadFile(smodule.GetModuleFileName(app.DefaultModuleName))
		fmt.Print(string(t))
		c.Error(err)
	} else {
		// check the added item exist
		c.Assert(r.Items()[name], check.NotNil)
	}
}

func (s *CmdSuite) TestModAddItemDependency(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module
	c.Assert(s.Mod("init", lang()), check.IsNil)
	// add a new dependency item (application) to the main item
	cmd := CmdSuite{}
	cmd.SetUpTest(nil)
	name := "hello"
	resolver := "\"Hello World!\""
	err := cmd.Mod("add", smodule.MainItemName, name, resolver)
	c.Assert(err, check.IsNil)
	// TODO verify the added dependency...
}

// test the del subcommand

func (s *CmdSuite) TestModDelModuleMissing(c *check.C) {
	c.Assert(s.Mod("del", "helloItem"), check.IsNil)
}

func (s *CmdSuite) TestModDelItemMissing(c *check.C) {
	c.Assert(s.Mod("del"), check.ErrorMatches, src.ItemMissing)
}

func (s *CmdSuite) TestModDelItemMissing2(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module use a new cmd
	cmd := CmdSuite{}
	cmd.SetUpTest(nil)
	c.Assert(cmd.Mod("init", lang()), check.IsNil)
	// try to delete the missing item
	err := s.Mod("del", "helloItem")
	c.Assert(err, check.IsNil)
}

func (s *CmdSuite) TestModDelItem(c *check.C) {
	// create a temporary folder and change the current working directory
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(c.MkDir())
	// initialize a new module use a new cmd
	cmd := CmdSuite{}
	cmd.SetUpTest(nil)
	c.Assert(cmd.Mod("init", lang()), check.IsNil)
	// add a new item use a new cmd
	cmd = CmdSuite{}
	cmd.SetUpTest(nil)
	name := "helloItem"
	err := cmd.Mod("add", name, app.DefaultModuleName)
	c.Assert(err, check.IsNil)
	// delete the added item
	err = s.Mod("del", name)
	c.Assert(err, check.IsNil)
	// check the item does not exist
	found, _ := smodule.IsItemExists(lang(), name)
	c.Assert(found, check.Equals, false)
}

// mod del|edit|list
// NameMissing             = "\"--name\" parameter is required"
// ModuleMissing           = "\"--mod\" parameter is required"
// LanguageMissing         = "Language parameter is required"
// ResolverMissing         = "\"--resolver\" parameter is required"
// DependencyMissing       = "\"--dep\" parameter is required"
// ItemDoesNotExistF       = "\"%s\" item does not exist"
// DependencyDoesNotExistF = "\"%s\" dependency item does not exist"

//"unknown flag: --lang go"