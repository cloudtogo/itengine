#!/usr/bin/env bash

bazel build -c opt //im2txt:run_inference

export DATASET=data.tgz
export MODEL_POS=model
curl -Lo ${DATASET} ${MODEL}
[ -d ${MODEL_POS} ] && rm -rf ${MODEL_POS}
tar zxf ${DATASET} -C ${MODEL_POS}

./engine -ckeckpoint "${MODEL_POS}/train" -wordlist "${MODEL_POS}/data/word_counts.txt"