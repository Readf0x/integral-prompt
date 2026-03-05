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

  PROMPT="%F{$integral_prompt_color}${integral_prompt[4]}%F{15}"
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
      integral helpers cursor-shape 1
      ;;
    viins|main)
      VI_KEYMAP="INSERT"
      integral helpers cursor-shape
      ;;
  esac

  if [[ $VI_KEYMAP != ${previous_vi_keymap} ]]; then
    integral render
    zle reset-prompt
  fi
}

