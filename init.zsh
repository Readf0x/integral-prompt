local version='0.0.5'

source './module_loader.zsh'
autoload -U colors; colors
autoload -Uz add-zsh-hook
autoload -U add-zle-hook-widget

export VI_KEYMAP=${VI_KEYMAP:-"INSERT"}

# === OPTIONS ===
export integral_vim_color="true"
export integral_vim_indicators=(
  "○" # insert
  "◒" # visual
  "◐" # v-line
  "●" # normal
)
export integral_top="⌠"
export integral_mid="⎮"
export integral_bot="⌡"

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
for f in $rc_locations; do
  if [[ -f $f ]]; then
    source $f
  fi
done

# === HELPERS ===
integral:helpers:cursor-shape() {
  if [[ $1 ]]; then
    echo -ne '\e[1 q'
  else
    echo -ne '\e[5 q'
  fi
}
integral:helpers:version() {
  print "v$version"
}
# TODO: Create "class" functions

# === MAIN LOGIC ===
integral:prompt() {
  integral:loop_modules
}

# === ZLE ===
# TODO: rerender on sigwinch
#   This can supposedly be done with zsh-async
# shamelessly stolen from p10k
# https://github.com/romkatv/powerlevel10k/issues/888
integral:zle-line-init() {
  emulate -L zsh

  [[ $CONTEXT == start ]] || return 0

  while true; do
    zle .recursive-edit
    local -i ret=$?
    [[ $ret == 0 && $KEYS == $'\4' ]] || break
    [[ -o ignore_eof ]] || exit 0
  done

  PROMPT="%{%F{11}%}∫%{$reset_color%}"
  RPROMPT=''
  zle .reset-prompt

  if (( ret )); then
    zle .send-break
  else
    zle .accept-line
  fi
  return ret
}

integral:line-pre-redraw() {
  local previous_vi_keymap="$VI_KEYMAP"

  case $KEYMAP in
    vicmd)
      case $REGION_ACTIVE in
        1)
          VI_KEYMAP="VISUAL"
          ;;
        2)
          VI_KEYMAP="V-LINE"
          ;;
        *)
          VI_KEYMAP="NORMAL"
          ;;
      esac
      integral:helpers:cursor-shape 1
      ;;
    viins|main)
      VI_KEYMAP="INSERT"
      integral:helpers:cursor-shape
      ;;
  esac

  if [[ $VI_KEYMAP != ${previous_vi_keymap} ]]; then
    integral:prompt
    zle reset-prompt
  fi
}

# === INIT ===
TRAPWINCH() {
  integral:prompt
  zle reset-prompt
}
add-zsh-hook preexec error_hook
add-zsh-hook precmd integral:prompt
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral:prompt
zle -N zle-line-init integral:zle-line-init

