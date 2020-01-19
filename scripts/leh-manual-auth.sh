#!/bin/sh

BASEDIR=$(dirname "$0")
$BASEDIR/letsencrypt_helper manual-auth -domain=$CERTBOT_DOMAIN -validation=$CERTBOT_VALIDATION