package config

type Cache struct {
	LoginUserCacheDataFile string `config:"login_user_cache_data_file" toml:"login_user_cache_data_file" yaml:"login_user_cache_data_file" json:"login_user_cache_data_file"`
	IDUserCacheDataFile    string `config:"id_user_cache_data_file" toml:"id_user_cache_data_file" yaml:"id_user_cache_data_file" json:"id_user_cache_data_file"`
	TaskCacheDataFile      string `config:"task_cache_data_file" toml:"task_cache_data_file" yaml:"task_cache_data_file" json:"task_cache_data_file"`
}

func (c *Cache) copy() *Cache {
	if c == nil {
		return nil
	}

	return &Cache{
		LoginUserCacheDataFile: c.LoginUserCacheDataFile,
		IDUserCacheDataFile:    c.IDUserCacheDataFile,
		TaskCacheDataFile:      c.TaskCacheDataFile,
	}
}
