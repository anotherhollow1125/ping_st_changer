#!/usr/bin/bash

cp ping_st_changer /usr/local/bin/ping_st_changer
chmod +x /usr/local/bin/ping_st_changer
cp ping_st_changer.service /etc/systemd/system/ping_st_changer.service
chmod -x ping_st_changer.service /etc/systemd/system/ping_st_changer.service
cp ping_st_changer.timer /etc/systemd/system/ping_st_changer.timer
chmod -x ping_st_changer.timer /etc/systemd/system/ping_st_changer.timer

echo "Please Overwrite /etc/systemd/system/ping_st_changer.(service | timer) with your own settings"
echo "after setting, please run the under command as root:"
echo "  systemctl enable --now ping_st_changer.timer"
