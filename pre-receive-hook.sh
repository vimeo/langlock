#!/bin/bash

echo "Server-side secret scanner (Langlock) has been configured in this instance of GitHub Enterprise. If needed, you may bypass scanning of your pushes with the following command: git push --push-option=SKIP_LANGLOCK"

# ------------------------------------------------------------------------------
# Debug
# ------------------------------------------------------------------------------

function debug() {
  [[ "$DEBUG" = 'true' ]] && echo -e "DEBUG $1" 1>&2
}

# ------------------------------------------------------------------------------
# Configuration
# ------------------------------------------------------------------------------

# Defaults
# SLACK_WEBHOOK_URL="$(cat /home/githook/slack_webhook_url.txt)"
SLACK_WEBHOOK_URL="$(/home/githook/slack_webhook_url.txt)"
#SLACK_USERS_TO_ALERT='<@UT171F66B>'
SLACK_USERS_TO_ALERT=''
GRACEFUL_TIMEOUT_SECS=3.2
BLOCK_PUSH_ON_TIMEOUT='false'
ALLOW_LIST_PATH='langlock.json'
ENABLE_OPT_IN='false'
OPT_IN_FLAG='SECRET_SCAN'
OPT_OUT_FLAG='SKIP_LANGLOCK'
DEBUG='true'
PROFILE='true'
NUM_THREADS=4
ENROLLED_USERS=""
USERS_WITH_COLOR_ENABLED=""
COLOR='true' # Many developers' IDEs and GUI tools with git integrations do not support color,
              # so in those contexts, the color codes (eg. `\e[31m`) will be printed directly

LANGLOCK_CONFIG_PATH="$(dirname "${BASH_SOURCE[0]}")/config.json"
readonly LANGLOCK_CONFIG_PATH
if [[ -f "$LANGLOCK_CONFIG_PATH" ]]; then
  for s in $(cat "$LANGLOCK_CONFIG_PATH" | jq -r "to_entries|map(\"\(.key)=\(.value|tostring)\")|.[]" ); do
    export $s;
  done
fi

readonly SLACK_WEBHOOK_URL
readonly GRACEFUL_TIMEOUT_SECS
readonly BLOCK_PUSH_ON_TIMEOUT
readonly ALLOW_LIST_PATH
readonly ENABLE_OPT_IN
readonly OPT_IN_FLAG
readonly OPT_OUT_FLAG
readonly DEBUG
readonly PROFILE
readonly NUM_THREADS
readonly ENROLLED_USERS
readonly USERS_WITH_COLOR_ENABLED
readonly COLOR

#debug "GITHUB_USER_LOGIN $GITHUB_USER_LOGIN"
#debug "GIT_DIR $GIT_DIR"
#debug "GITHUB_USER_IP $GITHUB_USER_IP"
#debug "GITHUB_REPO_NAME $GITHUB_REPO_NAME"
#debug "GITHUB_PULL_REQUEST_AUTHOR_LOGIN $GITHUB_PULL_REQUEST_AUTHOR_LOGIN"
#debug "GITHUB_REPO_PUBLIC $GITHUB_REPO_PUBLIC"
#debug "GITHUB_PUBLIC_KEY_FINGERPRINT $GITHUB_PUBLIC_KEY_FINGERPRINT"
#debug "GITHUB_PULL_REQUEST_HEAD $GITHUB_PULL_REQUEST_HEAD"
#debug "GITHUB_PULL_REQUEST_BASE $GITHUB_PULL_REQUEST_BASE"
#debug "GITHUB_VIA $GITHUB_VIA"
#debug "GIT_PUSH_OPTION_COUNT $GIT_PUSH_OPTION_COUNT"
HOOK_ENV="$(cat /home/githook/pre_receive_hook_env_version.txt)"
#debug "HOOK_ENV $HOOK_ENV"
#debug "SLACK_WEBHOOK_URL $SLACK_WEBHOOK_URL"

# ------------------------------------------------------------------------------
# Opt-in
# ------------------------------------------------------------------------------

