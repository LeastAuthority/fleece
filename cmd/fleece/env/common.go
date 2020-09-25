package env

import (
	"bufio"
	"fmt"
	"os"

	"github.com/leastauthority/fleece/bindata"
)

func RestoreBindata(dir string) error {
	overwrites, err := collectOverwrites(dir)
	if err != nil {
		return err
	}

	fmt.Println("The following files will be overwritten:")
	for _, file := range overwrites {
		fmt.Println("\t" + file)
	}
	fmt.Print("Continue? [y/(n)]: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// NB: input ends with '\n'
	switch input[:len(input)-1] {
	case "y", "yes":
	default:
		return nil
	}

	if err := bindata.RestoreAssets(dir, "docker"); err != nil {
		return err
	}

	if err := bindata.RestoreAssets(dir, "fuzzing"); err != nil {
		return err
	}
	return nil
}

func collectOverwrites(dir string) (overwrites []string, err error) {
	for _, name := range bindata.AssetNames() {
		info, err := bindata.AssetInfo(name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		if !info.IsDir() {
			overwrites = append(overwrites, info.Name())
		}
	}
	return overwrites, nil
}