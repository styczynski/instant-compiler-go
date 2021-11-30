package compiler_pipeline

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type RunCompilerPipelineInteractiveCliInterface struct {
	compilerPipeline *HeadlessCompilerPipeline

	gui *gocui.Gui

	editorView    *gocui.View
	outputView    *gocui.View
	runResultView *gocui.View
}

func (c *RunCompilerPipelineInteractiveCliInterface) handleInputChange() {
	input := strings.Join(c.editorView.BufferLines(), "\n")
	formattedText, err := c.compilerPipeline.Printer.FormatRaw(input, true)
	// if err != nil {
	// 	c.errView.Clear()
	// 	fmt.Fprintf(c.errView, "Error: %s", err.Error())
	// } else {
	// 	c.errView.Clear()
	// 	fmt.Fprintf(c.errView, "No errors")
	// 	c.editorView.Clear()
	if err == nil {
		c.editorView.Clear()
		fmt.Fprintf(c.editorView, strings.ReplaceAll(formattedText, "\\033", "\033"))
	}

	c.compilerPipeline.ProcessAsync(input, func(response CompilationResponse) {
		c.gui.Update(func(g *gocui.Gui) error {
			c.outputView.Clear()
			c.runResultView.Clear()
			if response.Ok {
				fmt.Fprintf(c.runResultView, "%s", strings.Join(response.Program.ProgramOutput, "\n"))
				compiledProgramText := response.Program.Program.CompiledProgram.ProgramToText()
				fmt.Fprintf(c.outputView, "%s", compiledProgramText)
			} else {
				fmt.Fprintf(c.outputView, "%s", response.Summary)
			}
			return nil
		})
	})
}

func (c *RunCompilerPipelineInteractiveCliInterface) Run() {
	c.compilerPipeline = CreateHeadlessCompilerPipeline()
	var err error

	c.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer c.gui.Close()

	c.gui.SetManagerFunc(c.layout)
	c.gui.InputEsc = true
	c.gui.Cursor = true

	if err := c.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		log.Panicln(err)
	}

	if err := c.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (c *RunCompilerPipelineInteractiveCliInterface) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	var err error
	if c.editorView, err = g.SetView("editor", 0, 0, maxX/2, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if c.outputView, err = g.SetView("output", maxX/2+1, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		if c.runResultView, err = g.SetView("run-output", 0, maxY-6, maxX/2, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			c.runResultView.Autoscroll = true
			c.runResultView.Wrap = true
			c.runResultView.Title = "Program output:"

			c.outputView.Autoscroll = true
			c.outputView.Wrap = true
			c.outputView.Title = "Compiler output:"

			c.editorView.Editable = true
			c.editorView.Wrap = true
			c.editorView.Title = "Input:"

			c.editorView.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
				switch {
				case ch != 0 && mod == 0:
					v.EditWrite(ch)
				case key == gocui.KeySpace:
					v.EditWrite(' ')
				case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
					v.EditDelete(true)
				case key == gocui.KeyDelete:
					v.EditDelete(false)
				case key == gocui.KeyInsert:
					v.Overwrite = !v.Overwrite
				case key == gocui.KeyEnter:
					v.EditNewLine()
				case key == gocui.KeyArrowDown:
					v.MoveCursor(0, 1, false)
				case key == gocui.KeyArrowUp:
					v.MoveCursor(0, -1, false)
				case key == gocui.KeyArrowLeft:
					v.MoveCursor(-1, 0, false)
				case key == gocui.KeyArrowRight:
					v.MoveCursor(1, 0, false)
				}
				c.handleInputChange()
			})

			if _, err := g.SetCurrentView("editor"); err != nil {
				return err
			}
			//fmt.Fprintln(v, s)
		}
	}
	return nil
}

func (c *RunCompilerPipelineInteractiveCliInterface) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
