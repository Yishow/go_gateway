package main

import (
	"context"
	"log"

	"go-gateway/internal/datalink/modbusshare"
	datalinksettings "go-gateway/internal/datalink/settings"
)

func loadPersistedShareSettings(ctx context.Context, svc *modbusshare.Service, repo datalinksettings.Repository) bool {
	if err := svc.SetSettingsRepository(repo); err != nil {
		log.Printf("本機 Modbus 分享服務設定儲存庫初始化失敗: %v", err)
		return false
	}
	if err := svc.LoadPersistedSettings(ctx); err != nil {
		log.Printf("本機 Modbus 分享服務持久化設定還原失敗，hydration 維持 pending: %v", err)
		return false
	}
	return true
}

func startConfiguredShareListener(ctx context.Context, svc *modbusshare.Service) {
	if !svc.Settings().Enabled {
		return
	}
	if err := svc.StartConfiguredListener(ctx); err != nil {
		log.Printf("本機 Modbus 分享服務啟動失敗，listener 維持停止: %v", err)
		return
	}
}
