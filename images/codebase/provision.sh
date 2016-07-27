#!/bin/sh
cp /tmp/provision/common/block-indefinitely /root/block-indefinitely
chmod +x /root/block-indefinitely

mkdir -p /ak/projects/www/codebase

tar -C /ak/projects/www/codebase -xzf /tmp/provision/codebase/package-latest.tar.gz

# Todo: Remove this
# cp /ak/projects/www/codebase/configs/_application.ini \
#  /ak/projects/www/codebase/configs/application.ini

# TodoL Remove this
# sed -i -f /tmp/provision/codebase/application.ini.sed \
#  /ak/projects/www/codebase/configs/application.ini

rm -rf /tmp/*
