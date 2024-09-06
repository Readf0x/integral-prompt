local version='0.0.5'

export integral_plugins=(
  './config.zsh'
  './helpers.zsh'
  './module_loader.zsh'
  './zle.zsh'
)
for f in $integral_plugins; do
  source $f
done
autoload -U colors; colors
autoload -Uz add-zsh-hook
autoload -U add-zle-hook-widget

export VI_KEYMAP=${VI_KEYMAP:-"INSERT"}

integral:prompt() {
  integral:loop_modules
}

# === INIT ===
TRAPWINCH() {
  integral:prompt
  zle reset-prompt
}
add-zsh-hook preexec error_hook
add-zsh-hook precmd integral:prompt
add-zsh-hook precmd integral:helpers:cursor-shape
zle -N integral:line-pre-redraw
add-zle-hook-widget zle-line-pre-redraw integral:line-pre-redraw
integral:prompt
zle -N zle-line-init integral:zle-line-init

