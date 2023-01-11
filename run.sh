#!/bin/bash
cd /home/ec2-user/go-aws-ec2
docker-compose build --no-cache
docker-compose up -d
