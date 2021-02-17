
#!/bin/bash

service_name=${PWD##*/}

credentials=claudius@85.186.121.41
target_dir=/home/claudius/services/$service_name/

# stop the service
#ssh -t $credentials sudo systemctl stop $service_name

time rsync -P -rsh=ssh $service_name $credentials:$target_dir

# start the service
#ssh -t $credentials sudo systemctl start $service_name
