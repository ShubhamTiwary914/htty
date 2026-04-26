#!/bin/bash

USER=ShubhamTiwary914
REPO=htty

CONFIG_DIR=$HOME/.config/htty
COMP_DIR=$HOME/.cache/htty
LOG_DIR=/var/log/htty
LOCAL_BIN=$HOME/.local/bin


BOLD='\033[1m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RESET='\033[0m'


step()  { echo -e "\n${CYAN}${BOLD}==> $1${RESET}"; }
ok()    { echo -e "${GREEN}  ✓ $1${RESET}"; }
warn()  { echo -e "${YELLOW}  ! $1${RESET}"; }

prompt_overwrite() {
  local label=$1
  read -rp "    ${label} already exists. Overwrite? [y/N] " ans
  [[ "$ans" =~ ^[Yy]$ ]]
}

#clean previous fail install (if any)
rm -f /tmp/htty.zip
rm -rf /tmp/htty-master


#pull ------
step "Downloading repository"
curl -fsSL "https://github.com/$USER/$REPO/archive/refs/heads/master.zip" -o /tmp/htty.zip
unzip -q /tmp/htty.zip -d /tmp
ok "Downloaded and extracted"

#config -----------
step "Installing config"
mkdir -p $CONFIG_DIR
if [ -f "$CONFIG_DIR/config.json" ]; then
  if prompt_overwrite "config.json"; then
    cp /tmp/htty-master/config.json "$CONFIG_DIR/"
    ok "Config overwritten"
  else
    warn "Skipped config"
  fi
else
  cp /tmp/htty-master/config.json "$CONFIG_DIR/"
  ok "Config installed -> $CONFIG_DIR/config.json"
fi


#completions -------
step "Installing completions"
mkdir -p $COMP_DIR
existing=false
for f in /tmp/htty-master/setup/default_cmps/*; do
  fname=$(basename "$f")
  if [ -f "$COMP_DIR/$fname" ]; then
    existing=true
    break
  fi
done

if $existing; then
  if prompt_overwrite "completions in $COMP_DIR"; then
    cp /tmp/htty-master/setup/default_cmps/* "$COMP_DIR/"
    ok "Completions overwritten"
  else
    warn "Skipped completions"
  fi
else
  cp /tmp/htty-master/setup/default_cmps/* "$COMP_DIR/"
  ok "Completions installed -> $COMP_DIR/"
fi


# log file ---
step "Setting up log file"
sudo mkdir -p $LOG_DIR
if [ ! -f "$LOG_DIR/htty.log" ]; then
  sudo touch "$LOG_DIR/htty.log"
  ok "Log file created -> $LOG_DIR/htty.log"
else
  ok "Log file already exists, skipping"
fi


# build ---
step "Building htty"
cd /tmp/htty-master/ || exit 1
make build
ok "Build complete"


# install binary -------
step "Installing binary"
if [ -f "$LOCAL_BIN/htty" ]; then
  if prompt_overwrite "htty binary at $LOCAL_BIN/htty"; then
    cp ./htty "$LOCAL_BIN/"
    ok "Binary overwritten -> $LOCAL_BIN/htty"
  else
    warn "Skipped binary install"
  fi
else
  cp ./htty "$LOCAL_BIN/"
  ok "Binary installed -> $LOCAL_BIN/htty"
fi

# cleanup ---
rm -f /tmp/htty.zip
rm -rf /tmp/htty-master


echo -e "\n${GREEN}${BOLD}htty setup complete.${RESET}"
