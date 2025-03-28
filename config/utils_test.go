package config

import (
	"testing"
)

func TestReturnAnArrayOfTaskDef(t *testing.T) {
	const input = `[{
		"name": "cron-test-1",
		"exec": "/bin/ls",
		"cron": "* * * * *"
		}]`

	result := ParseTaskDef(input)
	if len(result) != 1 {
		t.Errorf("Result should have len of 1")
	}
	task := result[0]

	if task.Name != "cron-test-1" {
		t.Errorf("Taks should have name : cron-test-1")
	}

	if task.Exec != "/bin/ls" {
		t.Errorf("Taks should have exec: /bin/ls")
	}

	if task.Cron != "* * * * *" {
		t.Errorf("Taks should have exec: * * * * *")
	}
}
