source <(integral init zsh)
fpath=($PWD/share/zsh/site-functions $fpath)
autoload -U compinit && compinit
