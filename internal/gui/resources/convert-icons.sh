#!/bin/bash

for iconName in trayicon_default trayicon_error trayicon_running trayicon_stopped; do
  sizes=(16 32 48 128 256)
  pngNames=()
  for size in ${sizes[@]}; do
    pngName="$iconName.$size.png"
    inkscape "$iconName.svg" -w $size -h $size -o $pngName >/dev/null 2>/dev/null
    pngNames+=("$pngName")
  done
  convert ${pngNames[@]} "$iconName.ico"
  for pngName in ${pngNames[@]}; do
    if [[ "$pngName" == "$iconName.256.png" ]]; then
      mv $pngName "$iconName.png"
    else
      rm $pngName
    fi
  done
done
