#!/bin/sh

BASEDIR=$(dirname "$0")
$BASEDIR/letsencrypt_helper dane-update -domain=$RENEWED_DOMAINS -cert=$RENEWED_LINEAGE/fullchain.pem