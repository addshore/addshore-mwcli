#!/usr/bin/env bash

# Fail on errors
set -e
# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "\"${last_command}\" command filed with exit code $?."' EXIT
# Output commands
set -x

# Output version
./mw version

# Setup & Create
./mw docker env set PORT 8080
./mw docker env set MEDIAWIKI_VOLUMES_CODE $(pwd)/mediawiki
./mw docker mediawiki create

# Validate the basic stuff
./mw docker docker-compose ps
./mw docker env list
cat ~/.mwcli/mwdd/default/.env
CURL=$(curl -s -L -N http://default.mediawiki.mwdd.localhost:8080) && echo $CURL && echo $CURL | grep -q "The MediaWiki logo"

# Add the needed LocalSettings
echo "<?php" >> mediawiki/LocalSettings.php
echo "//require_once "$IP/includes/PlatformSettings.php";" >> mediawiki/LocalSettings.php
echo "require_once '/mwdd/MwddSettings.php';" >> mediawiki/LocalSettings.php

# Install sqlite & check
./mw docker mediawiki install
CURL=$(curl -s -L -N http://default.mediawiki.mwdd.localhost:8080) && echo $CURL && echo $CURL | grep -q "MediaWiki has been installed"

# docker-compose: Make sure it appears to work
./mw docker docker-compose ps -- --services | grep -q "mediawiki"

# cd to mediawiki
cd mediawiki

# composer: Make sure a command works in root of the repo
./../mw docker mediawiki composer home | grep -q "https://www.mediawiki.org/"
# phpunit: Make sure a command works in the root of the repo
./../mw docker mediawiki phpunit ./tests/phpunit/unit/includes/PingbackTest.php | grep -q "OK "
# exec: Make sure a command works in the root of the repo
./../mw docker mediawiki exec ls | grep -q "LocalSettings.php"

# cd to Vector
cd skins/Vector

# composer: Make sure a command works from the Vector directory
./../../../mw docker mediawiki composer home | grep -q "http://gerrit.wikimedia.org/g/mediawiki/skins/Vector"
# phpunit: Make sure a command works from the Vector directory
./../../../mw docker mediawiki phpunit ./../../tests/phpunit/unit/includes/PingbackTest.php | grep -q "OK "
# exec: Make sure a command works from the Vector directory
# Right now this just executes in the MediaWiki directory
./../../../mw docker mediawiki exec ls | grep -q "LocalSettings.php"

# cd back again
cd ./../../../

# Destroy it all
./mw docker destroy
# And make sure only 1 exists after
docker ps
docker ps | wc -l | grep -q "1"
