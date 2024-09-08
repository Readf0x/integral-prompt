# TODO: Make modules configurable
# === OPTIONS ===
export integral_vim_color="true"
export integral_vim_indicators=(
  "○" # insert
  "◒" # visual
  "◐" # v-line
  "●" # normal
)
export integral_top="%{%F{11}%}⌠"
export integral_mid="%{%F{11}%}⎮"
export integral_bot="%{%F{11}%}⌡"

# === CONFIG LOADING ===
local rc_locations=(
  ~/.integralrc
  $XDG_CONFIG_HOME/integralrc
  $XDG_CONFIG_HOME/integral/rc
  $XDG_CONFIG_HOME/integral/rc.zsh
  ~/.config/integralrc
  ~/.config/integral/rc
  ~/.config/integral/rc.zsh
)
export integral_modules=(
  "visym"
  "error"
  "dir"
  "git"
)
for f in $rc_locations; do
  if [[ -f $f ]]; then
    source $f
  fi
done

