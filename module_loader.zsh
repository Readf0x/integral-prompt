# === MODULES ===
integral:module:git() {
  # TODO: improve efficiency by storing repeated calls in variables
  # Should each icon be it's own module? Would need to export branch name...
  local format_str length
  if ! $(git rev-parse --is-bare-repository 2>/dev/null); then
    if [ -d .git ] || git rev-parse --git-dir >/dev/null 2>&1; then
      local branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
      length=$((${#branch} + 1))
      format_str="%{%F{11}%}$branch⎇"
      if ! git diff --quiet --ignore-submodules 1>/dev/null 2>&1 || [[ $(git ls-files -o --exclude-standard) ]]; then
        local num=$(($(git ls-files -o --exclude-standard | wc -l) + $(git diff --name-only | wc -l)))
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{9}%}$num✘"
      fi
      if ! git diff --quiet --ignore-submodules --cached 1>/dev/null 2>&1; then
        local num=$(git diff --ignore-submodules --cached --name-only | wc -l)
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{11}%}$num+"
      fi
      if [[ $(git remote) ]] && git cherry >/dev/null 2>&1 && [[ $(git cherry | wc -l) -gt 0 ]]; then
        local num=$(git cherry | wc -l)
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{14}%}$num↑"
      fi
      if [[ $(git rev-list origin/${branch} --not HEAD 2>/dev/null) ]]; then
        local num=$(git rev-list origin/${branch} --not HEAD --count)
        length=$(($length + ${#num} + 2))
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
  if [[ $1 == "w" ]]; then
    return 1
  elif [[ $1 == "r" ]]; then
    print "$dir"
  elif [[ $1 == "c" ]]; then
    print "%{%F{12}%}"
  elif [[ $1 ]]; then
    print "${#dir}"
  else
    print "%{%F{12}%}$dir"
  fi
}

export sig=0
error_hook() {
  export sig=$?
}
integral:module:error() {
  local format_str
  if [[ $1 ]]; then
    if [[ $sig == 0 ]]; then
      print "0"
    else
      print "1"
    fi
  else
    case $sig in
      1) format_str="%{%F{9}%}✘" ;;
      2|127) format_str="%{%F{11}%}?" ;;
      126) format_str="%{%F{9}%}⚠" ;;
      130) format_str="%{%F{15}%}☠" ;;
      *) format_str="%{%F{9}%}✘" ;;
    esac

    print "$format_str"
  fi
}

export integral_modules=(
  "visym"
  "error"
  "dir"
  "git"
)

# BUG: leaves <space> at end of prompt
integral:loop_modules() {
  local -i position=0
  if [[ $2 ]]; then
    local -i max_len=$2
  else
    local -i max_len=$(($COLUMNS - 1))
  fi

  [[ $1 ]] && print -P "%{%F{14}%}\$max_len: %{%F{13}%}$max_len"

  local newline=$'\n'
  PROMPT="$newline$integral_top"

  for module in $integral_modules; do
    [[ $1 ]] && print -P "%{%F{12}%}current module: %{%F{11}%}$module%{%F{15}%}"
    local -i length=$(integral:module:$module 1)
    local format_str=$(integral:module:$module)

    [[ $1 == "1" ]] && print "$length:$format_str"
    [[ $1 == "2" ]] && print -P "$length:$format_str%{%F{15}%}%}"

    if [[ $length -gt 0 ]]; then
      local new_pos=$(($position + $length + 1))
      [[ $1 ]] && print -P "%{%F{14}%}\$new_pos: %{%F{13}%}$new_pos"
      if [[ $length -gt $max_len ]] && ! integral:module:$module w; then
        local raw_str=$(integral:module:$module r)
        local color=$(integral:module:$module c)
        [[ $1 ]] && print -P "%{%F{14}%}length > max_len: %{%F{13}%}$length > $max_len"
        local -i i=0
        while [[ $i -le $((${#raw_str} / $max_len)) ]]; do
          PROMPT+="$newline$integral_mid$color${raw_str:$(($i * $max_len)):$max_len}"
          [[ $1 == "1" ]] && print "${(%):-%F{15}PROMPT: $PROMPT" #}" #idk why but treesitter is freaking out over this
          [[ $1 == "2" ]] && print -P $PROMPT
          i+=1
        done
        PROMPT+=" "
        position=$((${#raw_str} % $max_len))
        [[ $1 ]] && print -P "%{%F{14}%}position changed to: %{%F{13}%}$position"
      elif [[ $new_pos -gt $max_len ]]; then
        [[ $1 ]] && print -P "%{%F{14}%}new_pos > max_len: %{%F{13}%}$new_pos > $max_len"
        PROMPT+="$newline$integral_mid$format_str "
        [[ $1 == "1" ]] && print "${(%):-%F{15}PROMPT: $PROMPT" #}" #idk why but treesitter is freaking out over this
        [[ $1 == "2" ]] && print -P $PROMPT
        position=$length
        [[ $1 ]] && print -P "%{%F{14}%}position changed to: %{%F{13}%}$position"
      else
        PROMPT+="$format_str "
        [[ $1 == "1" ]] && print "${(%):-%F{15}PROMPT: $PROMPT" #}" #idk why but treesitter is freaking out over this
        [[ $1 == "2" ]] && print -P $PROMPT
        position=$new_pos
      fi
    elif [[ $1 ]]; then
      print -P "%{%F{14}%}length is 0, skipping...$newline"
    fi
  done
  PROMPT+="$newline%{%F{11}%}$integral_bot%{%F{15}%}"
  [[ $1 ]] && print -P "%{%F{15}%}%====="
}

