echo -ne '\e[5 q'

__integral_render() {
  eval "$(integral render bash $COLUMNS $? $(jobs | wc -l))"
}

# Unfortunately just setting PS1 doesn't actually update the prompt
# __integral_debug() {
#   PS1=$(integral transient)
# }
#
# trap '__integral_debug' DEBUG

PROMPT_COMMAND="__integral_render; $PROMPT_COMMAND"