if [ $ENABLE_OPT_IN = 'true' ]; then
  if [[ ! ( "$GIT_PUSH_OPTION_0" == "$OPT_IN_FLAG" || "$GIT_PUSH_OPTION_1" == "$OPT_IN_FLAG" ) ]]; then
    # Did not opt in to the secret scanning, so immediately returning success
    exit 0
  fi
fi


# ------------------------------------------------------------------------------
# Determine repo name and user name
# ------------------------------------------------------------------------------

# Github Enterprise sets the value of GITHUB_REPO_NAME, but other pre-receive
# environments do not
if [ ! -z "$GITHUB_REPO_NAME" ]; then
  REPO_NAME=$GITHUB_REPO_NAME
elif [ $(git rev-parse --is-bare-repository) = true ]
then
  REPO_NAME=$(basename "$PWD")
else
  REPO_NAME=$(basename $(readlink -nf "$PWD"/..))
fi
readonly REPO_NAME

UNKNOWN_USER_IDENTIFIER='UNKNOWN_USER'
if [ -z "$GITHUB_USER_LOGIN" ]; then
  GITHUB_USER_LOGIN="$UNKNOWN_USER_IDENTIFIER"
fi
readonly GITHUB_USER_LOGIN


# ------------------------------------------------------------------------------
# Bypass if not enrolled
# ------------------------------------------------------------------------------

# If the push was not made by a user, then skip all scanning (exit script with success integer code)
if [[ ( "$GITHUB_USER_LOGIN" == "$UNKNOWN_USER_IDENTIFIER" ) ]]; then
  exit 0
fi

# If a list of ENROLLED_USERS is specified in the config file AND the
# user who is making the current push is not one of the enrolled
# users, then skip all scanning (exit script with success integer code))
if [ ! -z "$ENROLLED_USERS" ]; then
  echo "$ENROLLED_USERS" | grep "\"$GITHUB_USER_LOGIN\"" > /dev/null
  if [ ! $? -eq 0 ]; then
    debug "Skip scanning because user not enrolled for Langlock"
    exit 0
  fi
  # Check allow and deny list. If both provided, only respect allow list. If
  # neither an allow list nor a deny list is provided, treat it as "allow all."
  ENABLED_REPOS_FOR_USER=$(jq ".$GITHUB_USER_LOGIN.scanRepos" <(echo $ENROLLED_USERS))
  DISABLED_REPOS_FOR_USER=$(jq ".$GITHUB_USER_LOGIN.skipRepos" <(echo $ENROLLED_USERS))
  if [[ ( "$ENABLED_REPOS_FOR_USER" != "null" ) ]]; then
    echo "$ENABLED_REPOS_FOR_USER" | grep "$REPO_NAME" > /dev/null
    if [ ! $? -eq 0 ]; then
      debug "Skip scanning because repo not on user's allow list"
      exit 0
    fi
  elif [[ ( "$DISABLED_REPOS_FOR_USER" != "null" ) ]]; then
    echo "$DISABLED_REPOS_FOR_USER" | grep "$REPO_NAME"
    if [ $? -eq 0 ]; then
      debug "Skip scanning because repo blocked by user's deny list"
      exit 0
    fi
  else : ;
  fi
fi

# ------------------------------------------------------------------------------
# Enable or disable color
# ------------------------------------------------------------------------------
EFFECTIVE_COLOR="false"
if [[ ! -z "$COLOR" && $COLOR = 'true' && -z "$USERS_WITH_COLOR_ENABLED" ]]; then
  EFFECTIVE_COLOR="true"
  #debug "Color enabled for all users"
elif [[ ( -z "$COLOR" || $COLOR = 'true' ) && ! -z "$USERS_WITH_COLOR_ENABLED" ]]; then
    echo "$USERS_WITH_COLOR_ENABLED" | grep "\"$GITHUB_USER_LOGIN\"" > /dev/null
    if [ $? -eq 0 ]; then
      EFFECTIVE_COLOR="true"
      #debug "Color enabled for user $GITHUB_USER_LOGIN"
    fi
