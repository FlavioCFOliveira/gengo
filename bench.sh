#!/bin/bash

go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s .