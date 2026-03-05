local version='0.1.0'

integral() {
  integral:$1 $@[2,-1]
}

integral:--version() {
  print "v$version"
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

# === INIT ===
TRAPWINCH() {
  integral loop-modules
  zle && zle reset-prompt
}
add-zsh-hook precmd error_hook
add-zsh-hook precmd integral:loop-modules
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral:loop-modules
zle -N zle-line-init integral:zle-line-init

