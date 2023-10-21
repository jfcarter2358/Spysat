package basestation

import (
	"spysat/analyst"
	"spysat/datastore"
	"spysat/operation"
)

func ReceiveData(data, name, observer, group string) error {
	s := operation.Operations.Observers[observer].Streams[name]
	a := operation.Operations.Analysts[s.Analyst]
	output, err := analyst.DoAnalyze(a, name, group, observer, data, s.Arguments)
	if err != nil {
		return err
	}
	datastore.AddData(output, name, group, observer)
	return nil
}

func GetData(name, observer, group string) map[string][]interface{} {
	data := datastore.Data[group][observer][name]
	return data
}
