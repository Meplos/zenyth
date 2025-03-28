package config

import (
	"encoding/json"

	"github.com/Meplos/zenyth/tasks"
)

func ParseTaskDef(input string) []tasks.TaskDef {
	var tasks []tasks.TaskDef

	if err := json.Unmarshal([]byte(input), &tasks); err != nil {
		panic("Cannot unmarshall file")
	}
	return tasks
}
