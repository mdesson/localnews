#!/usr/bin/env sh

# build
go build -o builds/localnews

# deploy
scp builds/localnews $DEPLOY_USERNAME@$DEPLOY_IP:/bin/localnews

ssh $DEPLOY_USERNAME@$DEPLOY_IP -t systemctl restart localnews.service