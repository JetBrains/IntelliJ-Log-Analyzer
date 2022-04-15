#!/bin/sh
cd ..
version=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['info']['productVersion'])")
name=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['name'])")
wails build -clean -platform darwin -ldflags "-X 'main.Version=$version'"
wails build -platform windows -ldflags "-X 'main.Version=$version'"
if [ -e ./build/sign.sh ]
then
    echo "Signing artifacts"
    bash ./build/sign.sh "$name" "$version"
else
    echo "Signing script not found. Build finished without signing"
fi
