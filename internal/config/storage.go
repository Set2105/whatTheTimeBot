package config

import (
	"os"
	"path/filepath"
	"sync"
)

type ConfigStorage struct {
	mx            sync.Mutex
	configDirPath string
}

func InitConfigStorage(envName, defaultPath string) (csPointer *ConfigStorage, err error) {
	cs := ConfigStorage{}
	cs.configDirPath, err = createConfigDirPath(envName, defaultPath)
	if err != nil {
		return nil, err
	}
	if err = createDirs(cs.configDirPath); err != nil {
		return nil, err
	}
	return &cs, nil
}

func createDirs(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func createConfigDirPath(envName, defaultPath string) (string, error) {
	path := os.Getenv(envName)
	if path == "" {
		path = defaultPath
	}
	if !filepath.IsAbs(path) {
		abspath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(abspath, path)
	}
	return path, nil
}

func (cs *ConfigStorage) Save(key string, data []byte) error {
	filePath := filepath.Join(cs.configDirPath, key)
	cs.mx.Lock()
	defer cs.mx.Unlock()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ConfigStorage) Read(key string) ([]byte, error) {
	filePath := filepath.Join(cs.configDirPath, key)
	cs.mx.Lock()
	data, err := os.ReadFile(filePath)
	cs.mx.Unlock()
	if err != nil {
		return nil, err
	}
	return data, nil
}
