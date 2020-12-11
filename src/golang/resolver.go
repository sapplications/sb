package golang

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sapplications/sbuilder/src/smod"
)

type item struct {
	kind     uint
	name     string
	pkg      string
	path     string
	original string
	deps     items
}

type items map[string]item

var itemKind = struct {
	Func   uint
	Struct uint
	String uint
}{
	1,
	2,
	3,
}

type resolver struct {
	configuration string
	entryPoint    string
	config        *smod.ConfigFile
}

func (r *resolver) resolve() (items, error) {
	list := r.getItems()
	// details, err := r.getDetails(list)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(details)
	return list, nil
}

func (r *resolver) getItems() (list items) {
	list = make(items)
	r.getItem(r.entryPoint, list)
	return list
}

func (r *resolver) getItem(itemName string, list items) *item {
	if _, found := list[itemName]; found {
		return nil
	}
	// parse item and add it to the list
	pkg := ""
	name := ""
	kind := itemKind.Struct
	path := ""
	pathSep := "/"
	nameSep := "."
	if strings.HasPrefix(itemName, "\"") {
		kind = itemKind.String
		name = itemName
	} else {
		// get path
		data := strings.Split(itemName, pathSep)
		dataLen := len(data)
		fullName := data[dataLen-1]
		if dataLen > 1 {
			data = data[:dataLen-1]
			path = strings.Join(data, pathSep) + pathSep
		}
		// get pkg and item
		if fullName != "" {
			data = strings.Split(fullName, nameSep)
			dataLen = len(data)
			name = data[dataLen-1]
			if dataLen > 1 {
				pkg = data[0]
			}
		}
		// check and set type of func
		if name != "" && strings.HasSuffix(name, "()") {
			kind = itemKind.Func
		}
	}
	// create an item
	it := item{
		kind,
		name,
		pkg,
		path,
		itemName,
		make(items),
	}
	// process item dependencies
	var refIt *item
	deps := r.config.Items[itemName]
	for dep, res := range deps {
		refIt = r.getItem(res, list)
		if refIt != nil {
			it.deps[dep] = *refIt
		}
	}
	// add item to the result set and return it
	list[itemName] = it
	return &it
}

func (r *resolver) getDetails(list items) (string, error) {
	var unit []string
	unit = append(unit, "package main\n")
	if len(list) > 0 {
		unit = append(unit, "import (")
		unit = append(unit, "\t\"fmt\"")
		unit = append(unit, "\t\"reflect\"")
		// unit = append(unit, "\t\"strconv\"")
		for _, i := range list {
			if i.path != "" && i.pkg != "" {
				unit = append(unit, fmt.Sprintf("\t\"%s%s\"", i.path, i.pkg))
			}
		}
		unit = append(unit, ")\n")
	}
	unit = append(unit, "func main() {")
	if len(list) > 0 {
		for _, x := range list {
			switch x.kind {
			case itemKind.Func:
				unit = append(unit, r.checkFunc(x))
			case itemKind.Struct:
				unit = append(unit, r.checkStruct(x))
			case itemKind.String:
				unit = append(unit, r.checkString(x))
			default:
				fmt.Printf("\tfmt.Printf(\"%s=undefined\")", x.original)
			}
		}
	}
	unit = append(unit, "}")
	// generate a main unit and run it
	wd, _ := os.Getwd()
	root := filepath.Join(wd, r.configuration, "deps")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.Mkdir(root, os.ModePerm)
	}
	file, err := os.Create(filepath.Join(root, "main.go"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(strings.Join(unit, "\n"))
	writer.Flush()
	if err := goBuild(root, root); err != nil {
		return "", err
	}
	out, err := exec.Command(filepath.Join(root, "deps.exe")).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (r *resolver) checkFunc(x item) string {
	name := x.name
	if strings.HasSuffix(name, "()") {
		name = name[0 : len(name)-2]
	}
	return fmt.Sprintf("\tfmt.Printf(\"%s=", x.original) + "%s\\n\", " +
		fmt.Sprintf("reflect.ValueOf(%s.%s).Kind())", x.pkg, name)
}

func (r *resolver) checkStruct(x item) string {
	return fmt.Sprintf("\tfmt.Printf(\"%s=\")\n", x.original) +
		fmt.Sprintf("\tif reflect.ValueOf((*%s.%s)(nil)).CanInterface() {\n", x.pkg, x.name) +
		fmt.Sprint("\t\tfmt.Printf(\"struct\\n\")\n") +
		fmt.Sprint("\t} else {\n") +
		fmt.Sprint("\t\tfmt.Printf(\"undefined\\n\")\n") +
		fmt.Sprint("\t}\n")
}

func (r *resolver) checkString(x item) string {
	str := x.original
	if strings.HasPrefix(str, "\"") {
		str = str[1:]
	}
	if strings.HasSuffix(str, "\"") {
		str = str[0 : len(str)-1]
	}
	return fmt.Sprintf("\tfmt.Printf(\"\\\"%s\\\"=string\\n\")", str)
}