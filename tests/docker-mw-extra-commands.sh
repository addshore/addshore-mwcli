#!/usr/bin/env bash

set -e # Fail on errors
set -x # Output commands

# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "\"${last_command}\" command filed with exit code $?."' EXIT

# Create
./bin/mw docker mediawiki create

# Validate the basic stuff
./bin/mw docker docker-compose ps
./bin/mw docker env list
# TODO enable logo check again once the page no longer shown "Unable to connect to PostgreSQL server"
#CURL=$(curl -s -L -N http://default.mediawiki.mwdd.localhost:8080) && echo $CURL && echo $CURL | grep -q "The MediaWiki logo"
CURL=$(curl -s -L -N http://default.mediawiki.mwdd.localhost:8080) && echo $CURL && echo $CURL | grep -q "Unable to connect to PostgreSQL server"

# Install sqlite & check
./bin/mw docker mediawiki install
CURL=$(curl -s -L -N http://default.mediawiki.mwdd.localhost:8080) && echo $CURL && echo $CURL | grep -q "MediaWiki has been installed"

# docker-compose: Make sure it appears to work
./bin/mw docker docker-compose ps -- --services | grep -q "mediawiki"

# cd to mediawiki
cd mediawiki

# composer: Make sure a command works in root of the repo
./../bin/mw docker mediawiki composer home | grep -q "https://www.mediawiki.org/"
# phpunit: Make sure a command works in the root of the repo
./../bin/mw docker mediawiki phpunit ./tests/phpunit/unit/includes/PingbackTest.php
./../bin/mw docker mediawiki phpunit ./tests/phpunit/unit/includes/PingbackTest.php | grep -q "OK "
# exec: Make sure a command works in the root of the repo
./../bin/mw docker mediawiki exec ls | grep -q "api.php"

# cd to Vector
cd skins/Vector

# composer: Make sure a command works from the Vector directory
./../../../bin/mw docker mediawiki composer home | grep -q "http://gerrit.wikimedia.org/g/mediawiki/skins/Vector"
# phpunit: Make sure a command works from the Vector directory
./../../../bin/mw docker mediawiki phpunit ./../../tests/phpunit/unit/includes/PingbackTest.php | grep -q "OK "
# exec: Make sure a command works from the Vector directory
# Right now this just executes in the MediaWiki directory
./../../../bin/mw docker mediawiki exec ls | grep -q "api.php"

# cd back again
cd ./../../../

# Destroy it all
./bin/mw docker destroy
# And make sure only 1 exists after
docker ps
docker ps | wc -l | grep -q "1"