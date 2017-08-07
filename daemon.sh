set -ex
PIDFILE=$PROGDIR/$PROGNAME.pid
cd $PROGDIR;
pwd;
LOGDIR=$PROGDIR
OUTLOG=$PROGNAME.log
start-stop-daemon  -c www-data --make-pidfile --pidfile "$PIDFILE" --background \
    --no-close --exec  "$PROGDIR/${PROGNAME}" --start -- $ARGS \
    >> "${LOGDIR}/${OUTLOG}" 2>> "${LOGDIR}/${OUTLOG}" </dev/null