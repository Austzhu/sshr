package main

import (
	"fmt"
)

var bash = `
__dbg()
{
	[ "${SSHRCOMPLETE}" != "yes" ] && return
	echo -e "-----"$@"-----"
}

_complete_scp()
{
	local tmp=($(cat ~/.sshr/sshr.conf | grep -v "^[ \t]*#" | awk '{print $1}'))
	local uhost=()
	for p in ${tmp[@]}; do
		uhost=(${uhost[@]} $p":")
	done

	# 将输入按照空格切分
	OIFS=${IFS}
	IFS=' '
	read -ra line <<<${COMP_LINE}
	IFS=${OIFS}

	#
	local c=${#line[@]}
	local dir=${line[$c-1]##*:}
	[ $c == 2 ] && dir=""
	local views=(${uhost[@]} $(compgen -o default -- ${dir}))

	case $c in
		2)	COMPREPLY=(${views[@]}) ;;
		3)
			if [ "${COMP_LINE: -1}" == " " ]; then
				if [ x$(echo ${line[2]} | grep "@") == x ]; then
					COMPREPLY=(${uhost[@]} $(compgen -o default))
				else
					COMPREPLY=($(compgen -o default))
				fi
				return
			fi

			for p in ${views[@]} ; do
				if [ x$(echo $p | grep "@") != x ]; then
					if [[ ${line[$c-1]} =~ ^"$p" ]]; then
						COMPREPLY=($(compgen -o default -- ${dir}))
						return
					fi
				fi

				[[ $p =~ ^"${line[$c-1]}" ]] && COMPREPLY=(${COMPREPLY[@]} $p)
			done
		;;

		4)
			[ "${COMP_LINE: -1}" == " " ] && return
			for p in ${views[@]} ; do
				if [ x$(echo $p | grep "@") != x ]; then
					if [[ ${line[$c-1]} =~ ^"$p" ]]; then
						COMPREPLY=($(compgen -o default -- ${dir}))
						return
					fi
				fi

				[[ $p =~ ^"${line[$c-1]}" ]] && COMPREPLY=(${COMPREPLY[@]} $p)
			done
		;;

	esac
}

_complete_sshr()
{
	local cmds=(%s)
	local uhost=($(cat ~/.sshr/sshr.conf | grep -v "^[ \t]*#" | awk '{print $1}'))
	local views=(${cmds[@]} ${uhost[@]})

	COMPREPLY=()

	__dbg ${COMP_CWORD} ${#COMP_WORDS[@]} ${COMP_WORDS[@]}
	__dbg "\$1=$1,\$2=$2,\$3=$3"

    case "$3" in
    "$1")
        [ "$2" ] || { COMPREPLY=(${views[@]}) && return; }
		for p in ${views[@]}; do
			[[ $p =~ ^"$2" ]] && COMPREPLY=(${COMPREPLY[@]} $p)
		done
		return
        ;;
    esac

	# scp相关的complete
	if [ "${COMP_WORDS[1]}" == "scp" ]; then
		_complete_scp $@
		return
	fi
}

complete  -F _complete_sshr sshr
`

func doComplete(c *cmd) error {
	cmds := ""
	for _, v := range sCmd.c {
		if v.name != SSHR {
			cmds += " " + v.name
		}
	}

	fmt.Printf(bash, cmds[1:])
	return nil
}

func init() {
	name := "complete"
	register(&cmd{
		name:    name,
		brief:   "生成bash的TAB补全代码",
		flagSet: nil,
		cb:      doComplete,
	})
}
