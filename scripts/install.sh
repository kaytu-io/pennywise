#!/bin/bash

sudo -v
curl -s -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kaytu-io/pennywise/releases/latest | jq -r .assets[].browser_download_url | grep linux | grep amd | xargs wget -O ./pennywise -qi -
chmod +x ./pennywise
sudo mv ./pennywise /usr/bin/