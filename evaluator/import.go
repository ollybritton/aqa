package evaluator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ollybritton/aqa/ast"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/parser"
)

func pathToModuleName(path string) string {
	// collatz.aqa => collatz
	// collatz-the-best.aqa => collatz_the_best

	path = strings.TrimSuffix(path, "/")

	extension := filepath.Ext(path)
	components := strings.Split(path, string(filepath.Separator))

	var name string

	if extension == "" && len(components) == 0 {
		// Not a folder, just no extension given
		name = path

	} else if extension == "" && len(components) > 0 {
		// No extension, just use folder name
		name = components[len(components)-1]

	} else if len(extension) > 0 {
		// Full path, just use file name without extension
		last := components[len(components)-1]
		name = last[0 : len(last)-len(extension)]
	} else {
		// Use name of file without extension
		name = path[0 : len(path)-len(extension)]
	}

	reg := regexp.MustCompile("[^a-zA-Z0-9_/]+")

	return reg.ReplaceAllString(name, "_")
}

func evalImport(node *ast.ImportStatement, env *object.Environment) object.Object {
	fi, err := os.Stat(node.Path)
	if err != nil {
		return newError("could not read import %q", node.Path)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return evalDirectoryImport(node, env)
	case mode.IsRegular():
		return evalFileImport(node, env)
	}

	return newError("path specified was not a file or directory: %q", node.Path)
}

func evalFileImport(node *ast.ImportStatement, env *object.Environment) object.Object {
	bytes, err := ioutil.ReadFile(node.Path)
	if err != nil {
		return newError("could not read file %q", node.Path)
	}

	str := string(bytes)

	l := lexer.New(str)
	p := parser.New(l)

	program := p.Parse()
	if len(p.Errors()) > 0 {
		errors := []string{}
		for _, e := range p.Errors() {
			errors = append(errors, e.Error())
		}

		return newError("could not parse file: %v", strings.Join(errors, "\n"))
	}

	fileEnv := object.NewEnvironment()
	eval := Eval(program, fileEnv)

	var moduleName string

	if node.As == "" {
		moduleName = pathToModuleName(node.Path)
	} else {
		moduleName = node.As
	}

	module := &object.Module{
		Env:     fileEnv,
		Exposed: fileEnv.Keys(),

		Path:      node.Path,
		IsBuiltin: false,
	}

	if eval == nil || eval.Type() == "NULL" {
		env.Set(moduleName, module)
		return nil
	}

	if eval.Type() == object.ERROR_OBJ {
		return newError("error importing file, error during evaluation: %v", eval.Inspect())
	}

	return nil
}

func evalDirectoryImport(node *ast.ImportStatement, env *object.Environment) object.Object {
	return newError("directory imports coming soon!")
}
