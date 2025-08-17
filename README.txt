███╗   ███╗ ██████╗ ███╗   ██╗ ██████╗ ██╗     ██╗████████╗██╗  ██╗
████╗ ████║██╔═══██╗████╗  ██║██╔═══██╗██║     ██║╚══██╔══╝██║  ██║
██╔████╔██║██║   ██║██╔██╗ ██║██║   ██║██║     ██║   ██║   ███████║
██║╚██╔╝██║██║   ██║██║╚██╗██║██║   ██║██║     ██║   ██║   ██╔══██║
██║ ╚═╝ ██║╚██████╔╝██║ ╚████║╚██████╔╝███████╗██║   ██║   ██║  ██║
╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚═╝  ╚═╝

  ░▒▓█ _LUCH_ █▓▒░
  A Telegram bot that brings remote control to your monolith lighting.

  ───────────────────────────────────────────────────────────────
  ▓ OVERVIEW
  **LUCH** connects to a WebSocket server and exposes Telegram commands
  for managing lighting effects and notification preferences.
  Simple. Reliable. Radiant.

  ───────────────────────────────────────────────────────────────
  ▓ REQUIREMENTS
  ▪ Go 1.20+
  ▪ `TELEGRAM_TOKEN` environment variable
  ▪ WebSocket server at `ws://localhost:8092`
  ▪ `notify.json` initialized with `[]`

  ───────────────────────────────────────────────────────────────
  ▓ FEATURES
  ▪ Slash commands for setup and help
  ▪ Reply keyboard: `Lamp On`, `Lamp Off`, `Led Off`, `Next Effect`
  ▪ Persistent notification subscriptions
  ▪ Auto-reconnect on WebSocket drop

  ───────────────────────────────────────────────────────────────
  ▓ RUNNING
  Create a `.env` with `TELEGRAM_TOKEN=<your_token>` and launch:

  ```sh
  go run cmd/luch/main.go
  ```

  ───────────────────────────────────────────────────────────────
  ▓ FINAL WORDS
  Flip the switch from afar.
  This is **LUCH**.

