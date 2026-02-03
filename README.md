# FCKGOBACK

![FCKGOBACK](assets/fckgoback.png)

Utility for rolling back Arch Linux updates using the Arch Linux Archive.

## What is this?

Interactive utility for temporarily switching to an archived date and rolling back problematic updates. Automatically manages mirrorlist and restores it after use.

## Quick Start

```bash
# Build
git clone https://github.com/Yoxi7/fckgoback.git
cd fckgoback
go build -o fckgoback ./cmd/fckgoback

# Create alias (recommended)
echo "alias fckgoback='sudo $(pwd)/fckgoback'" >> ~/.bashrc  # or ~/.zshrc
source ~/.bashrc

# Usage
fckgoback
```

## How it works

1. **Date selection** — interactive menu to choose year/month/day from archive
2. **Availability check** — automatic verification that archive exists
3. **Mirrorlist backup** — save current configuration
4. **Temporary switch** — write archive URL to `/etc/pacman.d/mirrorlist`
5. **Update** — execute `pacman -Syyuu` to rollback packages
6. **Restore** — automatically restore original mirrorlist


## Features

- Automatic archive availability check before use
- Mirrorlist backup and restore on any outcome
- Russian and English language support (auto-detection)
- Graceful shutdown — restore on Ctrl+C
- All repositories check (core, extra, multilib)

## Security

- Automatic backup before changes
- Restore on errors and interruption
- Requires root privileges (system file operations)

## License

MIT
