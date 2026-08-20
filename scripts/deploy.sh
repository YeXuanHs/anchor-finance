#!/bin/bash
# Deploy script - run on server to download latest build and restart
set -e

REPO="YeXuanHs/anchor-finance"
INSTALL_DIR="/opt/anchorfinance"
BINARY_URL="https://github.com/${REPO}/releases/download/latest/anchorfinance"

echo "[1/3] Downloading latest binary..."
curl -sL "$BINARY_URL" -o "${INSTALL_DIR}/anchorfinance.new"
chmod +x "${INSTALL_DIR}/anchorfinance.new"

echo "[2/3] Stopping old process..."
kill -9 $(pgrep -f './anchorfinance') 2>/dev/null || true
sleep 1

echo "[3/3] Starting new binary..."
mv "${INSTALL_DIR}/anchorfinance.new" "${INSTALL_DIR}/anchorfinance"
cd "${INSTALL_DIR}"
nohup ./anchorfinance > /tmp/anchorfinance.log 2>&1 &
echo "PID: $!"
sleep 2

# Verify
if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo "OK - Server is running"
else
    echo "WARNING - Server may not have started correctly"
    tail -5 /tmp/anchorfinance.log
fi
