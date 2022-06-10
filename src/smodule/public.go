package smodule

import (
	"fmt"
	"os"
	"strings"
)

const (
	MainItemName string = "main"
	// notifications
	ModuleIsCreatedF string = "%s file has been created\n"
	// errors
	ItemExistsF             string = "the %s item already exists in %s module"
	ItemIsMissingF          string = "the %s item does not exist"
	ModuleFilesMissingF     string = "no sb files in %s"
	ModuleLanguageMismatchF string = "the %s language of %s module is mismatch the %s selected language"
)

type Manager struct {
	Lang func() string
}

type Reader interface {
	Lang() string
	Items() map[string]map[string]string
	Dependency(string, string) string
	Main() (map[string]string, error)
}

type Formatter struct {
}

func GetModuleFileName(name string) string {
	if strings.HasSuffix(name, moduleExt) {
		return name
	} else {
		return name + moduleExt
	}
}

func IsModuleExists(name string) bool {
	_, err := os.Stat(GetModuleFileName(name))
	return err == nil
}

func IsItemExists(lang, item string) (bool, string) {
	wd, _ := os.Getwd()
	mods, err := loadModules(lang)
	if (err != nil) && (err.Error() != fmt.Sprintf(ModuleFilesMissingF, wd)) {
		return false, ""
	}
	for _, m := range mods {
		if _, found := m.items[item]; found {
			return true, m.name
		}
	}
	return false, ""
}
