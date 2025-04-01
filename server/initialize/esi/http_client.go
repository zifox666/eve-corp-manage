package esi

import (
	"eve-corp-manager/config"
	"eve-corp-manager/core/esi"
)

// InitESIClient 初始化ESI HTTP客户端
func InitESIClient() {
	proxyHost := config.AppConfig.Proxy.Host
	proxyPort := config.AppConfig.Proxy.Port

	userAgent := "EveCorpManage/1.0.0 (zifox666@gmail.com; +https://github.com/zifox666/eve-corp-manage)"

	apiType := "esi"

	esi.EsiClient = esi.NewClient(proxyHost, proxyPort, userAgent, apiType)
}

// InitJaniceClient 初始化Janice HTTP客户端
func InitJaniceClient() {
	proxyHost := config.AppConfig.Proxy.Host
	proxyPort := config.AppConfig.Proxy.Port

	userAgent := "EveCorpManage/1.0.0 (zifox666@gmail.com; +https://github.com/zifox666/eve-corp-manage)"

	apiType := "janice"

	esi.JaniceClient = esi.NewClient(proxyHost, proxyPort, userAgent, apiType)
}
