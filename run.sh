#!/bin/bash

pids=()

go build
go build ./cmd/manager
go build ./cmd/provaBluetooth
go build ./cmd/provagui

./miriam &         pids+=($!)

./manager &        pids+=($!)
./provaBluetooth & pids+=($!)

trap 'kill -SIGTERM ${pids[@]}' INT TERM EXIT

sleep infinity