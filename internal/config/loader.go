package config

import (
	"time"
)

type Loader struct {
	MaxFilesPerTask       int    `config:"max_files_per_task" json:"max_files_per_task" toml:"max_files_per_task" yaml:"max_files_per_task"`
	MaxConcurrentTasks    int    `config:"max_concurrent_tasks" json:"max_concurrent_tasks" toml:"max_concurrent_tasks" yaml:"max_concurrent_tasks"`
	ZipDataFolder         string `config:"zip_data_folder" json:"zip_data_folder" toml:"zip_data_folder" yaml:"zip_data_folder"`
	ZipDataTeardown       int    `config:"zip_data_teardown" json:"zip_data_teardown" toml:"zip_data_teardown" yaml:"zip_data_teardown"`
	TempFilesDataFolder   string `config:"temp_files_data_folder" json:"temp_files_data_folder" toml:"temp_files_data_folder" yaml:"temp_files_data_folder"`
	TempFilesDataTeardown int    `config:"temp_files_data_teardown" json:"temp_files_data_teardown" toml:"temp_files_data_teardown" yaml:"temp_files_data_teardown"`
	AppRootDir            string `config:"app_root_dir" json:"app_root_dir" toml:"app_root_dir" yaml:"app_root_dir"`
}

func (c *Loader) ZipDataTeardownSeconds() time.Duration {
	if c == nil {
		return -1
	}

	return time.Second * time.Duration(c.ZipDataTeardown)
}

func (c *Loader) TempFilesDataTeardownSeconds() time.Duration {
	if c == nil {
		return -1
	}

	return time.Second * time.Duration(c.TempFilesDataTeardown)
}

func (c *Loader) copy() *Loader {
	if c == nil {
		return nil
	}

	return &Loader{
		MaxFilesPerTask:       c.MaxFilesPerTask,
		MaxConcurrentTasks:    c.MaxConcurrentTasks,
		ZipDataFolder:         c.ZipDataFolder,
		ZipDataTeardown:       c.ZipDataTeardown,
		TempFilesDataFolder:   c.TempFilesDataFolder,
		TempFilesDataTeardown: c.TempFilesDataTeardown,
	}
}
