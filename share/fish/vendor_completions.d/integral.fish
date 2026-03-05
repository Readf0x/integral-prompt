complete -c integral -f

# Subcommands
complete -c integral -n '__fish_use_subcommand' -a transient -d 'render transient prompt'
complete -c integral -n '__fish_use_subcommand' -a render -d 'render full prompt'
complete -c integral -n '__fish_use_subcommand' -a init -d 'print shell init script'
complete -c integral -n '__fish_use_subcommand' -a version -d 'print version info'
complete -c integral -n '__fish_use_subcommand' -a config -d 'print config as json'
complete -c integral -n '__fish_use_subcommand' -a help -d 'show help message'

# ---- init ----
complete -c integral -n '__fish_seen_subcommand_from init; and not __fish_seen_argument 1' -a 'raw bash zsh fish' -d 'shell'

# ---- render ----
# arg 1: shell
complete -c integral -n '__fish_seen_subcommand_from render; and not __fish_seen_argument 1' -a 'raw bash zsh fish' -d 'shell'
# arg 2: number of columns
complete -c integral -n '__fish_seen_subcommand_from render; and __fish_seen_argument 1; and not __fish_seen_argument 2' -a "$COLUMNS" -d 'number of columns'
# arg 3: exit signal
complete -c integral -n '__fish_seen_subcommand_from render; and __fish_seen_argument 2; and not __fish_seen_argument 3' -a '0 1 2 126 127 128 130 137 139 143 255' -d 'exit signal (decimal)'
# arg 4: number of background jobs
complete -c integral -n '__fish_seen_subcommand_from render; and __fish_seen_argument 3; and not __fish_seen_argument 4' -a "(jobs | grep '^\[' | wc -l)" -d 'number of background jobs'

# ---- commands with no args ----
for sub in transient version config help
    complete -c integral -n "__fish_seen_subcommand_from $sub" -d 'no additional arguments'
end
