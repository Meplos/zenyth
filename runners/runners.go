package runners

import (
	"log"
	"os/exec"
	"path"
)

type Runner interface {
	Exec(string) ([]byte, error)
}

type BunRunner struct {
	Key string
}
type BaseRunner struct {
	Key string
}

func GetRunner(key string) Runner {
	switch key {
	case "bin":
		return NewBaseRunner()
	case "bun":
		return NewBunRunner()
	default:
		log.Fatalf("Invalid runner %v", key)
	}
	return nil
}

func NewBunRunner() *BunRunner {
	return &BunRunner{
		Key: "bun",
	}
}

func NewBaseRunner() *BaseRunner {
	return &BaseRunner{
		Key: "bin",
	}
}

func GetDir(file string) string {
	folder := path.Dir(file)
	return folder
}

func (r *BunRunner) Exec(file string) ([]byte, error) {
	dir := GetDir(file)
	cmd := exec.Command("bun", "run", file)
	cmd.Dir = dir
	b, err := cmd.Output()
	return b, err
}

func (r *BaseRunner) Exec(file string) ([]byte, error) {
	return exec.Command(file).Output()
}
