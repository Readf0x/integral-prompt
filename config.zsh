# TODO: Make modules configurable
# === OPTIONS ===
export integral_modules=(
  "nix"
  "visym"
  "error"
  "dir"
  "git"
  "jobs"
)
export integral_right_modules=(
  "time"
)
export integral_kitty_integration="false"
export integral_prompt_color="11"
export integral_prompt=(
  "⌠"
  "⎮"
  "⌡"
  "∫"
)
export integral_nix_icons=(
  "❄"
  ""
)
export integral_nix_color=(
  "14"
  "13"
)
export integral_vim_indicators=(
  "○" # insert
  "◒" # visual
  "◐" # v-line
  "●" # normal
)
export integral_vim_colors=(
  "10" # insert
  "13" # visual
  "13" # v-line
  "9"  # normal
)
export integral_error_format() {
  case $1 in
    1) print "%{%F{9}%}✘" ;;
    2|127) print "%{%F{11}%}?" ;;
    126) print "%{%F{9}%}⚠" ;;
    130) print "%{%F{15}%}☠" ;;
    *) print "%{%F{9}%}✘" ;;
  esac
}
export integral_dir_format() { print ${PWD/$HOME/\~} }
export integral_dir_color="12"
export integral_git_icons=(
  "⎇"
  "✘"
  "+"
  "↑"
  "↓"
)
export integral_git_colors=(
  "11"
  "9"
  "11"
  "14"
  "14"
)
export integral_jobs_icon="⚙"
export integral_jobs_color="13"
export integral_time_format="%T"
export integral_uptime_icon="⏲"
export integral_uptime_color="12"

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
for f in $rc_locations; do
  if [[ -f $f ]]; then
    source $f
  fi
done
if $integral_kitty_integration && [[ $KITTY_PID ]] && [[ $(kitty +kitten query_terminal | grep font) =~ "NF|Nerd ?Font" ]]; then
  export integral_nerd_fonts="true"
else
  export integral_nerd_fonts="false"
fi
