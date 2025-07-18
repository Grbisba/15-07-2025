# 15-07-2025
downloader

## config

**must be .json format**

```text
{
  "controller": {
    "host": string, // server host.
    "port": int, // server port.
    "read_timeout": int, // timeout to server read.
    "write_timeout": int, // timeout to server write.
    "idle_timeout": int, // idle server timeout.
    "shutdown_timeout": int, // timeout server shutdown.
    "cert_file": string, // server cert file if use TLS.
    "key_file": string, // server key file if use TLS.
    "use_tls": bool, // use TLS flag.
    "enabled": bool // should be ctrl enabled.
  },
  "loader": {
    "max_files_per_task": int, // Max files which task can handle.
    "max_concurrent_tasks": int, // Max task which server can handle.
    "zip_data_folder": string, // Path to zip data folder for save ".zip"-files (must be named like "zip").
    "zip_data_teardown": int, // Duration after which zip file will be cleaned.
    "temp_files_data_folder": string, // Path to temp file data folder for save ".pdf" of ".jpeg"-files (must be named like "temp").
    "temp_files_data_teardown": int, // Duration after which folder with temp files will be cleaned.
    "app_root_dir": string, // Path to "data" folder (must be named like "data").
  },
  "jwt": {
    "access_token_life_time": int, // access token life time.
    "refresh_token_life_time": int, // refresh token life time.
    "public_key_path": string, // path to public key path (make key-gen). 
    "private_key_path": string, // path to private key path (make key-gen). 
    "signing_algorithm": string, // signing algorithm.
  },
  "cache": {
    "login_user_cache_data_file": "data/cache/login-users-cache",
    "id_user_cache_data_file": "data/cache/id-users-cache",
    "task_cache_data_file": "data/cache/cache-data"
  }
}
```

## run
1)
envs:

    CONFIG_FILE_PATH - env to config file path from app root dir.

2)
```bash
make run
```

**if you see an error on server shutdown: create a folder in "data/cache"**