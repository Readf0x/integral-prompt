function fish_prompt
	set sig $status
	if set -q int_transient
		echo -ne '\r\e[J'
		integral transient
		set -e int_transient
	else
		integral render fish $COLUMNS $sig $(jobs | wc -l) | source
	end
end
function fish_right_prompt
	echo "$integral_right_prompt_string"
end

function __int_resize --on-signal WINCH
	commandline -f repaint
end

function __int_transient_execute
	commandline -f expand-abbr suppress-autosuggestion
	if commandline --is-valid || test -z "$(commandline)"
		if commandline --paging-mode && test -n "$(commandline)"
			commandline -f accept-autosuggestion
			return 0
		end

		set --global int_transient 1
		commandline -f repaint execute
		return 0
	end

	commandline -f execute
end

bind --user --mode default \r __int_transient_execute
bind --user --mode insert \r __int_transient_execute
bind --user --mode default \cj __int_transient_execute
bind --user --mode insert \cj __int_transient_execute
echo -ne '\e[5 q'
