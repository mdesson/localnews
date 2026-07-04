#!/usr/bin/env sh

# build
go build -o builds/localnews

# delete old version
ssh $DEPLOY_USERNAME@$DEPLOY_IP -t rm /bin/localnews

# deploy
scp builds/localnews $DEPLOY_USERNAME@$DEPLOY_IP:/bin/localnews

ssh $DEPLOY_USERNAME@$DEPLOY_IP -t systemctl restart localnews.service

# configure server
scp Caddyfile $DEPLOY_USERNAME@$DEPLOY_IP:/etc/caddy/Caddyfile
ssh $DEPLOY_USERNAME@$DEPLOY_IP systemctl reload caddy
