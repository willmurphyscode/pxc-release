function pid_is_running() {
  declare pid="$1"
  ps -p "${pid}" >/dev/null 2>&1
}

# wait_pid_death
#
# @param pid
# @param timeout
#
# Watch a :pid: for :timeout: seconds, waiting for it to die.
# If it dies before :timeout:, exit 0. If not, exit 1.
#
# Note that this should be run in a subshell, so that the current
# shell does not exit.
#
function wait_pid_death() {
  declare pid="$1" timeout="$2"

  local countdown
  countdown=$(( timeout * 10 ))

  while true; do
    if ! pid_is_running "${pid}"; then
      return 0
    fi

    if [ ${countdown} -le 0 ]; then
      return 1
    fi

    countdown=$(( countdown - 1 ))
    sleep 0.1
  done
}

# kill_and_wait
#
# @param pidfile
# @param timeout [default 25s]
#
# For a pid found in :pidfile:, send a `kill`, then wait for :timeout: seconds to
# see if it dies on its own. If not, send it a `kill -9`. If the process does die,
# exit 0 and remove the :pidfile:. If after all of this, the process does not actually
# die, exit 1.
#
# Note:
# Monit default timeout for start/stop is 30s
# Append 'with timeout {n} seconds' to monit start/stop program configs
#
function kill_and_wait() {
  declare pidfile="$1" timeout="${2:-25}" sigkill_on_timeout="${3:-1}"

  if [ ! -f "${pidfile}" ]; then
    echo "Pidfile ${pidfile} doesn't exist"
    return 0
  fi

  local pid
  pid=$(head -1 "${pidfile}")

  if [ -z "${pid}" ]; then
    echo "Unable to get pid from ${pidfile}"
    return 1
  fi

  if ! pid_is_running "${pid}"; then
    echo "Process ${pid} is not running"
    rm -f "${pidfile}"
    return 0
  fi

  echo "Killing ${pidfile}: ${pid} "
  kill "${pid}"

  if ! wait_pid_death "${pid}" "${timeout}"; then
    if [ "${sigkill_on_timeout}" = "1" ]; then
      echo "Kill timed out, using kill -9 on ${pid}"
      kill -9 "${pid}"
      sleep 0.5
    fi
  fi

  if pid_is_running "${pid}"; then
    echo "Timed Out"
    return 1
  else
    echo "Stopped"
    rm -f "${pidfile}"
    return 0
  fi
}

log(){
  message=$1
  echo "$(date +"%Y-%m-%d %H:%M:%S %z") ----- $message"
}
