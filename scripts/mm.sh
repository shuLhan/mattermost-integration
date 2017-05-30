#!/bin/sh
##
## mm.sh: script to send any text to Mattermost.
##
## Parameters,
##
##  1: error level in string "INFO | ERROR", default to INFO
##  2: message to send
##
## NOTE:
## - message will be truncated to 4000 characters.
## - use double quote to sending message with space
##
## Example:
##
##	$ mm.sh ERROR "This is an error message"
##	$ cat log | mm.sh
##

## For testing
#MM_URL="https://your.server.com/hooks/t934mqsdm3rxpytuhzfearyifc"
#MM_CHANNEL="log_test"
#MM_BOTNAME="test"

MM_URL="{{ chat_url }}"
MM_CHANNEL="{{ chat_channel }}"
MM_BOTNAME="{{ shorthostname }}"

MM_ICON=":large_blue_circle:"
MM_HEADER=""
MM_TEXT=""

## Read pipe from STDIN.

if [ ! -t 0 ]; then
	echo ">> Reading from STDIN $0"
	MM_TEXT=$(cat -)
fi

while [ $# != 0 ]; do
	case "$#" in
	2)
		if [ $1 == "ERROR" ]; then
			MM_ICON=":red_circle:";
		fi
		;;
	1)
		if [ "${MM_TEXT}" == "" ]; then
			MM_HEADER="$1"
		else
			MM_HEADER="**$1**"
		fi
		;;
	esac

	shift
done

## Escape backslash char
MM_HEADER=${MM_HEADER//\\/\\\\}
MM_TEXT=${MM_TEXT//\\/\\\\}

## Escape quote char
MM_HEADER=${MM_HEADER//\"/\\\"}
MM_TEXT=${MM_TEXT//\"/\\\"}

MM_PAYLOAD="{\
\"username\":\"${MM_BOTNAME}\"\
,\"channel\":\"${MM_CHANNEL}\"\
,\"text\":\"${MM_ICON} ${MM_HEADER}\n${MM_TEXT}\"\
}"

## Cut text to 4000 characters.
MM_PAYLOAD=${MM_PAYLOAD:0:4000}

echo ">> Payload: ${MM_PAYLOAD}"

curl --insecure --silent --show-error \
	--header "Content-type: application/json" \
	--data "${MM_PAYLOAD}" \
	-X POST "${MM_URL}"
