#!/bin/sh
cd ..
version=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['info']['productVersion'])")
name=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['name'])")
wails build -clean -platform darwin -ldflags "-X 'main.Version=$version'"
wails build -platform windows -ldflags "-X 'main.Version=$version'"
cd ./build/bin
zip -vr ./"$name".app.zip ./"$name".app
rm -rf ./"$name".app
cd ..
if [ -e sign.sh ]
then
    echo "Signing artifacts"
    bash sign.sh "$name"
else
    echo "Signing script not found. Build finished without signing"
fi
