#!/bin/sh
/usr/bin/certbot renew --manual-auth-hook /root/letsencrypt_helper/leh-manual-auth.sh --manual-cleanup-hook /root/letsencrypt_helper/leh-manual-cleanup.sh --deploy-hook /root/letsencrypt_helper/leh-dane-update.sh