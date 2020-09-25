package env

import (
	"github.com/leastauthority/fleece/bindata"
)

func RestoreBindata(dir string) error {
	if err := bindata.RestoreAssets(dir, "docker"); err != nil {
		return err
	}

	if err := bindata.RestoreAssets(dir, "fuzzing"); err != nil {
		return err
	}
	return nil
}