fi
readonly EFFECTIVE_COLOR


#if [[ ( -z "$GITHUB_USER_LOGIN" ||  ( "$GITHUB_USER_LOGIN" != "ed" && "$GITHUB_USER_LOGIN" != "edward-sullivan" ) ) ]]; then
#  debug "Exiting because user is not Ed"
#  exit 0
#fi

#if [[ -f "./README.md" ]]; then
#    echo -e "\n\nREADME exists.\n\n"
#else
 #   echo '.'
    #echo -e "\n\nNOT FOUND: $(pwd)\n"
    # echo -e "Trying named pipes"
    # mkfifo testfifo
    # ls > testfifo &
    # cat testfifo
    # echo -e "done\n"
    # ls -la /data/repositories/c/nw/c4/ca/42/1
    # ls -la /data/repositories/c/nw/c4/ca/42
    # echo -e "whoami: $(whoami)\n"
    #echo -e "About to try cat /home/githook/test2.txt"
    # cat /home/githook/test2.txt
    # echo -e "About to try cat /test.txt"
    # cat /test.txt
    # echo -e "\nls -la /home/githook: $(ls -la /home/githook)\n"
    # echo -e "\ncat /home/githook/slack_webhook_url.txt: $(cat /home/githook/slack_webhook_url.txt)\n"
    # whoami
    # # chown githook:root /home/githook/test.txt
    # cat test.txt
    # cat /home/githook/test.txt
    # echo -e "su nobody then cat\n"
    # su nobody
    # cat /home/githook/test123.txt
    # cat /home/githook/test.txt
    # echo -e "Trying process substitution\n"
    # wc -c <(echo "abc") # should print '4'
    # echo -e "Trying to run command\n"
    # /home/githook/detectcreds.bin
    # echo -e "Trying to run test\n"
    # /home/githook/test.txt
    # echo -e "ls -la /\n"
    # ls -la /
    # echo -e "ls -la /dev\n"
    # ls -la /dev
    # echo -e "ls -la /home\n"
    # ls -la /home
    # echo -e "\ncat /home/githook/test.txt: $(cat /home/githook/test.txt)\n"
    # echo -e "\nsudo cat /home/githook/slack_webhook_url.txt: $(sudo cat /home/githook/slack_webhook_url.txt)\n"
    #echo -e "LS -la: $(ls -la)\n"
    # echo -e "LS refs: $(ls refs)\n"
    # echo -e "CAT Description: $(cat description)\n"
    #echo -e "Script dir: $(dirname "${BASH_SOURCE[0]}")\n"
    #echo -e "LS in script dir: $(ls $(dirname "${BASH_SOURCE[0]}"))\n"
    #echo -e "\n"
#fi




# ------------------------------------------------------------------------------
# Generate a unique id for the push (just to associate log messages)
# ------------------------------------------------------------------------------

