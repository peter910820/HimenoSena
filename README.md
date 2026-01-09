# HimenoSena

多功能 discord 管理機器人(自用)

## 專案架構:

```bash
.
├── cmd/           # 主程式進入點
│   └── main.go    # 機器人初始化與啟動
├── bot/           # 機器人核心處理模組
│   ├── register.go      # 註冊斜線指令
│   └── exp_feature.go   # 經驗值相關功能
├── commands/      # 指令模組 (一個指令一個檔案)
│   ├── ping.go          # 心跳檢測指令
│   ├── send.go          # 發送訊息指令
│   ├── get_roles.go     # 取得身分組指令
│   ├── get_level.go     # 取得等級指令
│   ├── get_all_level.go # 取得群組等級排行指令
│   └── get_log.go       # 查詢 log 指令
├── handlers/      # 機器人事件處理模組
│   ├── event.go         # Discord 事件處理器
│   ├── interaction.go   # 互動事件處理器
│   └── common.go        # 通用處理函式(Ready等事件)
├── utils/         # 通用函式模組
│   └── reuse.go         # 可重用工具函式
└── model.go       # 型別定義
```
