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
integral:helpers:real-length() {
  local x=$(print $1 | sed 's/%{\(%F\|%B\)\{0,2\}{[0-9]*}%}//g')
  export debug_len=$x
  print ${#x}
}

# === MAIN LOGIC ===
integral:prompt() {
  # Variables
  case $VI_KEYMAP in
    INSERT)
      visym="%{%F{10}%}${integral_vim_indicators[1]}"
      ;;
    VISUAL)
      visym="%{%F{13}%}${integral_vim_indicators[2]}"
      ;;
    V-LINE)
      visym="%{%F{13}%}${integral_vim_indicators[3]}"
      ;;
    NORMAL)
      visym="%{%F{9}%}${integral_vim_indicators[4]}"
      ;;
  esac
  local dir=${PWD/$HOME/\~}
  local git

  # Constants
  local newline=$'\n'
  local prompt_top="$newline%{%F{11}%}⌠$visym%(?.. %{%F{9}%}⚠)"
  local prompt_bot="$newline%{%F{11}%}⌡%{%F{255}%}"

  if [ -d .git ] || git rev-parse --git-dir >/dev/null 2>&1; then
    # TODO: add support for unstaged changes, untracked files, etc
    git="%{%F{11}%}$(git rev-parse --abbrev-ref HEAD)⎇"
    if ! git diff --quiet --ignore-submodules 2>/dev/null; then
      git="$git %{%F{9}%}$(git diff --no-ext-diff --ignore-submodules --stat | tail -n1 | cut -d' ' -f2)✘"
    fi
    if ! git diff --quiet --ignore-submodules --cached 2>/dev/null; then
      git="$git %{%F{11}%}$(git diff --no-ext-diff --ignore-submodules --stat --cached | tail -n1 | cut -d' ' -f2)+"
    fi
  fi
  # BUG: this will not work if other modules push the prompt past the terminal width
  # TODO: Change the loop method to be more efficient
  #   Instead of whatever the hell is happening here, we can create static module
  #   "classes" (just a function that follows a certain schema in this case) that output
  #   a preformatted string as well as their raw length. Then we can loop through all
  #   the modules and add their length to the current length of the line. If it exceeds
  #   the terminal width, we just insert our newline and set the length back to 0.
  if (( ${#dir} >= $COLUMNS )); then
    # This is where stuff gets messy, shell string concatenation is less than readable.
    # Set the first line, must be done seperately because of the indicator
    prompt_top="$prompt_top %{%F{12}%}${dir:0:$(($COLUMNS - 3))}"
    # Determine the number of lines the rest of the directory path will take up
    local x=$((${#dir} / $COLUMNS + 1))
    for ((i = 1; i < x; i++)); do
      # Append the correct section of the path by offsetting the substring based on $i
      prompt_top="$prompt_top$newline%{%F{11}%}⎮%{%F{12}%}${dir:$((($COLUMNS - 1) * $i - 3)):$((($COLUMNS - 1)))}"
    done
    # TODO: detect if the last line is short enough to not need a new line and if the
    # rest of the "modules" need to go through a similar process to the directory path.
    if [[ $git ]]; then
      prompt_top="$prompt_top$newline%{%F{11}%}⎮%{%F{12}%}$git"
    fi
  else
    # real simple if the directory path is short enough
    prompt_top="$prompt_top %{%F{12}%}$dir $git"
  fi

  PROMPT="$prompt_top$prompt_bot"
}

# === ZLE ===
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
add-zsh-hook precmd integral:prompt
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral:prompt
zle -N zle-line-init integral:zle-line-init

