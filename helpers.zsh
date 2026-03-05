# === HELPERS ===
integral:helpers() {
  integral:helpers:$1 $@[2,-1]
}
integral:helpers:cursor-shape() {
  if [[ $1 ]]; then
    echo -ne '\e[1 q'
  else
    echo -ne '\e[5 q'
  fi
}
integral:helpers:add-prompt() {
  PROMPT+="$1"
}
integral:helpers:newline() {
  local newline=$'\n'
  if [[ $2 == "reset" ]]; then
    PROMPT="$newline%{%F{$integral_prompt_color}%}${integral_prompt[${1:-2}]}"
  else
    PROMPT+="$newline%{%F{$integral_prompt_color}%}${integral_prompt[${1:-2}]}"
  fi
}

