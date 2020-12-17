package context

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type SelectionBlock interface {
	Begin() (int, int)
	End() (int, int)
	Describe(src SelectionBlock, id int, mappingID func(SelectionBlock) int) []string
	GetID() int
}

type SelectionBlocks []SelectionBlock

func (s SelectionBlocks) Len() int {
	return len(s)
}
func (s SelectionBlocks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SelectionBlocks) Less(i, j int) bool {
	b1, e1 := s[i].Begin()
	b2, e2 := s[j].Begin()
	if b1 == b2 {
		return e1 < e2
	}
	return b1 < b2
}

func isInRange(line int, col int, block SelectionBlock) bool {
	bl, bc := block.Begin()
	el, ec := block.End()
	return (line > bl || (line == bl && col >= bc)) && (line < el || (line == el && col < ec))
}

func hashBlock(block SelectionBlock) string {
	bl, bc := block.Begin()
	el, ec := block.End()
	return fmt.Sprintf("%d|%d|%d|%d", bl, bc, el, ec)
}

func (c *ParsingContext) PrintSelectionBlocksList(inputBlocks SelectionBlocks) string {
	blocks := inputBlocks

	sort.Sort(blocks)
	blockIDs := map[string]int{}
	for _, block := range blocks {
		blockIDs[hashBlock(block)] = block.GetID()
	}

	output := []string{}
	for _, block := range blocks {
		blockID := blockIDs[hashBlock(block)]
		description := block.Describe(block, blockID, func(block SelectionBlock) int {
			id, ok := blockIDs[hashBlock(block)]
			if !ok {
				panic("Could not find ID for matching selection block.")
			}
			return id
		})
		bl, bc := block.Begin()
		el, ec := block.End()
		selDescription := fmt.Sprintf("%s %-20s %s %-20s",
			"starts at",
			fmt.Sprintf("(line %d, col %d)", bl, bc),
			"and ends at",
			fmt.Sprintf("(line %d, col %d)", el, ec))
		output = append(output, fmt.Sprintf("%-30s %s", fmt.Sprintf("Block: %3d(%s)", blockID, description), selDescription))
	}
	return fmt.Sprintf("    %s\n", strings.Join(output, "\n    "))
}

