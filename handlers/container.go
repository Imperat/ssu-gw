package handlers

import (
	"net/http"
	"strconv"

	"github.com/dekarti/ssu-gw/common"
	"github.com/dekarti/ssu-gw/models"
	"github.com/dekarti/ssu-gw/util"
	"github.com/labstack/echo"
)

type (
	TestUnit struct {
		Input    string `json:"input"`
		Output   string `json:"output"`
		Result   string `json:"result"`
		Expected string `json:"expected"`
	}
)

// POST /task/:id/work/launch
func LaunchWorkHandler(c echo.Context) error {
	taskNumber, _ := strconv.Atoi(c.Param("id"))

	work := &models.Work{
		DockerClient: common.CLI,
		Task:         models.Tasks[taskNumber],
		Path:         "/Users/aarutyunyan/Desktop/ssu-fl/tasks",
		Command:      []string{"python", "exp_nfa.py"},
	}

	testUnit := &TestUnit{}

	if err := c.Bind(testUnit); err != nil {
		return err
	}

	if testUnit.Input == "" {
		testUnit.Input = models.Tasks[taskNumber].DefaultInput
		testUnit.Expected = models.Tasks[taskNumber].ExpectedOutput
	}

	if err := work.WriteInput(testUnit.Input); err != nil {
		return err
	}

	if err := work.LaunchWork(); err != nil {
		return err
	}

	if err := util.IsContainerSuccessfullyExited(work.ContainerName, 5); err != nil {
		return err
	}

	if output, err := work.ReadOutput(); err != nil {
		return err
	} else {
		testUnit.Output = output
		if testUnit.Expected != "" {
			if output == testUnit.Expected {
				testUnit.Result = "Output matches expected value"
				return c.JSON(http.StatusOK, testUnit)
			} else {
				testUnit.Result = "Output doesn't match expected value"
				return c.JSON(http.StatusBadRequest, testUnit)
			}
		}
		return c.JSON(http.StatusOK, testUnit)
	}
}
