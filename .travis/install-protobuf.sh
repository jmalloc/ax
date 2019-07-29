#!/bin/sh
set -e
# check to see if protobuf folder is empty
if [ ! -d "$HOME/protobuf/lib" ]; then
  wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
  unzip protoc-3.5.1-linux-x86_64.zip
  chmod +x bin/protoc
  sudo cp bin/protoc /usr/local/bin/
  ls -la /usr/local/bin/protoc
  sudo cp -R include/google/protobuf/* /usr/local/include/
else
  echo "Using cached directory."
fi
