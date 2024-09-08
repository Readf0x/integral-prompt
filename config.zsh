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
  $ZDOTDIR/.integralrc
  $XDG_CONFIG_HOME/integralrc
  $XDG_CONFIG_HOME/integral/rc
  $XDG_CONFIG_HOME/integral/rc.zsh
  ~/.config/integralrc
  ~/.config/integral/rc
  ~/.config/integral/rc.zsh
)
export integral_modules=(
  "nix"
  "visym"
  "error"
  "dir"
  "git"
  "jobs"
)
export integral_kitty_integration="false"
if $integral_kitty_integration && [[ $KITTY_PID ]] && [[ $(kitty +kitten query_terminal | grep font) =~ "NF|Nerd ?Font" ]]; then
  export integral_nerd_fonts="true"
else
  export integral_nerd_fonts="false"
fi
for f in $rc_locations; do
  if [[ -f $f ]]; then
    source $f
  fi
done

