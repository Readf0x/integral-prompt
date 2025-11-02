function fish_prompt
	integral render fish $COLUMNS $status $(jobs | wc -l) | source
end
function fish_right_prompt
	echo "$integral_right_prompt_string"
end
