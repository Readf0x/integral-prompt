#!/usr/bin/env zsh

color_names=(
  "Black" "Red" "Green" "Yellow"
  "Blue" "Magenta" "Cyan" "White"
  "Bright Black (Gray)" "Bright Red" "Bright Green" "Bright Yellow"
  "Bright Blue" "Bright Magenta" "Bright Cyan" "Bright White"
)

for i in {0..15}; do
  if [[ $i = (0|8) ]]; then
    print -P "%K{$i}  %k%K{7} %F{$i}$i: ${color_names[$i + 1]}%f%k"
  else
    print -P "%K{$i}  %k %F{$i}$i: ${color_names[$i + 1]}%f"
  fi
done
