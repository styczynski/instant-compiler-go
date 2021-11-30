package events_collector

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/styczynski/latte-compiler/src/config"
)

var formatStatusBg = color.New(color.BgBlack).SprintFunc()
var formatStatusFg = color.New(color.FgHiBlue).SprintFunc()

type StatusUpdater interface {
	UpdateStatus(description string)
	Init()
	Done()
}

type SilentStatusUpdater struct{}

type CliProgressBarStatusUpdater struct {
	spinner *spinner.Spinner
}

func CreateStatusUpdater(isSilent bool) StatusUpdater {
	if isSilent {
		return SilentStatusUpdater{}
	}
	return &CliProgressBarStatusUpdater{}
}

func init() {
	config.RegisterEntityFactory(config.ENTITY_STATUS_UPDATER, CliProgressBarStatusUpdaterFactory{})
	config.RegisterEntityFactory(config.ENTITY_STATUS_UPDATER, SilentStatusUpdaterFactory{})
}

type SilentStatusUpdaterFactory struct{}

func (SilentStatusUpdaterFactory) CreateEntity(c config.EntityConfig) interface{} {
	return SilentStatusUpdater{}
}

func (SilentStatusUpdaterFactory) Params(argSpec *config.EntityArgSpec) {
}

func (SilentStatusUpdaterFactory) EntityName() string {
	return "updater-silent"
}

type CliProgressBarStatusUpdaterFactory struct{}

func (CliProgressBarStatusUpdaterFactory) CreateEntity(c config.EntityConfig) interface{} {
	return &CliProgressBarStatusUpdater{}
}

func (CliProgressBarStatusUpdaterFactory) Params(argSpec *config.EntityArgSpec) {
}

func (CliProgressBarStatusUpdaterFactory) EntityName() string {
	return "updater-cli-progress"
}

func (SilentStatusUpdater) Init() {
	// No-op
}

func (SilentStatusUpdater) Done() {
	// No-op
}

func (SilentStatusUpdater) UpdateStatus(description string) {
	// No-op
}

func (u *CliProgressBarStatusUpdater) Init() {
	u.spinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	u.spinner.Start()
}

func (u *CliProgressBarStatusUpdater) Done() {
	u.spinner.Stop()
}

func (u *CliProgressBarStatusUpdater) UpdateStatus(description string) {
	u.spinner.Prefix = fmt.Sprintf("\033[36m%s\033[m  ", formatStatusFg(formatStatusBg((description))))
}
