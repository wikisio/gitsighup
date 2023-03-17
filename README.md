# gitsighup
gitsighup run as a proxy, will git pull and send a HUP signal to the observed service process to reload configuration.

## Prerequisites

1. the observed service should reload the configuration by trapping the HUP signal. 
2. the target process should be managed by systemd. gitsighub sends signal with command `systemctl kill --signal=HUP {serviceName}`.
3. the configuration file should be managed in a git repository. gitsighup updates configuration by `cd configPath; git pull origin {tagOrBranchOrCommit}`


## API

  PUT /api/v1/services/{serviceName}?tag={tagOrBranchOrCommit}

## Configuraiton

```yaml
services:
- name: prometheus
  configPath: /opt/prometheus/conf
- name: alertmanager
  configPath: /opt/promethues/altermanger/conf
- name: gitsighup
  configPath: /opt/gitsigup/conf/gitsighup.yml
```

## startup

```ini
[Unit]
Description=Git Sig HUP Proxy
Wants=network-online.target
After=network-online.target


[Service]
User=gitsighup
Group=gitsighup
Type=simple
ExecStart=/opt/gitsigup/gitsigup -c /opt/gitsigup/conf/gitsighup.yml
Restart=on-failure

[Install]
WantedBy=default.target
```
  
