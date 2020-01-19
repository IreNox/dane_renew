#!/bin/sh

BASEDIR=$(dirname "$0")
$BASEDIR/letsencrypt_helper manual-cleanup -domain=$CERTBOT_DOMAIN