function gen_rand_len32_base62_id {
  echo $(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z1-9' | fold -w 32 | head -n 1)
}

PUSH_ID=$(gen_rand_len32_base62_id)
readonly PUSH_ID




# ------------------------------------------------------------------------------
# Profile utils
# ------------------------------------------------------------------------------

function current_millisecond_timestamp {
  echo $(python3 -c 'import time; print(int(round(time.time() * 1000)))')
}

PROFILE_START_TIME=$(current_millisecond_timestamp)
readonly PROFILE_START_TIME

function elapsed_milliseconds {
   local current_time=$(current_millisecond_timestamp)
   echo $(($current_time - $PROFILE_START_TIME))
}

function profile() {
  [[ $PROFILE = 'true' ]] && echo "PROFILE $1: $(elapsed_milliseconds) ms" 1>&2
}

# ------------------------------------------------------------------------------
# Git utils
# ------------------------------------------------------------------------------

# Commit sha with all zeros
ZERO_COMMIT='0000000000000000000000000000000000000000'
readonly ZERO_COMMIT

# ------------------------------------------------------------------------------
# Misc util
# ------------------------------------------------------------------------------

# NOTE: For performance reasons, I have made the log() function async, since a round
#       trip to an outside service may be slow. However, there
#       is a subtle and critical pitfall to be mindful of when using asynchronous
#       subprocesses. Specifically, if one of the subprocesses finishes after the main
#       parent process finishes, then the script's ultimate exit status code will
#       the exit code of the last subprocess to finish, instead of the explicit
#       exit code of the parent process, which is not acceptable, because the
#       exit code of the script determines whether or not the push is accepted
#       upstream. Hence we want to make sure that always set the exit status
#       explicitly and intentionally.
#       .....

async_subprocesses=()

function log {
  #local message="$1"
  local message=$(echo "$1" | sed 's/"/\\"/g')
  local slack_users_to_alert="$2"
  # curl -X POST -H 'Content-type: application/json' --max-time 0.7 --data "{\"text\":\"$message ($(elapsed_milliseconds) ms, Push ID: $PUSH_ID) $slack_users_to_alert\"}" $SLACK_WEBHOOK_URL &
  curl -X POST -H 'Content-type: application/json' --max-time 0.7 --data "{\"text\":\"$message ($(elapsed_milliseconds) ms, Push ID: $PUSH_ID) $slack_users_to_alert\"}" $SLACK_WEBHOOK_URL 2>/dev/null >/dev/null &
  async_subprocesses+=( $! )
  #debug "IN log FUNC async_subprocesses[*]: ${async_subprocesses[*]} $1"
}

function wait_for_async_subprocesses {
  #debug "IN wait FUNC async_subprocesses[*]: ${async_subprocesses[*]}"
  #debug "IN wait FUNC async_subprocesses[-1]: ${async_subprocesses[-1]}"
  wait ${async_subprocesses[-1]}
  #wait ${async_subprocesses[*]}
}

# Allows joining elements in an array in a single string, with a multi-character separator
function join {
   local separator
   local result
   local first

   separator="$1"
   local -n arr=$2
   result=""

   first=true
   for i in "${arr[@]}"
   do
      if [ $first = 'true' ]; then
        result="$i"
        first='false'
      else
        result="$result$separator$i"
      fi
   done
   echo $result
}


# ------------------------------------------------------------------------------
# Output styling
# ------------------------------------------------------------------------------

RED='\033[31m'
GREEN='\033[92m'
BLUE='\033[96m'
BOLD='\033[1m'
PLAIN='\033[0m'
readonly RED
readonly GREEN
readonly BLUE
readonly BOLD
readonly PLAIN

# Note that the trailing whitespace in the HEADER is necessary
HEADER=$(cat <<-END

==========================================================================

_-_-                      ,,            ,,
 /,       _           _   ||            ||
 ||      < \, \\\\\\/\\\\\\  / \\\\\\ ||  /'\\\\\\  _-_ ||/\ 
~||      /-|| || || || || || || || ||   ||_<
 ||     (( || || || || || || || || ||   || |
(  -__,  \/\\\\\\ \\\\\\ \\\\\\ \\\\\\_-| \\\\\\ \\\\\\,/  \\\\\\,/ \\\\\\,\ 
                     /  \ 
                    '----'   From the Vimeo security team

Langlock: a server-side pre-receive hook to prevent developers from
accidentally pushing hardcoded secrets upstream.
==========================================================================
END
)
readonly HEADER

FOOTER=$(cat <<-END
If you have questions or need help, please contact your security team or Github administrators.
 
END
)
readonly FOOTER

ALLOW_LIST_FORMAT_MESSAGE=$(cat <<-END
Allow list must be JSON, with the format:

   {
      "allowStrings": {
         "hashOfSecret" : {<metadata>},
         "hashOfAnotherSecret": {<metadata>},
         ...
      },
      "allowPaths": [
         {
            "regex": "<Perl-like regex to specify a file name/path>",
            "reason": "<Justification for this entry>"
         },
         ...
      ]
   }

   Please note that plaintext secrets (even if they're false positives)
   should not appear in the allow list, only their hashes should.
END
)
readonly ALLOW_LIST_FORMAT_MESSAGE

BLOCKED_TIMEOUT_MESSAGE=$(cat <<-END
ERROR: The pre-receive hook that scans pushes for hardcoded secrets took too long
analyzing your push! Try pushing again, or to bypass the secret scanning,
use the option --push-option=$OPT_OUT_FLAG with your push command.
END
)
if [ $EFFECTIVE_COLOR = 'true' ]; then
  BLOCKED_TIMEOUT_MESSAGE=$(echo -e "$BLOCKED_TIMEOUT_MESSAGE" | awk -v bold=$BOLD -v red=$RED -v plain=$PLAIN '{print bold red $0 plain}')
fi
readonly BLOCKED_TIMEOUT_MESSAGE


# ------------------------------------------------------------------------------
# Parse input to pre-receive hook
# ------------------------------------------------------------------------------

#profile 'start-parsing-hook-input'

oldrevs=()
newrevs=()
refnames=()
while read oldrev newrev refname; do
  oldrevs+=( $oldrev )
  newrevs+=( $newrev )
  refnames+=( $refname )
done
OLDREV_STRINGS=$(join ', ' oldrevs)
NEWREV_STRINGS=$(join ', ' newrevs)
BRANCH_NAMES=$(join ', ' refnames)
readonly OLDREV_STRINGS
readonly NEWREV_STRINGS
readonly BRANCH_NAMES

#profile 'done-parsing-hook-input'




# ------------------------------------------------------------------------------
# Log that the pre-receive hook has begun
# ------------------------------------------------------------------------------

log "*Pre-receive hook started* for user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\` for branches \`$BRANCH_NAMES\` (old revs: \`$OLDREV_STRINGS\`; new revs: \`$NEWREV_STRINGS\`; github_via: \`$GITHUB_VIA\`; hook env version: \`$HOOK_ENV\`)."
INITIAL_LOG_ID=$!
readonly INITIAL_LOG_ID


# ------------------------------------------------------------------------------
# Skip scanning if we are merging a PR (and all branches are already in cloud)
# ------------------------------------------------------------------------------

# GITHUB_VIA will have a value of "pull request merge button" when a user attempts
# to merge a PR via the browser. A regular push via git cli or IDE will have
# an empty/undefined value. Directly editing, commiting, and pushing from
# github within the browser will have a value of "blob edit"
if [[ ( ! -z "$GITHUB_VIA" && "$GITHUB_VIA" != "blob edit" ) ]]; then
  #wait $INITIAL_LOG_ID
  log "*Push succeeeded, scanning skipped for upstream merge* for user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\`."
  wait_for_async_subprocesses
  exit 0
fi


# ------------------------------------------------------------------------------
# Bypass / Opt-out (will notify slack)
# ------------------------------------------------------------------------------

if [[ ( "$GIT_PUSH_OPTION_0" == "$OPT_OUT_FLAG" || "$GIT_PUSH_OPTION_1" == "$OPT_OUT_FLAG" ) ]]; then

  # Git user explicitly bypassing the secret scanning with a --push-option

  #wait $INITIAL_LOG_ID
  log "*Bypass applied* by user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\` to branches \`$BRANCH_NAMES\`." "$SLACK_USERS_TO_ALERT"

  echo "Server-side secret scanner (Langlock) bypassed."
  wait_for_async_subprocesses
  exit 0
fi


# ------------------------------------------------------------------------------
# Start timer to trigger graceful timeout, to avoid the abrupt external-caused timeout
# ------------------------------------------------------------------------------

#profile 'start-setup-graceful-timer'

# I'll refer to this backgrounded subshell as the "sleep wrapper", since it contains
# the sleep
MAIN_PROC=$$
trap "profile 'start-executing-trap'; exit 0" SIGINT SIGQUIT SIGTERM
{
  sleep $GRACEFUL_TIMEOUT_SECS
  #profile 'done-sleep-start-perform-graceful-timeout-post-sleep'
  #debug "ps 0"
  #ps
  #kill -s SIGINT -1
  #kill -s SIGINT -1
  #kill $$
  #debug "ps 1"
  #ps
  if [ $BLOCK_PUSH_ON_TIMEOUT = 'true' ]; then
    log "*Pre-receive hook timed out* for user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\`. *PUSH BLOCKED.*" "$SLACK_USERS_TO_ALERT"
    echo -e "$HEADER\n\n$BLOCKED_TIMEOUT_MESSAGE\n\n$FOOTER"
    #profile 'done-perform-graceful-timeout-post-sleep'
    # Note: not necessary to wait for async subproceeses since we terminate them above
    wait_for_async_subprocesses 2>/dev/null
    kill -s SIGINT -1
    #kill -10 $MAIN_PROC
    #exit 1
  fi
  log "*Pre-receive hook timed out* for user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\`. *PUSH ACCEPTED.*" "$SLACK_USERS_TO_ALERT"
  echo "Server-side secret scanner (Langlock) timed out, but git push still accepted."
  #profile 'done-perform-graceful-timeout-post-sleep'
  #wait $INITIAL_LOG_ID
  # Note: not necessary to wait for async subproceeses since we terminate them above
  #wait_for_async_subprocesses 2>/dev/null
  wait_for_async_subprocesses
  #debug "MAIN_PROC: $MAIN_PROC"
  #debug "ps"
  #ps
  kill -s SIGINT -1
  #kill -10 $MAIN_PROC
  #exit 0
} &
SLEEP_WRAPPER_ID=$!
PGREP_SLEEP_PATTERN="pgrep -P $SLEEP_WRAPPER_ID sleep"
SLEEP_ID=$($PGREP_SLEEP_PATTERN)

#profile 'done-setup-graceful-timer'

# Invoke this function when we no longer need the graceful timeout timer
function disable_graceful_timeout() {
  #profile 'start-disable-graceful-timeout'
  if [ -z "$SLEEP_ID" ]; then
    SLEEP_ID=$($PGREP_SLEEP_PATTERN)
    #debug "Trying again to figure out sleep id. It is: $SLEEP_ID"
  fi
  #profile 'killing-sleep-wrapper'
  kill -n 15 $SLEEP_WRAPPER_ID 2>/dev/null
  wait $SLEEP_WRAPPER_ID 2>/dev/null
  #profile 'killing-sleep'
  kill -n 15 $SLEEP_ID 2>/dev/null
  wait $SLEEP_ID 2>/dev/null
  #profile 'done-disable-graceful-timeout'
}

# ------------------------------------------------------------------------------
# Grab global allow list
# ------------------------------------------------------------------------------

    GLOBAL_ALLOW_LIST_PATH="$(dirname "${BASH_SOURCE[0]}")/${ALLOW_LIST_PATH}"
    readonly GLOBAL_ALLOW_LIST_PATH

    GLOBAL_ALLOWED_STRINGS=$((jq .allowStrings $GLOBAL_ALLOW_LIST_PATH 2>/dev/null || echo "{}") | sed 's/null/{}/g')

    GLOBAL_ALLOWED_PATHS=$((jq .allowPaths $GLOBAL_ALLOW_LIST_PATH 2>/dev/null || echo "[]") | sed 's/null/[]/g')

    #if [[ -f "$GLOBAL_ALLOW_LIST_PATH" ]]; then
    #    echo -e "\n\nREADME exists.\n\n"
    #else
    #allowed path
    #allowed secrets
    ##jq ".allowed_ips += $(echo "[\"7819\",\"123123\"]") | .ppp += $(echo "[\"7777\",\"333\"]")" <(echo '{"allowed_ips": ["123","456"], "ppp": ["4"]}') 



# ------------------------------------------------------------------------------
# Pre-receive hook
# ------------------------------------------------------------------------------

#profile 'start-secret-scanning'

branches_that_modify_allow_list=()
branches_with_invalid_allow_list=()

declare -A secrets_per_branch
declare -A commits_that_modify_allowlist_per_branch
declare -A commit_diffs_that_modify_allowlist_per_branch
for (( i=0; i<${#refnames[@]}; i++ ))
do
  oldrev=${oldrevs[$i]}
  newrev=${newrevs[$i]}
  refname=${refnames[$i]}

  # debug "Processing branch ${refname}: ${oldrev} -> ${newrev}"


  # ----------------------------------------------------------------------------
  # Get the list of all the commits to scan
  # ----------------------------------------------------------------------------
  span=`git rev-list $(git for-each-ref --format='%(refname)' refs/heads/* | sed 's/^/\^/') ${newrev}`

  # debug "Commits pushed for branch $refname: $span"


  # ----------------------------------------------------------------------------
  # Iterate over all commits in the branch
  # ----------------------------------------------------------------------------

  #profile 'start-secret-scanning-for-single-branch'

  if [ -z "$span" ]
  then
      wait_for_async_subprocesses
      echo "Server-side secret scanner (Langlock) detected no potential secrets in the push."
      exit 0
  fi

  function git_log_for_branch() {
    git show $span -m --no-prefix --ignore-space-change --unified=0 -- . ":(exclude)$ALLOW_LIST_PATH" | grep -v '^@@' | egrep -v '^-[^-]'
    #profile 'done-git_log_for_branch'
  }
  function load_allowlist_for_branch() {
    # git show $newrev:$ALLOW_LIST_PATH 2>/dev/null | jq '.["allowStrings"]' 2>&1
    git show $newrev:$ALLOW_LIST_PATH 2>/dev/null || echo '{"allowStrings": {}, "allowPaths": []}'
  }
  function load_combined_global_and_local_allowlist() {
    jq ".allowStrings += ${GLOBAL_ALLOWED_STRINGS} | .allowPaths += ${GLOBAL_ALLOWED_PATHS}" <(load_allowlist_for_branch) 2>/dev/null || echo "MISFORMATTED ALLOW LIST"
    #profile 'done-load_combined_global_and_local_allowlist'
  }
  function list_commits_that_modify_allow_list() {
    git show --no-patch --format=%H $span -- $ALLOW_LIST_PATH | tr '\n' ' ' | sed 's/ /, /g'  |  sed 's/, $//g'
  }
  function show_commit_diffs_that_modify_allow_list() {
    git show $span -- $ALLOW_LIST_PATH
  }
  function detectcreds() {
    #profile 'start-detectcreds'
    /home/githook/detectcreds.bin <(git_log_for_branch) --in-type "log" --out-type "pretty" --threads $NUM_THREADS --allow-list <(load_combined_global_and_local_allowlist) --allow-list-name "$ALLOW_LIST_PATH" 2>&1
    #profile 'done-detectcreds'
  }
  #debug "Allow-list: $(load_allowlist_for_branch)"
  #debug "Git log: $(git_log_for_branch)"
  secrets_in_branch="$(detectcreds)"
  #debug "Secrets detected: ${secrets_in_branch}"


  err=$?
  if [ $err -eq 4 ]
  then
    branches_with_invalid_allow_list+=( "$refname" )
    continue
  elif [ $err -ne 0 ]; then
    log "*FAILURE in Langlock pre-receive hook* for user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\`. *PUSH ACCEPTED.*\nExit code: $err\nError:\n$secrets_in_branch" "$SLACK_USERS_TO_ALERT"
    disable_graceful_timeout
    echo -e "Server-side secret scanner (Langlock) experienced fatal runtime error, but git push still accepted.\nA debugging report has already been sent to the Github administrators."
    #wait $INITIAL_LOG_ID
    wait_for_async_subprocesses
    exit 0
  fi


  if [ ! -z "$secrets_in_branch" ]; then
    secrets_per_branch["$refname"]="$secrets_in_branch"
  fi

  commits_that_modify_allowlist="$(list_commits_that_modify_allow_list)"
  if [ ! -z "$commits_that_modify_allowlist" ]; then
    commits_that_modify_allowlist_per_branch["$refname"]="$commits_that_modify_allowlist"
  fi

  commit_diffs_that_modify_allowlist="$(show_commit_diffs_that_modify_allow_list)"
  if [ ! -z "$commit_diffs_that_modify_allowlist" ]; then
    commit_diffs_that_modify_allowlist_per_branch["$refname"]="$commit_diffs_that_modify_allowlist"
  fi

  #profile 'done-secret-scanning-for-single-branch'

done

# ------------------------------------------------------------------------------
# Simulate timeout
# ------------------------------------------------------------------------------

# sleep 20


# ------------------------------------------------------------------------------
# Turn off graceful timeout timer, since our script is almost done
# ------------------------------------------------------------------------------

disable_graceful_timeout


# ------------------------------------------------------------------------------
# Block push if we detected malformed allow list (not necessarily updated)
# ------------------------------------------------------------------------------

if [ ${#branches_with_invalid_allow_list[@]} -gt 0 ]; then
  formatted_list=$(join '\n   • ' branches_with_invalid_allow_list)
  formatted_str="ERROR: Your push was blocked because we detected invalid syntax\nin the allow list files for branches:\n\n   • $formatted_list"
  if [ $EFFECTIVE_COLOR = 'true' ]; then
    formatted_str=$(echo -e "$formatted_str" | awk -v bold=$BOLD -v red=$RED -v plain=$PLAIN '{print bold red $0 plain}')
  fi
  echo -e "$HEADER\n\n$formatted_str\n\n$ALLOW_LIST_FORMAT_MESSAGE\n\n$FOOTER"
  #profile 'exiting-block-push-due-to-invalid-allow-list'
  log "*Push blocked due to invalid allow list*" "$SLACK_USERS_TO_ALERT"
  #wait $INITIAL_LOG_ID
  wait_for_async_subprocesses
  exit 1
fi


# ------------------------------------------------------------------------------
# Block push (return non-zero exit code) if secrets found
# ------------------------------------------------------------------------------

if [ ${#secrets_per_branch[@]} -gt 0 ]; then
  echo -e "$HEADER\n\nYour push was blocked because it contained the following errors:\n"
  for branch_name in "${!secrets_per_branch[@]}"
  do
    echo -e "\nERROR: detected secrets in branch ${branch_name}:\n\n${secrets_per_branch[$branch_name]}"
  done
  echo -e "\n\n$FOOTER"
  #profile 'exiting-block-push-due-to-detected-secrets'
  log "*Push blocked due to detected secrets*" "$SLACK_USERS_TO_ALERT"
  #wait $INITIAL_LOG_ID
  wait_for_async_subprocesses
  exit 1
fi


# ------------------------------------------------------------------------------
# Notify slack if there was a successful push/modification to the allow list
# ------------------------------------------------------------------------------

if [ ${#commits_that_modify_allowlist_per_branch[@]} -gt 0 ]; then
  #profile 'start-notify-slack-about-update-to-allow-list'
  formatted_str=""
  formatted_diff_str=""
  for branch_name in "${!commits_that_modify_allowlist_per_branch[@]}"
  do
    formatted_str="$formatted_str• Branch ${branch_name}: ${commits_that_modify_allowlist_per_branch["$branch_name"]}\n"
    formatted_diff_str="$formatted_diff_str\n\n${commit_diffs_that_modify_allowlist_per_branch["$branch_name"]}\n"
  done
  log "*Allow list updated* by user \`$GITHUB_USER_LOGIN\` in repo \`$REPO_NAME\`:\n$formatted_str$formatted_diff_str\n\n" "$SLACK_USERS_TO_ALERT"
  #profile 'done-notify-slack-about-update-to-allow-list'
fi


# ------------------------------------------------------------------------------
# Exit with success if no secrets found
# ------------------------------------------------------------------------------

#profile 'exiting-success'

echo "Server-side secret scanner (Langlock) detected no potential secrets in the push."
log "*Push succeeded, no new secrets detected*"
#wait $INITIAL_LOG_ID
wait_for_async_subprocesses
exit 0
