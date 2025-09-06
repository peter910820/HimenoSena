# HimenoSena
A disocrd bot is used to intercept messages from other bots to designated channels

## 專案架構:
```bash
.
├── bot      # 機器人核心處理模組 包含初始化資料庫連線以及註冊斜線指令等
│   └── ...
├── commands # 指令模組(一個指令一個檔案)
│   └── ...
├── handlers # 機器人事件handlers模組
│   └── ...
├── models   # 型別定義模組
│   └── ...
├── utils    # 通用函式模組
│   └── ...
└── main.go  # 機器人進入點
```