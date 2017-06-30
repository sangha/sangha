#! /bin/sh

### BEGIN INIT INFO
# Provides:             sangha
# Required-Start:       $syslog
# Required-Stop:        $syslog
# Default-Start:        2 3 4 5
# Default-Stop:
# Short-Description:    sangha daemon
### END INIT INFO

set -e

# /etc/init.d/sangha: start and stop the sangha daemon

umask 022

. /lib/lsb/init-functions

#export GOROOT="/usr/local/opt/go"
export PATH="$GOROOT/bin:$PATH"
export GOPATH="/home/sangha/go"

BINARY="$GOPATH/bin/sangha"
PIDFILE="/var/run/sangha.pid"

test -x $BINARY || exit 0

case "$1" in
  start)
        log_daemon_msg "Starting sangha daemon" "sangha" || true
        if start-stop-daemon --start -b --quiet --oknodo -m --pidfile $PIDFILE --exec $BINARY ; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;
  stop)
        log_daemon_msg "Stopping sangha daemon" "sangha" || true
        if start-stop-daemon --stop --quiet --oknodo --pidfile $PIDFILE; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  reload|force-reload)
        log_daemon_msg "Reloading sangha daemon's configuration" "sangha" || true
        if start-stop-daemon --stop --signal 1 --quiet --oknodo --pidfile $PIDFILE --exec $BINARY ; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  restart)
        log_daemon_msg "Restarting sangha daemon" "sangha" || true
        start-stop-daemon --stop --quiet --oknodo --retry 30 --pidfile $PIDFILE
        if start-stop-daemon --start -b --quiet --oknodo -m --pidfile $PIDFILE --exec $BINARY ; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  status)
        status_of_proc -p $PIDFILE $BINARY sangha && exit 0 || exit $?
        ;;

  *)
        log_action_msg "Usage: /etc/init.d/sangha {start|stop|reload|restart|status}" || true
        exit 1
esac

exit 0
