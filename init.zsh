autoload -U colors; colors
autoload -Uz add-zsh-hook
autoload -U add-zle-hook-widget

export VI_KEYMAP=${VI_KEYMAP:-"INSERT"}

integral:helpers:cursor-shape() {
  if [[ $1 ]]; then
    echo -ne '\e[1 q'
  else
    echo -ne '\e[5 q'
  fi
}
integral:prompt() {
  case $VI_KEYMAP in
    INSERT)
      visym="%{%F{112}%}○"
      ;;
    VISUAL)
      visym="%{%F{92}%}◒"
      ;;
    V-LINE)
      visym="%{%F{92}%}◐"
      ;;
    NORMAL)
      visym="%{%F{160}%}●"
      ;;
    REPLACE)
      visym="%{%F{32}%}◍"
      ;;
  esac
  local newline=$'\n'
  local prompt_top="$newline%{%F{11}%}⌠$visym "
  local prompt_bot="$newline%{%F{11}%}⌡%{%F{255}%}"

  local dir=${PWD/$HOME/\~}
  if (( ${#dir} >= $COLUMNS )); then
    prompt_top="$prompt_top%{%F{32}%}${dir:0:$(($COLUMNS - 3))}"
    for ((i = 1; i < $((${#dir} / $COLUMNS + 1)); i++)); do
	    prompt_top="$prompt_top$newline%{%F{11}%}⎮%{%F{32}%}${dir:$((($COLUMNS - 1) * $i - 3)):$((($COLUMNS - 1)))}"
    done
  else
    prompt_top="$prompt_top%{%F{32}%}$dir"
  fi

  PROMPT="$prompt_top$prompt_bot"
}

integral:zle-line-init() {
  emulate -L zsh

  [[ $CONTEXT == start ]] || return 0

  while true; do
    zle .recursive-edit
    local -i ret=$?
    [[ $ret == 0 && $KEYS == $'\4' ]] || break
    [[ -o ignore_eof ]] || exit 0
  done

  local saved_prompt=$PROMPT
  local saved_rprompt=$RPROMPT
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
      if [[ $ZLE_KEYSTATE == *overwrite* ]]; then
        VI_KEYMAP="REPLACE"
      else
        VI_KEYMAP="INSERT"
      fi
      integral:helpers:cursor-shape
      ;;
  esac

  if [[ $VI_KEYMAP != ${previous_vi_keymap} ]]; then
    integral:prompt
    zle reset-prompt
  fi
}

add-zsh-hook precmd integral:prompt
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral:prompt
zle -N zle-line-init integral:zle-line-init