func (c *ParsingContext) PrintSelectionBlocks(inputBlocks SelectionBlocks) string {
	blocks := inputBlocks
	input := []byte{}
	input = append(input, c.ParserInput...)
	input = append(input, []byte("\n\n")...)
	scanner := bufio.NewScanner(bytes.NewReader(input))
	curLineNo := 1
	lines := []string{}

	//curBlock := 0
	sort.Sort(blocks)
	seenBlocksMap := map[int]*color.Attribute{}

	var curColBg *color.Attribute = nil
	colorIndex := -1
	colors := []color.Attribute{
		color.BgRed,
		color.BgCyan,
		color.BgBlue,
		color.BgGreen,
		color.BgYellow,
	}

	blockIDs := map[string]int{}
	for _, block := range blocks {
		blockIDs[hashBlock(block)] = block.GetID()
	}

	lineWidth := 60
	for scanner.Scan() {
		//letterBuf := []rune{}
		lineOut := ""
		//lineOutActualLen := 0
		blocksInfo := []string{}
		originalLine := scanner.Text()
		//lineWidthCutoffPos := -1
		//if len(originalLine) > lineWidth {
		//	originalLine = originalLine[:lineWidth]
		//}
		for i, letter := range originalLine {
			curColNo := i+1
			var curBlock SelectionBlock = nil
			for j := 0; j<len(blocks); j++ {
				if isInRange(curLineNo, curColNo, blocks[j]) {
					if curBlock == nil {
						curBlock = blocks[j]
					} else {
						cl, ce := curBlock.End()
						nl, ne := blocks[j].End()
						if cl > nl || (cl == nl && ce > ne) {
							curBlock = blocks[j]
						}
					}
				}
			}
			if curBlock != nil {
				blockID := blockIDs[hashBlock(curBlock)]
				wasBlockStart := false
				descriptionText := ""
				if _, ok := seenBlocksMap[blockID]; !ok {
					wasBlockStart = true
				}
				if wasBlockStart {
					colorIndex = (colorIndex+1) % len(colors)
					curColBg = &colors[colorIndex]
					seenBlocksMap[blockID] = curColBg
					formatBlockSelection := color.New(*curColBg).SprintFunc()
					description := curBlock.Describe(curBlock, blockID, func(block SelectionBlock) int {
						id, ok := blockIDs[hashBlock(block)]
						if !ok {
							panic("Could not find ID for matching selection block.")
						}
						return id
					})
					if len(description) > 0 {
						descriptionText = fmt.Sprintf("(%s)", strings.Join(description, "; "))
					}
					blocksInfo = append(blocksInfo, fmt.Sprintf("%s%s", formatBlockSelection(fmt.Sprintf("Block: %3d", blockID)), descriptionText))
				}
				formatBlockSelection := color.New(*curColBg).SprintFunc()
				lineOut = lineOut + formatBlockSelection(string(letter))
			} else {
				lineOut = lineOut + string(letter)
			}
		}
		linePostfix := ""
		if len(blocksInfo) > 0 {
			linePostfix = fmt.Sprintf("[%s]", strings.Join(blocksInfo, ", "))
		}
		lineContent := fmt.Sprintf("%5d | %s", curLineNo, lineOut)
		if len(linePostfix) > 0 {
			if lineWidth-len(originalLine) > 0 {
				lineContent += strings.Repeat(" ", lineWidth-len(originalLine))
			}
		}
		lines = append(lines, fmt.Sprintf("%s%s", lineContent, linePostfix))
		//for i, letter := range originalLine {
		//	isBlockStart := false
		//	curColNo := i+1
		//	if isInRange(curLineNo, curColNo, blocks[curBlock]) {
		//		if !wasInCurBlock {
		//			isBlockStart = true
		//			colorIndex = (colorIndex+1) % len(colors)
		//			curColBg = &colors[colorIndex]
		//			lineOut = lineOut + string(letterBuf)
		//			lineOutActualLen += len(letterBuf)
		//			letterBuf = []rune{}
		//		}
		//		wasInCurBlock = true
		//	} else if(wasInCurBlock) {
		//		wasInCurBlock = false
		//		if curColBg == nil {
		//			lineOut = lineOut + string(letterBuf)
		//			lineOutActualLen += len(letterBuf)
		//		} else {
		//			formatBlockSelection := color.New(*curColBg).SprintFunc()
		//			lineOut = lineOut + formatBlockSelection(string(letterBuf))
		//			lineOutActualLen += len(letterBuf)
		//		}
		//		letterBuf = []rune{}
		//		curColBg = nil
		//		curBlock++
		//		if isInRange(curLineNo, curColNo, blocks[curBlock]) {
		//			isBlockStart = true
		//			colorIndex = (colorIndex+1) % len(colors)
		//			curColBg = &colors[colorIndex]
		//			wasInCurBlock = true
		//		}
		//	} else {
		//		// Do nothing
		//	}
		//	letterBuf = append(letterBuf, letter)
		//	if isBlockStart {
		//		formatBlockSelection := color.New(*curColBg).SprintFunc()
		//		blockID := blockIDs[hashBlock(blocks[curBlock])]
		//		descriptionText := ""
		//		description := blocks[curBlock].Describe(blocks[curBlock], blockID, func(block SelectionBlock) int {
		//			id, ok := blockIDs[hashBlock(block)]
		//			if !ok {
		//				panic("Could not find ID for matching selection block.")
		//			}
		//			return id
		//		})
		//		if len(description) > 0 {
		//			descriptionText = fmt.Sprintf("(%s)", strings.Join(description, "; "))
		//		}
		//		blocksInfo = append(blocksInfo, fmt.Sprintf("%s%s", formatBlockSelection(fmt.Sprintf("Block: %3d", blockID)), descriptionText))
		//	}
		//	if i+1 == lineWidth {
		//		offset := 0
		//		if curColBg != nil {
		//			formatBlockSelection := color.New(*curColBg).SprintFunc()
		//			offset = len(formatBlockSelection(string(letterBuf)))
		//		} else {
		//			offset = len(string(letterBuf))
		//		}
		//		lineWidthCutoffPos = len(string(letterBuf)) + offset
		//	}
		//}
		//if curColBg != nil {
		//	formatBlockSelection := color.New(*curColBg).SprintFunc()
		//	lineOut = lineOut + formatBlockSelection(string(letterBuf))
		//	lineOutActualLen += len(letterBuf)
		//} else {
		//	lineOut = lineOut + string(letterBuf)
		//	lineOutActualLen += len(letterBuf)
		//}
		//linePostfix := ""
		//if len(blocksInfo) > 0 {
		//	linePostfix = fmt.Sprintf("[%s]", strings.Join(blocksInfo, ", "))
		//}
		//lineContent := fmt.Sprintf("%5d | %s", curLineNo, lineOut)
		//if len(linePostfix) > 0 {
		//	if lineWidthCutoffPos != -1 {
		//		lineContent = lineContent[:lineWidthCutoffPos]
		//	} else {
		//		lineContent += strings.Repeat(" ", lineWidth-len(originalLine))
		//	}
		//}
		//lines = append(lines, fmt.Sprintf("%s%s", lineContent, linePostfix))
		curLineNo++
	}
	return strings.Join(lines, "\n")
}
