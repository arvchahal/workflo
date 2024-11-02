package githubactions

import (
	"os"
)

func directoryExists() (bool, error) {
	path := "./github/workflows"
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}
