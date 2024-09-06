# === HELPERS ===
integral:helpers:cursor-shape() {
  if [[ $1 ]]; then
    echo -ne '\e[1 q'
  else
    echo -ne '\e[5 q'
  fi
}
integral:helpers:version() {
  print "v$version"
}

