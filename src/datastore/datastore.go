package datastore

import (
	"spysat/logger"
	"spysat/observer"
	"strconv"
	"strings"
)

// group/observer/stream/headers/data
var Data map[string]map[string]map[string]map[string][]interface{}
var Headers map[string]map[string]map[string][]string
var Types map[string]map[string]map[string][]string

func Init() {
	Data = make(map[string]map[string]map[string]map[string][]interface{})
	Headers = make(map[string]map[string]map[string][]string)
	Types = make(map[string]map[string]map[string][]string)
}

func AddObserver(observer string, o observer.Observer) {
	if _, ok := Data[o.Group]; !ok {
		Data[o.Group] = make(map[string]map[string]map[string][]interface{})
		Headers[o.Group] = make(map[string]map[string][]string)
		Types[o.Group] = make(map[string]map[string][]string)
	}

	Headers[o.Group][observer] = make(map[string][]string)
	Types[o.Group][observer] = make(map[string][]string)

	observer_data := make(map[string]map[string][]interface{})
	for name, s := range o.Streams {
		stream_data := make(map[string][]interface{})
		for _, header := range s.Headers {
			stream_data[header] = make([]interface{}, 0)
		}
		observer_data[name] = stream_data
		Headers[o.Group][observer][name] = s.Headers
		Types[o.Group][observer][name] = s.Types
	}
	Data[o.Group][observer] = observer_data
}

func AddData(csv, stream, group, observer string) {
	if len(csv) == 0 {
		return
	}
	lines := strings.Split(csv, "\n")
	if len(lines) == 0 {
		return
	}

	headers := Headers[group][observer][stream]
	types := Types[group][observer][stream]

	for _, line := range lines {
		parts := strings.Split(line, ",")
		for idx, val := range parts {
			switch types[idx] {
			case "int":
				conv, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				Data[group][observer][stream][headers[idx]] = append(Data[group][observer][stream][headers[idx]], conv)
				break
			case "float":
				conv, err := strconv.ParseFloat(val, 64)
				if err != nil {
					continue
				}
				Data[group][observer][stream][headers[idx]] = append(Data[group][observer][stream][headers[idx]], conv)
				break
			case "string":
				Data[group][observer][stream][headers[idx]] = append(Data[group][observer][stream][headers[idx]], val)
				break
			case "bool":
				conv, err := strconv.ParseBool(val)
				if err != nil {
					continue
				}
				Data[group][observer][stream][headers[idx]] = append(Data[group][observer][stream][headers[idx]], conv)
				break
			default:
				logger.Errorf("", "Invalid variable type: %s", types[idx])
			}
		}
	}
}
