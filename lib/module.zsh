# === MODULES ===
# BUG: length gets printed to console in some repos
integral:module() {
  integral:module:$1 $@[2,-1]
}
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
      format_str="%F{${int_git_colors[1]}}$branch${int_git_icons[1]}"
      if ! git diff --quiet --ignore-submodules 1>/dev/null 2>&1 || [[ $(git ls-files -o --exclude-standard) ]]; then
        local num=$(($(git ls-files -o --exclude-standard | wc -l) + $(git diff --name-only | wc -l)))
        length=$(($length + ${#num} + ${#int_git_icons[2]} + 1))
        format_str="$format_str %F{${int_git_colors[2]}}$num${int_git_icons[2]}"
      fi
      if ! git diff --quiet --ignore-submodules --cached 1>/dev/null 2>&1; then
        local num=$(git diff --ignore-submodules --cached --name-only | wc -l)
        length=$(($length + ${#num} + ${#int_git_icons[3]} + 1))
        format_str="$format_str %F{${int_git_colors[3]}}$num${int_git_icons[3]}"
      fi
      if [[ $(git remote) ]] && git cherry >/dev/null 2>&1 && [[ $(git cherry | wc -l) -gt 0 ]]; then
        local num=$(git cherry | wc -l)
        length=$(($length + ${#num} + ${#int_git_icons[4]} + 1))
        format_str="$format_str %F{${int_git_colors[4]}}$num${int_git_icons[4]}"
      fi
      if [[ $(git rev-list origin/${branch} --not HEAD 2>/dev/null) ]]; then
        local num=$(git rev-list origin/${branch} --not HEAD --count)
        length=$(($length + ${#num} + ${#int_git_icons[5]} + 1))
        format_str="$format_str %F{${int_git_colors[5]}}$num${int_git_icons[5]}"
      fi
    fi
  else
    format_str="%F{11}bare⎇"
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
      format_str="%F{${int_vim_colors[1]}}${int_vim_indicators[1]}"
      ;;
    VISUAL)
      format_str="%F{${int_vim_colors[2]}}${int_vim_indicators[2]}"
      ;;
    V-LINE)
      format_str="%F{${int_vim_colors[3]}}${int_vim_indicators[3]}"
      ;;
    NORMAL)
      format_str="%F{${int_vim_colors[4]}}${int_vim_indicators[4]}"
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
  local dir=$(int_dir_format)
  if [[ $1 == "w" ]]; then
    return 1
  elif [[ $1 == "r" ]]; then
    print "$dir"
  elif [[ $1 == "c" ]]; then
    print "%F{$int_dir_color}"
  elif [[ $1 ]]; then
    print "${#dir}"
  else
    print "%F{$int_dir_color}$dir"
  fi
}

export sig=0
error_hook() {
  export sig=$?
}
integral:module:error() {
  local format_str
  if [[ $sig == 0 ]]; then
    print "0"
  else
    if [[ $1 ]]; then
      print "1"
    else
      int_error_format $sig
    fi
  fi
}

integral:module:jobs() {
  local -i num=$(jobs | wc -l)
  if [[ $num == 0 ]]; then
    print "0"
  else
    if [[ $1 ]]; then
      print $((${#num} + ${#int_jobs_icon}))
    else
      print "%F{$int_jobs_color}$num$int_jobs_icon"
    fi
  fi
}

integral:module:nix() {
  if [[ $IN_NIX_SHELL ]] || [[ $name ]]; then
    if [[ $1 ]]; then
      print "1"
    else
      local color=${int_nix_color[1]}
      [[ $IN_NIX_SHELL == "impure" ]] && color=${int_nix_color[2]}
      if $int_nerd_fonts; then
        print "%F{$color}${int_nix_icons[2]}"
      else
        print "%F{$color}${int_nix_icons[1]}"
      fi
    fi
  else
    print "0"
  fi
}

integral:module:time() {
  if [[ $1 ]]; then
    print "1"
  else
    print "%F{12}$(date +${int_time_format:-%T})"
  fi
}

integral:module:uptime() {
  local uptime=$(uptime | awk '{print $1}')
  if [[ $1 ]]; then
    print $((${#uptime} + ${#int_uptime_icon}))
  else
    print "%F{$int_uptime_color}$uptime$int_uptime_icon"
  fi
}

integral:module:battery() {
  local format_str=$(cat /sys/class/power_supply/*([1])/capacity 2>/dev/null)
  if [[ $1 ]] && [[ $format_str ]]; then
    print "$((${#format_str} + ${#int_battery_icon}))"
  elif [[ $1 ]]; then
    print "0"
  else
    print "%F{$int_battery_color}$format_str$int_battery_icon"
  fi
}

integral:module:ssh() {
  if [[ $SSH_CONNECTION ]]; then
    local format_str=$(int_ssh_format || print "%F{13}${USER}%F{8}@%F{14}${HOSTNAME}" )
    if [[ $1 ]]; then
      print ${#format_str};
    else
      print "$format_str"
    fi
  else
    print "0"
  fi
}

integral:module:direnv() {
  if [[ $DIRENV_DIR ]]; then
    if [[ $1 ]]; then
      print "1"
    else
      print "$(int_direnv_format $DIRENV_DIR)"
    fi
  else
    print "0"
  fi
}

integral:module:cpu() {
  local cpu="$(top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4}')%"
  if [[ $1 ]]; then
    print "${#cpu}"
  else
    print "%F{${int_cpu_color}}${cpu}"
  fi
}

integral:module:distrobox() {
  if [[ $CONTAINER_ID ]]; then
    if [[ $1 ]]; then
      print ${#CONTAINER_ID}
    else
      print "%F{${int_distrobox_color}}${CONTAINER_ID}"
    fi
  else
    print "0"
  fi
}

integral:module:sshplus() {
  local ssh=$(integral module ssh)
  local db=$(integral module distrobox)
  local format_str
  if [[ $ssh != 0 ]] && [[ $db != 0 ]]; then
    format_str="$ssh%F{8}[$db%F{8}]"
  elif [[ $ssh != 0 ]]; then
    format_str=$ssh
  elif [[ $db != 0 ]]; then
    format_str=$db
  fi
  if [[ $format_str ]]; then
    if [[ $1 ]]; then
      print ${#format_str}
    else
      print $format_str
    fi
  else
    print "0"
  fi
}

# BUG: leaves <space> at end of prompt
# TODO: add right prompt
#   Will require a refactor, this method will introduce complications.
#   Should create a subfunction to handle inserting newlines that inserts the right prompt.
#   Might require the entire right prompt to be rendered before the left one.
integral:render() {
  local -i position=0
  local -i max_len
  if [[ $1 ]]; then
    max_len=$1
  else
    max_len=$(($COLUMNS - 1))
  fi

  if [[ ${#int_separator} -gt 1 ]]; then echo "Unsupported separator size!"; fi

  integral helpers newline 1 reset

  for module in $int_modules; do
    [[ $1 ]] && print -P "%F{15}module: $module"
    local -i length=$(integral module $module 1)
    local format_str=$(integral module $module)

    if [[ $length -gt 0 ]]; then
      local new_pos=$(($position + $length + 1))
      if [[ $length -gt $max_len ]] && ! integral module $module w; then
        [[ $1 ]] && print -P "%F{15}wrapping!"
        local raw_str=$(integral module $module r)
        local color=$(integral module $module c)
        local -i i=0
        while [[ $i -le $((${#raw_str} / $max_len)) ]]; do
          if [[ $i == 0 ]]; then
            integral helpers add-prompt "$color${raw_str:$(($i * $max_len)):$(($max_len - $position))}" $1
          else
            integral helpers newline
            integral helpers add-prompt "$color${raw_str:$(($i * $max_len - $position)):$max_len}" $1
          fi
          i+=1
        done
        integral helpers add-prompt "${int_separator}" $1
        position=$(($position + (${#raw_str} % $max_len)))
      elif [[ $new_pos -gt $max_len ]]; then
        [[ $1 ]] && print -P "%F{15}new line"
        integral helpers newline
        integral helpers add-prompt "$format_str${int_separator}" $1
        position=$length
      else
        integral helpers add-prompt "$format_str${int_separator}" $1
        position=$new_pos
      fi
    elif [[ $1 ]]; then
      print -P "%F{15}skipping"
    fi
  done
  integral helpers newline 3
  integral helpers add-prompt "%F{15}" $1
  [[ $1 ]] && print "====="
}
