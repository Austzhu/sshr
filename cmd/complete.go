package main

import (
	"fmt"
)

var bash = `
__dbg()
{
	[ "${SSHRCOMPLETE}" != "yes" ] && return
	echo -e "-----"$@"-----\n"
}

_conplete_sshr()
{
	local cmds=(%s)
	local uhost=$(cat ~/.sshr/sshr.conf | grep -v "^[ \t]*#" | awk '{print $1}')
	local views=(${cmds[@]} ${uhost})

	__dbg ${#views[@]}  ${views[@]}
	__dbg $1"--"$2"--"$3

    COMPREPLY=()
    case "$3" in
    "$1")
        [ "$2" ] || { COMPREPLY=(${views[@]}) && return; }

        for ((i = 0; i < ${#views[@]}; i++)); do
			[[ ${views[i]} =~ ^"$2" ]] && COMPREPLY=(${COMPREPLY[@]} ${views[i]})
		done
        ;;
    esac
}

complete -F _conplete_sshr sshr
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
