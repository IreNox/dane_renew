#!/bin/sh
letsencrypt_helper dane-update -domain=$RENEWED_DOMAINS -cert=$RENEWED_LINEAGE/fullchain.pem