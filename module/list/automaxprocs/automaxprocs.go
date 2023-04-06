package automaxprocs

import (
	"fmt"
	"github.com/vvfock3r/gooey/module/list/logger"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

type AutoMaxProcs struct{}

func (p *AutoMaxProcs) Register(*cobra.Command) {}

func (p *AutoMaxProcs) MustCheck(*cobra.Command) {}

func (p *AutoMaxProcs) Initialize(*cobra.Command) error {
	_, err := maxprocs.Set(maxprocs.Logger(p.logFunc))
	return err
}

func (p *AutoMaxProcs) logFunc(format string, v ...any) {
	logger.Info(fmt.Sprintf(format, v))
}
