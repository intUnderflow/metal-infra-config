package bootstrap

import (
	"encoding/json"
	"github.com/go-git/go-billy/v5"
	"io"
)

func GetBootstrapRecords(file billy.File) (map[string]string, error) {
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	records := map[string]string{}
	err = json.Unmarshal(contents, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
