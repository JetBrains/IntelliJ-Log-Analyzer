#!/bin/bash
function RenameAndCompress() {
  local name=$1
  local version=$2
  local platformname=$3
  local pwd=$(pwd)
  cd "./build/bin/"
  for file in "$name".app; do
      local extension="${file##*.}"
      local standartName="$name"-"$version"-"$platformname"
      zip -vr "$standartName".zip "$file"
      rm -r "$file"
  done
  for file in "$name".exe; do
          local extension="${file##*.}"
          local standartName="$name"-"$version"-"$platformname"."$extension"
          mv "$file" "$standartName"
  done
  cd $pwd
}

cd ..
version=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['info']['productVersion'])")
name=$(cat ./wails.json | python3 -c "import sys, json; print(json.load(sys.stdin)['name'])")
rm -rf ./build/bin/*
ls -l ./build/bin/*

for platform in  "darwin/arm64" "windows/amd64" "darwin/amd64"
do
  platformname=$(echo $platform | sed 's/\//-/g')
  wails build -platform $platform -ldflags "-X 'main.Version=$version'"
  RenameAndCompress "$name" "$version" "$platformname"
done
if [ -e ./build/sign.sh ]
then
    echo "Signing artifact $filename"
    bash ./build/sign.sh "$name"
else
    echo "Signing script not found. Build finished without signing"
fi