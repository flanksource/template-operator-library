#!/bin/bash

kubectl -n elastic-system delete elasticsearch estest
kubectl -n elastic-system delete elasticsearchdb estest

kubectl -n redis-operator delete redisfailover redisdb-e2e
kubectl -n redis-operator delete redisdb redisdb-e2e