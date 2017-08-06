#!/usr/bin/env bash
set -e
addr=${1-http://localhost:8181}
for ep in temperature fake_ep current_speed ; do
    echo
    echo
    for act in  "/max" "/avg" "/fake" "/min" ""; do
        echo
        for dates in "?start=01/01/2016" "" "?start=01/01/2016&stop=01/01/2017" "?stop=01/01/2006"; do
                url="$addr/test/api/v1/${ep}${act}${dates}"
                cmd="curl $url"
                echo $cmd
                $cmd
                echo
#                curl -s -o /dev/null -w "%{http_code}" $url
#                echo
        done
    done
done
