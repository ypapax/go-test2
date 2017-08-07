package go_test2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
)

const apiVersion = "1"

var actions = map[string]bool{
	"min": true,
	"max": true,
	"avg": true,
	"":    true,
}

const reqDateLayout = "02/01/2006"

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
	context, err := newContext(connStr)
	if err != nil {
		glog.Error(err)
		return err
	}
	if err := fillDBIfEmpty(context.Session, dataURL); err != nil {
		glog.Error(err)
		return err
	}
	epMap := make(map[string]bool)
	dbIndices := []mgo.Index{{Key: []string{timeField}}}
	for _, ep := range endpoints {
		epMap[ep] = true
		dbIndices = append(dbIndices, mgo.Index{Key: []string{ep + qcSuffix, timeField, ep}})
		dbIndices = append(dbIndices, mgo.Index{Key: []string{ep + qcSuffix, ep}})
	}
	if err := ensureIndeces(context.Session, dbIndices); err != nil {
		glog.Error(err)
		return err
	}
	http.HandleFunc("/test/api/v"+apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		glog.V(4).Infof("parts %+v, len(parts) = %+v", parts, len(parts))
		if len(parts) != 5 && len(parts) != 6 {
			glog.V(4).Infof("not supported endpoint")
			http.Error(w, "not supported endpoint", http.StatusNotFound)
			return
		}
		parts = parts[2:]
		ep := parts[2]
		if _, ok := epMap[ep]; !ok {
			glog.V(4).Infof("not supported endpoint")
			http.Error(w, "not supported endpoint", http.StatusNotFound)
			return
		}
		var aggr string
		if len(parts) > 3 {
			aggr = parts[3]
		}
		if _, ok := actions[aggr]; !ok {
			glog.V(4).Infof("not supported endpoint")
			http.Error(w, "not supported endpoint", http.StatusNotFound)
			return
		}
		s := context.Session.Copy()
		defer s.Close()
		start, err := dateParse(r, "start")
		if err != nil {
			http.Error(w, "unable to parse date for start parameter", http.StatusBadRequest)
			return
		}
		stop, err := dateParse(r, "stop")
		if err != nil {
			http.Error(w, "unable to parse date for stop parameter", http.StatusBadRequest)
			return
		}
		var rslt interface{}
		switch aggr {
		case "":
			rslt, err = find(s, ep, start, stop)

		default:
			rslt, err = aggregate(s, ep, aggr, start, stop)
		}
		if err != nil && err != mgo.ErrNotFound {
			glog.Errorf("%+v for %+v %+v %+v %+v", ep, aggr, start, stop)
			http.Error(w, "internal service error", http.StatusInternalServerError)
			return
		}
		glog.V(4).Infof("rs %+v", rslt)
		if err == mgo.ErrNotFound {
			bnds, err := getBounds(s, ep)
			if err != nil {
				glog.Error(err)
				writeResp(&errorResp{Error: "Out of bounds. Internal error in getting supported date bounds."}, w)
				return
			}
			writeResp(&errorResp{
				Error: fmt.Sprintf("Out of bounds. Supported bounds: start=%+v stop=%+v",
					bnds.Min.Format(reqDateLayout), bnds.Max.Format(reqDateLayout))}, w)
			return
		}
		glog.V(4).Infof("rs %+v", rslt)
		writeResp(rslt, w)

	})
	glog.Infof("listening " + servePort)
	if err := http.ListenAndServe(":"+servePort, nil); err != nil {
		glog.Error(err)
	}
	return nil
}

func dateParse(r *http.Request, paramName string) (*time.Time, error) {
	param := r.URL.Query()[paramName]
	if len(param) == 0 {
		return nil, nil
	}
	t, err := time.Parse(reqDateLayout, param[0])
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return &t, nil
}

type errorResp struct {
	Error string `json:"error"`
}

func writeResp(obj interface{}, w http.ResponseWriter) error {
	b, err := json.Marshal(obj)
	if err != nil {
		glog.Error(err)
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}
