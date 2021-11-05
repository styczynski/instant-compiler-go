package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/fatih/color"
	"github.com/spf13/afero"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

var callErrorDefaultBg = color.New(color.Reset).SprintFunc()
var callErrorDefaultFg = color.New(color.FgHiRed).SprintFunc()

var callErrorCommandBg = color.New(color.BgBlack).SprintFunc()
var callErrorCommandFg = color.New(color.Underline).SprintFunc()

var callErrorOutBg = color.New(color.Reset).SprintFunc()
var callErrorOutFg = color.New(color.FgBlue).SprintFunc()

type BuildContext struct {
	Cmd         *sh.Session
	Fs          afero.Fs
	Out         afero.Fs
	variables   map[string]string
	tmpBuildDir string
	outLoc      string
	outputFiles map[string]string
}

func CreateBuildContext(program type_checker.LatteTypecheckedProgram, c *context.ParsingContext) *BuildContext {
	filePath := program.Program.Filename()
	fileBase := filepath.Base(filePath)
	fileLoc := filepath.Dir(filePath)
	fileBaseName := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))

	dname, err := os.MkdirTemp("", fmt.Sprintf("_compile_%s", fileBase))
	if err != nil {
		panic(err)
	}

	rootPath, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	variables := map[string]string{
		"INPUT_FILE_NAME": fileBase,
		"INPUT_FILE_LOC":  fileLoc,
		"INPUT_FILE_BASE": fileBaseName,
		"BUILD_DIR":       dname,
		"ROOT":            rootPath,
	}

	outFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(rootPath, fileLoc))

	fs := afero.NewBasePathFs(afero.NewOsFs(), dname)
	session := sh.NewSession()
	for varName, varContent := range variables {
		session = session.SetEnv(varName, varContent)
	}

	session.ShowCMD = false
	session.PipeFail = true

	return &BuildContext{
		Cmd:         session,
		Fs:          fs,
		Out:         outFs,
		tmpBuildDir: dname,
		variables:   variables,
		outLoc:      filepath.Join(rootPath, fileLoc),
		outputFiles: map[string]string{},
	}
}

func (c *BuildContext) Dispose() {
	os.RemoveAll(c.tmpBuildDir)
	c.Cmd = nil
	c.Fs = nil
}

func (c *BuildContext) WriteBuildFile(name string, content []byte) {
	afero.WriteFile(c.Fs, name, content, 0644)
}

func (c *BuildContext) WriteOutput(description string, extension string, content []byte) {
	afero.WriteFile(c.Out, fmt.Sprintf("%s.%s", c.variables["INPUT_FILE_BASE"], extension), content, 0644)
	c.outputFiles[fmt.Sprintf("%s.%s", c.variables["INPUT_FILE_BASE"], extension)] = description
}

func (c *BuildContext) GetOutputFiles() map[string]map[string]string {
	return map[string]map[string]string{
		c.outLoc: c.outputFiles,
	}
}

func (c *BuildContext) ReadBuildFile(name string) []byte {
	out, err := afero.ReadFile(c.Fs, name)
	if err != nil {
		panic(err)
	}
	return out
}

func (c *BuildContext) GetVariable(varName string) string {
	if val, ok := c.variables[varName]; ok {
		return val
	}
	return ""
}

func (c *BuildContext) Substitute(content string) string {
	for varName, varContent := range c.variables {
		content = strings.ReplaceAll(content, "$"+varName, varContent)
	}
	return content
}

func (c *BuildContext) Call(name string, errorPattern string, args ...interface{}) *CompilationError {
	runArgs := []interface{}{}
	for _, cmd := range args {
		arg := c.Substitute(fmt.Sprintf("%v", cmd))
		runArgs = append(runArgs, arg)
	}

	out, err := c.Cmd.Command(name, runArgs...).CombinedOutput()
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	if errorPattern != "" && strings.Contains(string(out), errorPattern) {
		errorMessage = "Output contains errors"
	}

	if errorMessage != "" {
		cmds := []string{}
		for _, cmd := range args {
			cmds = append(cmds, fmt.Sprintf("%v", cmd))
		}
		commandStr := callErrorCommandBg(callErrorCommandFg(fmt.Sprintf(" %s %s ", name, strings.Join(cmds, " "))))
		outLines := strings.Split(string(out), "\n")
		formattedOutLines := []string{callErrorOutBg(callErrorOutFg("    |> "))}
		for _, line := range outLines {
			formattedOutLines = append(formattedOutLines, callErrorOutBg(callErrorOutFg(fmt.Sprintf("    |> %s", line))))
		}

		messageLines := []string{
			callErrorDefaultBg(callErrorDefaultFg(fmt.Sprintf("    | On command: %s\n", commandStr))),
		}
		for varName, varContent := range c.variables {
			messageLines = append(messageLines, callErrorDefaultBg(callErrorDefaultFg(fmt.Sprintf("    |     %s = %s\n", varName, varContent))))
		}

		messageLines = append(messageLines, []string{
			callErrorDefaultBg(callErrorDefaultFg(fmt.Sprintf("    | Error: %v\n", errorMessage)),
				fmt.Sprintf("    %s \n%s", callErrorDefaultBg(callErrorDefaultFg("| Program output:")), strings.Join(formattedOutLines, "\n"))),
		}...)

		return CreateCompilationError(
			"Build command has failed",
			strings.Join(messageLines, ""))
	}
	return nil
}
