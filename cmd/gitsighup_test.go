package main

import "testing"

func TestSendRequest(t *testing.T) {
	url := "http://192.168.0.143:3000/api/v1/configsrv/"
	SendRequest(url, "dev", "alertmanager", "alertmanager.yml", "/opt/prometheus/alertmanager/conf/")
}
