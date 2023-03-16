# gitsighup
gitsighup run as a proxy, will git pull and send a HUP signal to the observed process to reload configuration.

## Configuraiton

services:
- name: prometheus
  configPath: /opt/prometheus/conf
- name: alertmanager
  configPath: /opt/promethues/altermanger/conf
  
