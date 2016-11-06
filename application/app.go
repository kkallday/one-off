package application

import (
	"flag"
	"fmt"
	"strings"
)

type App struct {
	oneOff oneOff
}

type oneOff interface {
	Run(oneOffInputs OneOffInputs) error
}

func New(oneOff oneOff) App {
	return App{
		oneOff: oneOff,
	}
}

func (a App) Execute(args []string) error {
	oneOffInputs, err := a.parseArgs(args)
	if err != nil {
		return err
	}

	err = a.validateInputs(oneOffInputs)
	if err != nil {
		return err
	}

	err = a.oneOff.Run(oneOffInputs)
	if err != nil {
		return err
	}

	return nil
}

func (App) parseArgs(args []string) (OneOffInputs, error) {
	var oneOffInputs OneOffInputs

	flags := flag.NewFlagSet("app", flag.ContinueOnError)
	flags.StringVar(&oneOffInputs.TargetAlias, "ta", "", "concourse target alias")
	flags.StringVar(&oneOffInputs.Pipeline, "p", "", "name of pipeline")
	flags.StringVar(&oneOffInputs.Job, "j", "", "name of job")
	flags.StringVar(&oneOffInputs.Task, "t", "", "name of task")
	flags.StringVar(&oneOffInputs.OutputDir, "out", "", "(optional) directory to write one off script")
	err := flags.Parse(args)
	if err != nil {
		return OneOffInputs{}, err
	}

	return oneOffInputs, nil
}

func (App) validateInputs(oneOffInputs OneOffInputs) error {
	var errs []string

	if oneOffInputs.TargetAlias == "" {
		errs = append(errs, "target alias -ta")
	}

	if oneOffInputs.Pipeline == "" {
		errs = append(errs, "pipeline -p")
	}

	if oneOffInputs.Job == "" {
		errs = append(errs, "job -j")
	}

	if oneOffInputs.Task == "" {
		errs = append(errs, "task -t")
	}

	if len(errs) > 0 {
		return fmt.Errorf("missing %s", strings.Join(errs, ", "))
	}

	return nil
}
