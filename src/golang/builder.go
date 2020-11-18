package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sapplications/sbuilder/src/smod"
)

type Builder struct {
	Configuration string
}

func (b *Builder) Build(config *smod.ConfigFile) error {
	var err error
	useCurrentConfig := b.Configuration == ""
	if b.Configuration, err = check(b.Configuration, config); err != nil {
		return err
	}
	// check the golang file with all dependencies is exist
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, b.Configuration, depsFileName)
	if _, err = os.Stat(filePath); err != nil {
		return fmt.Errorf("\"%s\" does not exist. Please use a \"generate\" command to create it.", filePath)
	}
	if !useCurrentConfig {
		// delete the configuration golang file
		if _, err := os.Stat(configFileName); err == nil {
			if err := os.Remove(configFileName); err != nil {
				return err
			}
		}
	}
	g := Generator{
		b.Configuration,
	}
	// generate a golang file and save current configuration
	if err := g.generateConfigFile(); err != nil {
		return err
	}
	// generate a golang main file if it is missing
	if _, err := os.Stat(mainFileName); err != nil {
		if os.IsNotExist(err) {
			if err := g.generateMainFile(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// build the application
	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
