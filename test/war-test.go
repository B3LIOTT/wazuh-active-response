package test

import (
	"fmt"
	"os"
)

func Add(keys []interface{}) error {
	file, err := os.OpenFile("ar-test-result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("Active response triggered by rule ID: <%v>\n", keys)); err != nil {
		return err
	}

	return nil
}

func Delete() error {
	if err := os.Remove("ar-test-result.txt"); err != nil {
		return err
	}

	return nil
}
