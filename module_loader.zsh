# === MODULES ===
integral:module:git() {
  # TODO: improve efficiency by storing repeated calls in variables
  local format_str length
  if ! $(git rev-parse --is-bare-repository 2>/dev/null); then
    if [ -d .git ] || git rev-parse --git-dir >/dev/null 2>&1; then
      local branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
      length=$((${#branch} + 1))
      format_str="%{%F{11}%}$branch⎇"
      if ! git diff --quiet --ignore-submodules 1>/dev/null 2>&1 || [[ $(git ls-files -o --exclude-standard) ]]; then
        local num=$(($(git ls-files -o --exclude-standard | wc -l) + $(git diff --name-only | wc -l)))
        length=$(($length + $num + 2))
        format_str="$format_str %{%F{9}%}$num✘"
      fi
      if ! git diff --quiet --ignore-submodules --cached 1>/dev/null 2>&1; then
        local num=$(git diff --ignore-submodules --cached --name-only | wc -l)
        length=$(($length + $num + 2))
        format_str="$format_str %{%F{11}%}$num+"
      fi
      if [[ $(git remote) ]] && git cherry >/dev/null 2>&1 && [[ $(git cherry | wc -l) -gt 0 ]]; then
        local num=$(git cherry | wc -l)
        length=$(($length + $num + 2))
        format_str="$format_str %{%F{14}%}$num↑"
      fi
      if [[ $(git rev-list origin/${branch} --not HEAD 2>/dev/null) ]]; then
        local num=$(git rev-list origin/${branch} --not HEAD --count)
        length=$(($length + $num + 2))
        format_str="$format_str %{%F{14}%}$num↓"
      fi
    fi
  else
    format_str="%{%F{11}%}bare⎇"
    length="5"
  fi

  if [[ $1 ]]; then
    print "$length"
  else
    print "$format_str"
  fi
}

integral:module:visym() {
  local format_str
  case $VI_KEYMAP in
    INSERT)
      format_str="%{%F{10}%}${integral_vim_indicators[1]}"
      ;;
    VISUAL)
      format_str="%{%F{13}%}${integral_vim_indicators[2]}"
      ;;
    V-LINE)
      format_str="%{%F{13}%}${integral_vim_indicators[3]}"
      ;;
    NORMAL)
      format_str="%{%F{9}%}${integral_vim_indicators[4]}"
      ;;
  esac

  if [[ $1 ]]; then
    # We can hardcode the length here kek
    print '1'
  else
    print "$format_str"
  fi
}

integral:module:dir() {
  # this one is gonna be pretty difficult, idk how I'm gonna handle it tbh
  #   potential solutions:
  #     - append 'w' to length to indicate wrapping
  #     - add settings option (e.g. "$format_str:$length:w")
  #     - hardcode wrapping into module loader (bad idea)

  local dir=${PWD/$HOME/\~}
  if [[ $1 ]]; then
    print "${#dir}"
  else
    print "%{%F{12}%}$dir"
  fi
}

error_hook() {
  export sig=$?
}
integral:module:error() {
  if [[ $1 ]]; then
    print "$sig"
  else
    if [[ $sig -eq 1 ]]; then
      print "%{%F{11}%}⚠"
    else
      print "%{%F{9}%}✘"
    fi
  fi
}

export INTEGRAL_MODULES=(
  "visym"
  "error"
  "dir"
  "git"
)

integral:loop_modules() {
  newline=$'\n'
  PROMPT="$newline%{%F{11}%}$integral_top"
  for module in $INTEGRAL_MODULES; do
    local length=$(integral:module:$module 1)
    local format_str=$(integral:module:$module)
    if [[ $length -gt 0 ]]; then
      PROMPT+="$format_str "
    fi
  done
  PROMPT+="$newline%{%F{11}%}$integral_bot%{${reset_color}%}"
}

