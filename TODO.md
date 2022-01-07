
# features

## install script

```bash
LATEST_DEVTOOLS=`curl -s https://api.github.com/repos/MrL0co/devtools/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")'`
curl -o- https://raw.githubusercontent.com/MrL0co/devtools/${LATEST_DEVTOOLS}/install.sh
```

- make install.sh 
  - downloads latest version of tool (check current major version)
  - prompt install directory
  - run `devtools --initial-setup` for initial configuration
- `--initial-setup`
  - check existing config
    - if exists check update of config format
    - if not create new and run initial `question setup`
- `question setup`
  - ask questions about user preferences
    1. alias names?
    2. 
