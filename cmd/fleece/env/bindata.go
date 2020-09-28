package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leastauthority/fleece/bindata"
)

func promptOverwrite(files []string) (bool, error) {
	fmt.Println("The following files will be overwritten (stash changes first):")
	for _, file := range files {
		fmt.Println("\t" + file)
	}

	fmt.Print("Continue? [y/(n)]: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	// NB: input ends with '\n'
	switch input[:len(input)-1] {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

func RestoreBindata(dir string) error {
	overwrites, err := collectOverwrites(dir)
	if err != nil {
		return err
	}

	if len(overwrites) != 0 {
		shouldOverwrite, err := promptOverwrite(overwrites)
		if err != nil {
			return err
		}
		if !shouldOverwrite {
			return nil
		}
	}

	if err := bindata.RestoreAssets(dir, "docker"); err != nil {
		return err
	}
	return nil
}

func collectOverwrites(dir string) (overwrites []string, err error) {
	for _, name := range bindata.AssetNames() {
		assetInfo, err := bindata.AssetInfo(name)
		if err != nil {
			return nil, err
		}
		if assetInfo.IsDir() {
			continue
		}

		outputPath := filepath.Join(dir, filepath.FromSlash(name))
		_, err = os.Stat(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		overwrites = append(overwrites, assetInfo.Name())
	}
	return overwrites, nil
}
