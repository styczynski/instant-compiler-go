package compiler_pipeline

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type RunCompilerPipelineInteractiveCliInterface struct {
	compilerPipeline *HeadlessCompilerPipeline

	errView    *gocui.View
	editorView *gocui.View
	outputView *gocui.View
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
		c.outputView.Clear()
		fmt.Fprintf(c.outputView, response.Summary)
	})
}

func (c *RunCompilerPipelineInteractiveCliInterface) Run() {
	c.compilerPipeline = CreateHeadlessCompilerPipeline()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(c.layout)
	g.InputEsc = true
	g.Cursor = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (c *RunCompilerPipelineInteractiveCliInterface) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	var err error
	if c.editorView, err = g.SetView("editor", 0, 0, maxX/2, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if c.outputView, err = g.SetView("output", maxX/2+1, 0, maxX-1, maxY-5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		if c.errView, err = g.SetView("error-view", 0, maxY-4, maxX/2, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			c.editorView.Editable = true
			c.editorView.Wrap = true
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
