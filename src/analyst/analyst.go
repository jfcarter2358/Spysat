package analyst

import (
	"fmt"
	"os"
	"os/exec"
	"spysat/logger"

	"github.com/google/uuid"
)

type Analyst struct {
	Language string
	Run      string
}

func DoAnalyze(a Analyst, stream, group, observer, input_contents string, args map[string]interface{}) (string, error) {
	id := uuid.New().String()
	logger.Debugf("", "Starting analyst run for %s/%s with ID %s", group, observer, id)

	runDir := fmt.Sprintf("/tmp/spysat/%s", id)
	err := os.MkdirAll(runDir, 0755)
	if err != nil {
		logger.Errorf("", "Error creating run directory %s", err.Error())
		return "", err
	}

	script_path := fmt.Sprintf("%s/run.%s", runDir, a.Language)
	input_path := fmt.Sprintf("%s/input.txt", runDir)

	script_command := []string{}
	script_contents := ""

	switch a.Language {
	case "sh":
		script_command = []string{"/bin/bash", "-c", fmt.Sprintf("%s/run.sh", runDir)}

		script_contents += "script_directory=$( cd -- \"$( dirname -- \"${BASH_SOURCE[0]}\" )\" &> /dev/null && pwd )\n"
		script_contents += "data_in=\"$(cat \"${script_directory}/input.txt\")\"\n"
	case "py":
		script_command = []string{"python", fmt.Sprintf("%s/run.py", runDir)}

		script_contents += "import sys\n"
		script_contents += "import os\n"
		script_contents += "script_directory = os.path.dirname(os.path.abspath(sys.argv[0]))\n"
		script_contents += "with open(f'{script_directory}/input.txt, 'r', encoding='utf-8') as data_file:\n"
		script_contents += "    data_in = data_file.read()\n"
	default:
		logger.Errorf("", "Invalid language type: %s", a.Language)
		return "", nil
	}

	script_contents += a.Run

	// Write out our run script
	script_data := []byte(script_contents)
	err = os.WriteFile(script_path, script_data, 0777)
	if err != nil {
		logger.Errorf("", "Error writing run file %s", err.Error())
		return "", err
	}

	// Write out our input script
	input_data := []byte(input_contents)
	err = os.WriteFile(input_path, input_data, 0777)
	if err != nil {
		logger.Errorf("", "Error writing input file %s", err.Error())
		return "", err
	}

	output, err := exec.Command(script_command[0], script_command[1:]...).Output()
	if err != nil {
		logger.Errorf("", "Error running script: %s", err.Error())
		return "", err
	}

	// datastore.AddData(string(output), stream, group, observer)

	return string(output), nil
}
