package main

import (
	"os"
)

// pr-pal @ https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang
func ensureDir(dirName string) error {

	err := os.Mkdir(dirName, 02666) //probably umasked with 22 not 02

	os.Chmod(dirName, 02666)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}

}
