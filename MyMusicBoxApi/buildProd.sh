sudo systemctl stop mymusic
sudo go build -buildvcs=false -o /home/admin/mymusicbox_production
#sudo ls -l /home/admin/mymusicbox_production
sudo systemctl start mymusic
echo "Press CTRL + C to escape"
sudo journalctl -u mymusic.service -f