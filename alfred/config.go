package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	Delimiter = "DELIMITER"
)

func GetDelimiter(wf *aw.Workflow) string {
	return wf.Config.Get(Delimiter)
}
