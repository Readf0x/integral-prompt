# TODO: Make modules configurable
# === OPTIONS ===
export int_modules=(
  "direnv"
  "nix"
  "visym"
  "error"
  "dir"
  "git"
  "jobs"
)
export int_right_modules=(
  "time"
)
export int_kitty_integration="false"
export int_prompt_color="11"
export int_prompt=(
  "⌠"
  "⎮"
  "⌡"
  "∫"
)
export int_nix_icons=(
  "❄"
  ""
)
export int_nix_color=(
  "14"
  "13"
)
export int_vim_indicators=(
  "○" # insert
  "◒" # visual
  "◐" # v-line
  "●" # normal
)
export int_vim_colors=(
  "10" # insert
  "13" # visual
  "13" # v-line
  "9"  # normal
)
export int_error_format() {
  case $1 in
    1) print "%F{9}✘" ;;
    2|127) print "%F{11}?" ;;
    126) print "%F{9}⚠" ;;
    130) print "%F{15}☠" ;;
    148) print "%F{13}✱" ;;
    *) print "%F{9}✘" ;;
  esac
}
export int_dir_format() { print ${PWD/$HOME/\~} }
export int_dir_color="12"
export int_git_icons=(
  "⎇"
  "✘"
  "+"
  "↑"
  "↓"
)
export int_git_colors=(
  "11"
  "9"
  "11"
  "14"
  "14"
)
export int_jobs_icon="⚙"
export int_jobs_color="13"
export int_time_format="%T"
export int_uptime_icon="⏲"
export int_uptime_color="12"
export int_battery_icons=(
  "🗲"
  "󰁹"
)
export int_battery_color="10"
export int_ssh_format="%F{12}${USER}%F{13}@%F{14}${HOSTNAME}"

export int_direnv_format() {
  case $1 in
    *) print "%F{11}⌁" ;;
  esac
}

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
if $int_kitty_integration && [[ $KITTY_PID ]] && [[ $(kitty +kitten query_terminal | grep font) =~ "NF|Nerd ?Font" ]]; then
  export int_nerd_fonts="true"
else
  export int_nerd_fonts="false"
fi