#!/usr/bin/env bash

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # this regex checks if a proto file has its go_package set to github.com/alice/checkers/api/...
    # gogo proto files SHOULD ONLY be generated if this is false
    # you don't want gogo proto to run for proto files which are natively built for google.golang.org/protobuf
    if grep -q "option go_package" "$file" && grep -H -o -c 'option go_package.*github.com/verana-labs/verana-blockchain/x/checkers/api' "$file" | grep -q ':0$'; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

echo "Generating pulsar proto code"
buf generate --template buf.gen.pulsar.yaml

cd ..

# cp -r github.com/verana-labs/verana-blockchain/x/checkers* ./
# rm -rf ./api/veranablockchain/checkers && mkdir ./api/veranablockchain/checkers

# cd proto
# mv verana-labs/verana-blockchain/x/checkers/* ./api/veranablockchain/checkers
# rm -rf github.com verana-labs

# cd ..

