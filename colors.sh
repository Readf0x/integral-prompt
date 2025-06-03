#!/usr/bin/env zsh

# Define color names for 0-15 (standard ANSI colors)
color_names=(
  "Black" "Red" "Green" "Yellow"
  "Blue" "Magenta" "Cyan" "White"
  "Bright Black (Gray)" "Bright Red" "Bright Green" "Bright Yellow"
  "Bright Blue" "Bright Magenta" "Bright Cyan" "Bright White"
)

for i in {0..15}; do
  # Background block, then number and name in foreground color
  printf "\e[48;5;%sm  \e[0m \e[38;5;%sm%2d: %-20s\e[0m\n" "$i" "$i" "$i" "${color_names[$i]}"
done
