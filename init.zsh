local version='0.2.0'

integral() {
  case $1 in
    --version)
      print "v$version"
      ;;
    -h|--help)
      cat <<EOF
integral [ group ... ] <function> [ args ... ]

$(print -l ${(ok)functions} | rg -e '^integral:' | sed s/integral:// | sed 's/:/ /')
EOF
      ;;
    *)
      integral:$1 $@[2,-1]
      ;;
  esac
}

# https://github.com/spaceship-prompt/spaceship-prompt/commit/111c6f160c4376001d5469f8e8771ee89ea4158a
local int_path=${${(%):-%x}:A:h}
export integral_plugins=(
  "$int_path/config.zsh"
  "$int_path/helpers.zsh"
  "$int_path/module.zsh"
  "$int_path/zle.zsh"
)
for f in $integral_plugins; do
  if [[ -f $f ]]; then
    source $f
  else
    print "Plugin not found: $f"
    exit 1
  fi
done
fpath+="$int_path/comp"
autoload -Uz add-zsh-hook
autoload -U add-zle-hook-widget

export VI_KEYMAP=${VI_KEYMAP:-"INSERT"}
export HOSTNAME=${HOSTNAME:-"$(hostname)"}

# === INIT ===
TRAPWINCH() {
  integral render
  zle && zle reset-prompt
}
add-zsh-hook precmd error_hook
add-zsh-hook precmd integral:render
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral render
zle -N zle-line-init integral:zle-line-init

