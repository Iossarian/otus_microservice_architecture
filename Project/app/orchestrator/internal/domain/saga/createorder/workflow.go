package createorder

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Workflow struct {
	Steps []SagaStep
}

type SagaStep struct {
	Name         string
	Transaction  Action
	Compensation Action
}

type Action func(ctx echo.Context, request Request) error

func NewWorkflow() *Workflow {
	return &Workflow{
		Steps: make([]SagaStep, 0),
	}
}

func (w *Workflow) AddStep(step SagaStep) {
	w.Steps = append(w.Steps, step)
}

func (w *Workflow) Execute(ctx echo.Context, request Request) error {
	for _, step := range w.Steps {
		fmt.Printf("Executing step %s\n", step.Name)

		if err := step.Transaction(ctx, request); err != nil {
			fmt.Printf("Transaction step %s failed: %s\n", step.Name, err)

			for _, s := range w.Steps {
				fmt.Println("stepName", step.Name)
				if s.Name == step.Name {
					return errors.Wrapf(err, "executing step %s", s.Name)
				}

				fmt.Printf("Compensating step %s\n", s.Name)

				if err := s.Compensation(ctx, request); err != nil {
					return errors.Wrapf(err, "compensating step %s", s.Name)
				}
			}

			return errors.Wrapf(err, "executing step %s", step.Name)
		}
	}

	return nil
}
