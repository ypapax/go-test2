#!/usr/bin/env bash
set -ex
if [ -z $ZIP_NAME ]; then
    >&2 echo missing ZIP_NAME
    exit 1
fi
monitFile=/tmp/go-test2monit
monitDir=/etc/monit/conf.d
rm -rf $monitFile
dir=$HOME/go/go-test2
pushd $dir
unzip -o $ZIP_NAME
popd
bin=go-test2

daemonFile=$dir/daemon.sh
ARGS=$@
echo -e "#!/usr/bin/env bash\nPROGDIR=$dir; PROGNAME=$bin; ARGS=\"$ARGS\"; \n$(cat $daemonFile)" > $daemonFile

echo check process go-test2 with pidfile $dir/$bin.pid >> $monitFile
echo "   "start program = \"$daemonFile\" >> $monitFile
echo "   "stop program = \"/usr/bin/pkill -ef $bin\" >> $monitFile
cat $monitFile
mv $monitFile $monitDir
monit reload


