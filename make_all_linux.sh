#! /bin/bash

echo "welcome use aoko!"

go fmt ./src/...

cd linux

./make_GameServer.sh

./make_Externalgws.sh

./make_Innergws.sh

./make_loginserver.sh

./make_simulate.sh
