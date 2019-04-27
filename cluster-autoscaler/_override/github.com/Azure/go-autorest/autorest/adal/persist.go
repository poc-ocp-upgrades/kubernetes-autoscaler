package adal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadToken(path string) (*Token, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%s) while loading token: %v", path, err)
	}
	defer file.Close()
	var token Token
	dec := json.NewDecoder(file)
	if err = dec.Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode contents of file (%s) into Token representation: %v", path, err)
	}
	return &token, nil
}
func SaveToken(path string, mode os.FileMode, token Token) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory (%s) to store token in: %v", dir, err)
	}
	newFile, err := ioutil.TempFile(dir, "token")
	if err != nil {
		return fmt.Errorf("failed to create the temp file to write the token: %v", err)
	}
	tempPath := newFile.Name()
	if err := json.NewEncoder(newFile).Encode(token); err != nil {
		return fmt.Errorf("failed to encode token to file (%s) while saving token: %v", tempPath, err)
	}
	if err := newFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file %s: %v", tempPath, err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to move temporary token to desired output location. src=%s dst=%s: %v", tempPath, path, err)
	}
	if err := os.Chmod(path, mode); err != nil {
		return fmt.Errorf("failed to chmod the token file %s: %v", path, err)
	}
	return nil
}
