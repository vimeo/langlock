

function current_millisecond_timestamp {
  echo $(python3 -c 'import time; print(int(round(time.time() * 1000)))')
}

PROFILE_START_TIME=$(current_millisecond_timestamp)
readonly PROFILE_START_TIME

function elapsed_milliseconds {
   local current_time=$(current_millisecond_timestamp)
   echo "$1 $(($current_time - $PROFILE_START_TIME))"
}


MAIN_PROC=$$
#trap "elapsed_milliseconds 'TRAP'; kill -s SIGINT -1" 10

{
    sleep 20
    exit 1
} &
{
    sleep 2
    echo "SLEEP 2 DONE"
    exit 1
} &
{
    sleep 2.1
    echo "SLEEP 2.1 DONE"
    exit 1
} &
{
    sleep 2.5
    echo "SLEEP 2.5 DONE"
    exit 1
} &
{
    sleep 2.7
    echo "SLEEP 2.7 DONE"
    exit 1
} &

trap "elapsed_milliseconds 'TRAP'; exit 0" SIGINT SIGQUIT SIGTERM
{
    sleep 1
    #kill -s SIGINT -1

    sleep 1 &
    BACKGROUND_SLEEP_PROC_ID=$!
    wait $!
    elapsed_milliseconds 'After background sleep'

    kill -s SIGINT -1
    #kill $$
    #kill -10 $MAIN_PROC
    sleep 10
    #exit 0
} &

sleep 3
elapsed_milliseconds 'FAILED';
exit 1