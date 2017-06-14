#!/usr/bin/env bash

bazel build -c opt //im2txt:run_inference

export DATASET=data.tgz
export MODEL=model
curl -Lo ${DATASET} ${MODEL}
[ -d ${MODEL} ] && rm -rf ${MODEL}
tar zxf ${DATASET} -C ${MODEL}

./engine -ckeckpoint "${MODEL}/train" -wordlist "${MODEL}/data/word_counts.txt"