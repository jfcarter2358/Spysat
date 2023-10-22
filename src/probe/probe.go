package probe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"spysat/config"
	"spysat/logger"
	"strings"
	"time"

	"github.com/google/uuid"
)

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"spysat/auth"
// 	"spysat/cmd"
// 	"spysat/config"
// 	"spysat/constants"
// 	"spysat/container"
// 	"spysat/filestore"
// 	"spysat/health"
// 	"spysat/logger"
// 	"spysat/mongodb"
// 	"spysat/run"
// 	"spysat/state"
// 	"spysat/task"
// 	"strings"
// 	"time"
// )

type Probe struct {
	Language  string   `yaml:"language"`
	Run       string   `yaml:"run"`
	Arguments []string `yaml:"arguments"`
	Interval  int      `yaml:"interval"`
}

func Run(p Probe, args map[string]map[string]map[string]map[string]interface{}) {

	logger.Debugf("", "Probe params:")
	logger.Tracef("", "  language : %s", p.Language)
	logger.Tracef("", "  run      : %s", p.Run)
	logger.Tracef("", "  arguments: %v", p.Arguments)
	logger.Tracef("", "  interval : %v", p.Interval)
	logger.Debugf("", "Function args:")
	logger.Debugf("", "  args     : %v", args)
	id := uuid.New().String()
	logger.Debugf("", "Starting probe run with ID %s", id)

	runDir := fmt.Sprintf("/tmp/spysat/%s", id)
	err := os.MkdirAll(runDir, 0777)
	if err != nil {
		logger.Errorf("", "Error creating run directory %s", err.Error())
	}

	script_path := fmt.Sprintf("%s/run.%s", runDir, p.Language)

	script_command := []string{}
	script_contents := ""

	switch p.Language {
	case "sh":
		script_command = []string{"/bin/bash", fmt.Sprintf("%s/run.sh", runDir)}
		script_contents += "script_directory=\"$(dirname \"$(readlink -fm \"$0\")\")\"\n"
	case "py":
		script_command = []string{"python", fmt.Sprintf("%s/run.py", runDir)}
		script_contents += "import sys\n"
		script_contents += "import os\n"
		script_contents += "script_directory = os.path.dirname(os.path.abspath(sys.argv[0]))\n"
	default:
		logger.Errorf("", "Invalid language type: %s", p.Language)
	}

	script_contents += p.Run

	// Write out our run script
	script_data := []byte(script_contents)
	err = os.WriteFile(script_path, script_data, 0777)
	if err != nil {
		logger.Errorf("", "Error writing run file %s", err.Error())
	}

	for {
		for gName, g := range args {
			logger.Tracef("", "Checking group %s", gName)
			for oName, o := range g {
				logger.Tracef("", "Checking observer %s", oName)
				for sName, _ := range o {
					logger.Debugf("", "Doing probe for %s/%s/%s at /tmp/spysat/%s", gName, oName, sName, id)
					output, err := DoProbe(script_command, p.Arguments, args[gName][oName][sName], id)

					if err != nil {
						continue
					}

					jsonBody, err := json.Marshal(map[string]string{"data": output})
					if err != nil {
						continue
					}

					bodyReader := bytes.NewReader(jsonBody)
					requestURL := fmt.Sprintf("http://%s:%d/api/v1/basestation/%s/%s/%s", config.Config.Host, config.Config.Port, gName, oName, sName)
					req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
					if err != nil {
						continue
					}

					req.Header.Set("Content-Type", "application/json")
					client := http.Client{}
					client.Do(req)
				}
			}
		}
		logger.Tracef("", "Finished probing for id %s", id)
		time.Sleep(time.Duration(p.Interval) * time.Millisecond)
	}
}

func DoProbe(script_command, arguments []string, args map[string]interface{}, id string) (string, error) {
	for _, key := range arguments {
		val := args[key]
		script_command = append(script_command, fmt.Sprintf("%v", val))
	}
	logger.Tracef("", "Executing probe command %v", script_command)
	output, err := exec.Command(script_command[0], script_command[1:]...).CombinedOutput()
	outputString := strings.TrimSuffix(string(output), "\n")
	logger.Tracef("", "Probe output: %s", outputString)
	if err != nil {
		logger.Errorf("", "Error running script for probe /tmp/spysat/%s/: %s", id, err.Error())
		return "", err
	}
	return outputString, nil
}
