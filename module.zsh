# === MODULES ===
# BUG: length gets printed to console in some repos
integral:module:git() {
  # TODO: improve efficiency by storing repeated calls in variables
  local format_str length
  if ! $(git rev-parse --is-bare-repository 2>/dev/null); then
    if [ -d .git ] || git rev-parse --git-dir >/dev/null 2>&1; then
      local branch=$(git branch --show-current 2>/dev/null)
      if [[ $branch == "" ]]; then
        x=$(git rev-parse HEAD 2>/dev/null)
        branch=${x:0:7}
      fi
      if [[ $branch == "" ]]; then
        branch="detached"
      fi
      length=$((${#branch} + 1))
      format_str="%{%F{${integral_git_colors[1]}}%}$branch${integral_git_icons[1]}"
      if ! git diff --quiet --ignore-submodules 1>/dev/null 2>&1 || [[ $(git ls-files -o --exclude-standard) ]]; then
        local num=$(($(git ls-files -o --exclude-standard | wc -l) + $(git diff --name-only | wc -l)))
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{${integral_git_colors[2]}}%}$num${integral_git_icons[2]}"
      fi
      if ! git diff --quiet --ignore-submodules --cached 1>/dev/null 2>&1; then
        local num=$(git diff --ignore-submodules --cached --name-only | wc -l)
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{${integral_git_colors[3]}}%}$num${integral_git_icons[3]}"
      fi
      if [[ $(git remote) ]] && git cherry >/dev/null 2>&1 && [[ $(git cherry | wc -l) -gt 0 ]]; then
        local num=$(git cherry | wc -l)
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{${integral_git_colors[4]}}%}$num${integral_git_icons[4]}"
      fi
      if [[ $(git rev-list origin/${branch} --not HEAD 2>/dev/null) ]]; then
        local num=$(git rev-list origin/${branch} --not HEAD --count)
        length=$(($length + ${#num} + 2))
        format_str="$format_str %{%F{${integral_git_colors[5]}}%}$num${integral_git_icons[5]}"
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
  local dir=$(integral_dir_format)
  if [[ $1 == "w" ]]; then
    return 1
  elif [[ $1 == "r" ]]; then
    print "$dir"
  elif [[ $1 == "c" ]]; then
    print "%{%F{$integral_dir_color}%}"
  elif [[ $1 ]]; then
    print "${#dir}"
  else
    print "%{%F{$integral_dir_color}%}$dir"
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
    integral_error_format $sig
  fi
}

integral:module:jobs() {
  local -i num=$(jobs | wc -l)
  if [[ $1 ]]; then
    if [[ $num == 0 ]]; then
      print "0"
    else
      print $((${#num} + 1))
    fi
  else
    print "%{%F{$integral_jobs_color}%}$num$integral_jobs_icon"
  fi
}

integral:module:nix() {
  if [[ $1 ]]; then
    if [[ $IN_NIX_SHELL ]] || [[ $name ]]; then
      print "1"
    else
      print "0"
    fi
  else
    local color=${integral_nix_color[1]}
    [[ $IN_NIX_SHELL == "impure" ]] && color=${integral_nix_color[2]}
    if $integral_nerd_fonts; then
      print "%{%F{$color}%}${integral_nix_icons[2]}"
    else
      print "%{%F{$color}%}${integral_nix_icons[1]}"
    fi
  fi
}

integral:module:time() {
  if [[ $1 ]]; then
    print "1"
  else
    print "%{%F{12}%}$(date +${integral_time_format:-%T})"
  fi
}

integral:module:uptime() {
  local uptime=$(uptime | awk '{print $1}')
  if [[ $1 ]]; then
    print $((${#uptime} + 1))
  else
    print "%{%F{$integral_uptime_color}%}$uptime$integral_uptime_icon"
  fi
}

# BUG: leaves <space> at end of prompt
# TODO: add right prompt
#   Will require a refactor, this method will introduce complications.
#   Should create a subfunction to handle inserting newlines that inserts the right prompt.
#   Might require the entire right prompt to be rendered before the left one.
integral:loop_modules() {
  local -i position=0
  local -i max_len=$(($COLUMNS - 1))

  local newline=$'\n'
  PROMPT="$newline%{%F{$integral_prompt_color}%}${integral_prompt[1]}"

  for module in $integral_modules; do
    local -i length=$(integral:module:$module 1)
    local format_str=$(integral:module:$module)

    if [[ $length -gt 0 ]]; then
      local new_pos=$(($position + $length + 1))
      if [[ $length -gt $max_len ]] && ! integral:module:$module w; then
        local raw_str=$(integral:module:$module r)
        local color=$(integral:module:$module c)
        local -i i=0
        while [[ $i -le $((${#raw_str} / $max_len)) ]]; do
          if [[ $i == 0 ]]; then
            PROMPT+="$color${raw_str:$(($i * $max_len)):$(($max_len - $position))}"
          else
            PROMPT+="$newline%{%F{$integral_prompt_color}%}${integral_prompt[2]}$color${raw_str:$(($i * $max_len - $position)):$max_len}"
          fi
          i+=1
        done
        PROMPT+=" "
        position=$(($position + (${#raw_str} % $max_len)))
      elif [[ $new_pos -gt $max_len ]]; then
        PROMPT+="$newline%{%F{$integral_prompt_color}%}${integral_prompt[2]}$format_str "
        position=$length
      else
        PROMPT+="$format_str "
        position=$new_pos
      fi
    fi
  done
  PROMPT+="$newline%{%F{$integral_prompt_color}%}${integral_prompt[3]}%{%F{15}%}"
}
