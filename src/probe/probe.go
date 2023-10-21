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
	for {
		for gName, g := range args {
			for oName, o := range g {
				for sName, s := range o {
					output, err := DoProbe(p, s)

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
		time.Sleep(time.Duration(p.Interval) * time.Millisecond)
	}
}

func DoProbe(p Probe, args map[string]interface{}) (string, error) {
	id := uuid.New().String()
	logger.Debugf("", "Starting probe run with ID %s", id)

	runDir := fmt.Sprintf("/tmp/spysat/%s", id)
	err := os.MkdirAll(runDir, 0755)
	if err != nil {
		logger.Errorf("", "Error creating run directory %s", err.Error())
		return "", err
	}

	script_path := fmt.Sprintf("%s/run.%s", runDir, p.Language)

	script_command := []string{}
	script_contents := ""

	switch p.Language {
	case "sh":
		script_command = []string{"/bin/bash", "-c", fmt.Sprintf("%s/run.sh", runDir)}
		for _, name := range p.Arguments {
			script_command = append(script_command, fmt.Sprintf(" '%v'", args[name]))
		}
	case "py":
		script_command = []string{"python", fmt.Sprintf("%s/run.py", runDir)}
		for _, name := range p.Arguments {
			script_command = append(script_command, fmt.Sprintf(" '%v'", args[name]))
		}
	default:
		logger.Errorf("", "Invalid language type: %s", p.Language)
		return "", nil
	}

	script_contents += p.Run

	// Write out our run script
	script_data := []byte(script_contents)
	err = os.WriteFile(script_path, script_data, 0777)
	if err != nil {
		logger.Errorf("", "Error writing run file %s", err.Error())
		return "", err
	}

	output, err := exec.Command(script_command[0], script_command[1:]...).Output()
	if err != nil {
		logger.Errorf("", "Error running script: %s", err.Error())
		return "", err
	}

	return string(output), nil
}
