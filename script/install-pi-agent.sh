#!/usr/bin/env bash
set -euo pipefail

SERVICE_FILE="./installable/raspivid-stream.service"
DEST_DIR="/etc/systemd/system"
SERVICE_NAME="raspivid-stream.service"

# 1. Copy the service file, replacing it if it already exists
if [ -f "${DEST_DIR}/${SERVICE_NAME}" ]; then
  echo "${SERVICE_NAME} already exists in ${DEST_DIR}, replacing it."
  sudo cp "${SERVICE_FILE}" "${DEST_DIR}/"
  sudo chmod 644 "${DEST_DIR}/${SERVICE_NAME}"
else
  echo "Copying ${SERVICE_FILE} to ${DEST_DIR}/"
  sudo cp "${SERVICE_FILE}" "${DEST_DIR}/"
  sudo chmod 644 "${DEST_DIR}/${SERVICE_NAME}"
fi

# 2. Reload systemd to pick up any new or updated unit
echo "Reloading systemd manager configuration..."
sudo systemctl daemon-reload

# 3. Enable the service at boot if not already enabled
if systemctl is-enabled --quiet "${SERVICE_NAME}"; then
  echo "${SERVICE_NAME} is already enabled."
else
  echo "Enabling ${SERVICE_NAME} to start on boot..."
  sudo systemctl enable "${SERVICE_NAME}"
fi

echo "${SERVICE_NAME} is now installed, updated, and enabled but not started."