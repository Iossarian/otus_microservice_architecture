package createorder

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Workflow struct {
	Steps map[string]SagaStep
}

type SagaStep struct {
	Transaction  Action
	Compensation Action
}

type Action func(ctx echo.Context, request Request) error

func NewWorkflow() *Workflow {
	return &Workflow{
		Steps: make(map[string]SagaStep),
	}
}

func (w *Workflow) AddStep(name string, step SagaStep) {
	w.Steps[name] = step
}

func (w *Workflow) Execute(ctx echo.Context, request Request) error {
	for stepName, stepAction := range w.Steps {
		fmt.Printf("Executing step %s\n", stepName)

		if err := stepAction.Transaction(ctx, request); err != nil {
			fmt.Printf("Transaction failed: %s\n", err)

			for n, a := range w.Steps {
				fmt.Printf("Compensating step %s\n", n)

				if err := a.Compensation(ctx, request); err != nil {
					return errors.Wrapf(err, "compensating step %s", n)
				}

				if n == stepName {
					break
				}
			}
		}
	}

	return nil
}
