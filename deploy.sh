#! /bin/bash

cp -R ./templates ./bin

mkdir ./bin/videos

nohup ./api &
nohup ./scheduler &
nohup ./streamserver &
nohup ./web &

echo "deploy finished"