#!/bin/sh
mkdir -p ~/bin

chmod +x *.sh
chmod +x letsencrypt_helper

cp letsencrypt_helper ~/bin
cp leh-dane-update.sh /etc/letsencrypt/renewal-hooks/deploy