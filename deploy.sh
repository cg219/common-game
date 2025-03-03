#!/bin/bash

mkdir -p projects && cd projects

if [ -d "common-game" ]; then
    cd common-game
    git pull
else
    git checkout https://github.com/cg219/common-game.git common-game
    cd common-game
fi

echo $SECRETS_FILE > secrets.yaml
docker stack deploy -c stack.yaml commongame
