#!/bin/bash

for D in $(ls -1d sd*); do
  go test "github.com/gaorx/stardust3/$D"
done