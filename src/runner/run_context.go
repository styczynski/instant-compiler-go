package runner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/fatih/color"
	"github.com/spf13/afero"
	"github.com/styczynski/latte-compiler/src/compiler"
)

var callErrorDefaultBg = color.New(color.Reset).SprintFunc()
var callErrorDefaultFg = color.New(color.FgHiRed).SprintFunc()

var callErrorCommandBg = color.New(color.BgBlack).SprintFunc()
var callErrorCommandFg = color.New(color.Underline).SprintFunc()

var callErrorOutBg = color.New(color.Reset).SprintFunc()
var callErrorOutFg = color.New(color.FgBlue).SprintFunc()

type CompiledCodeRunContext struct {
	outputFiles  map[string]string
	compilerMeta map[string]interface{}
	fs           afero.Fs
	cmd          *sh.Session
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func CreateCompiledCodeRunContext(program compiler.LatteCompiledProgram) *CompiledCodeRunContext {
	meta := map[string]interface{}{}
	for k, v := range program.CompilerMeta {
		meta[k] = v
	}

	outputPath := ""
	for ext, path := range program.OutputFilesByExt {
		meta[fmt.Sprintf("OUTPUT_PATH_%s", strip(strings.ToUpper(ext)))] = path
		outputPath = filepath.Dir(path)
	}
	meta["OUTPUT_DIR"] = outputPath

	fs := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), outputPath))

	session := sh.NewSession()
	for varName, varContent := range meta {
		session = session.SetEnv(varName, fmt.Sprintf("%v", varContent))
	}

	return &CompiledCodeRunContext{
		outputFiles:  program.OutputFilesByExt,
		compilerMeta: meta,
		cmd:          session,
		fs:           fs,
	}
}

func (c *CompiledCodeRunContext) ReadFileByExt(extension string) ([]byte, error) {
	out, err := afero.ReadFile(c.fs, c.Substitute("$INPUT_FILE_BASE.%s", extension))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *CompiledCodeRunContext) Substitute(content string, args ...interface{}) string {
	content = fmt.Sprintf(content, args...)
	for varName, varContent := range c.compilerMeta {
		content = strings.ReplaceAll(content, "$"+varName, fmt.Sprintf("%v", varContent))
	}
	return content
}

func (c *CompiledCodeRunContext) GetCompilerMeta(key string) interface{} {
	return c.compilerMeta[key]
}

func (c *CompiledCodeRunContext) GetOutputFilePathByExtension(extension string) string {
	return c.outputFiles["."+extension]
}

func (c *CompiledCodeRunContext) Call(name string, errorPattern string, args ...interface{}) ([]string, *compiler.RunError) {
	runArgs := []interface{}{}
	for _, cmd := range args {
		arg := c.Substitute(fmt.Sprintf("%v", cmd))
		runArgs = append(runArgs, arg)
	}

	out, err := c.cmd.Command(name, runArgs...).CombinedOutput()
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
		for varName, varContent := range c.compilerMeta {
			messageLines = append(messageLines, callErrorDefaultBg(callErrorDefaultFg(fmt.Sprintf("    |     %s = %v\n", varName, varContent))))
		}

		messageLines = append(messageLines, []string{
			callErrorDefaultBg(callErrorDefaultFg(fmt.Sprintf("    | Error: %v\n", errorMessage)),
				fmt.Sprintf("    %s \n%s", callErrorDefaultBg(callErrorDefaultFg("| Program output:")), strings.Join(formattedOutLines, "\n"))),
		}...)

		return nil, compiler.CreateRunError(
			"Run command has failed",
			strings.Join(messageLines, ""))
	}
	return strings.Split(string(out), "\n"), nil
}
