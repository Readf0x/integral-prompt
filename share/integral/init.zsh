# shamelessly stolen from p10k
# https://github.com/romkatv/powerlevel10k/issues/888
__integral_zle-line-init() {
  emulate -L zsh

  [[ $CONTEXT == start ]] || return 0

  while true; do
    zle .recursive-edit
    local -i ret=$?
    [[ $ret == 0 && $KEYS == $'\4' ]] || break
    [[ -o ignore_eof ]] || exit 0
  done

  PROMPT=$(integral transient)
  RPROMPT=''
  zle .reset-prompt

  if (( ret )); then
    zle .send-break
  else
    zle .accept-line
  fi
  return ret
}

__integral_cursor-shape() {
  if [[ $1 ]]; then
    echo -ne '\e[1 q'
  else
    echo -ne '\e[5 q'
  fi
}

export sig=0
__integral_error_hook() {
  export sig=$?
}

__integral_render() {
  eval "$(integral render zsh $sig $(jobs | wc -l))"
}

__integral_line-pre-redraw() {
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
      __integral_cursor-shape 1
      ;;
    viins|main)
      VI_KEYMAP="INSERT"
      __integral_cursor-shape
      ;;
  esac

  if [[ $VI_KEYMAP != ${previous_vi_keymap} ]]; then
    __integral_render
    zle reset-prompt
  fi
}

autoload -Uz add-zsh-hook
autoload -Uz add-zle-hook-widget

export VI_KEYMAP=${VI_KEYMAP:-"INSERT"}
export HOSTNAME=${HOSTNAME:-"$(hostname)"}
export INTEGRAL_INSTALL_PATH=$int_path

# === INIT ===
TRAPWINCH() {
  __integral_render
  zle && zle reset-prompt
}
add-zsh-hook precmd __integral_error_hook
add-zsh-hook precmd __integral_render
add-zsh-hook precmd __integral_cursor-shape
zle -N __integral_line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw __integral_line-pre-redraw
__integral_render
zle -N zle-line-init __integral_zle-line-init
