package go_test2

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

const apiVersion = "1"
const dataURL = "http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z"

func Launch(connStr, servePort string, endpoints []string) error {
	if len(connStr) == 0 {
		return errors.New("connections string is required")
	}
	if len(servePort) == 0 {
		return errors.New("port flag is required")
	}
	if len(endpoints) == 0 {
		return errors.New("at least one endpoint is required")
	}
	context, err := NewContext(connStr)
	if err != nil {
		glog.Error(err)
		return err
	}
	if err := FillDBIfEmpty(context.Session, dataURL); err != nil {
		glog.Error(err)
		return err
	}
	epMap := make(map[string]bool)
	for _, ep := range endpoints {
		epMap[ep] = true
	}
	http.HandleFunc("/api/v"+apiVersion, func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "not supported endpoint", http.StatusNotFound)
			return
		}
		ep := parts[2]
		if _, ok := epMap[ep]; !ok {
			http.Error(w, "not supported endpoint", http.StatusNotFound)
			return
		}
		_ = context
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("ok"))

	})
	if err := http.ListenAndServe(":"+servePort, nil); err != nil {
		glog.Error(err)
	}
	return nil
}
