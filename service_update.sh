# Preemtive superuser access
sudo ls > /dev/null

# Rebuild the server
sudo go build -o /usr/local/bin/fit main.go

# Restart the service
sudo systemctl restart fit.service

# Log
sudo systemctl status fit.service

