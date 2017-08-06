package go_test2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const dbName = "go-test2"
const collectionName = "neracoos"
const dataURL = "http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z"
const qcSuffix = "_qc"
const timeField = "time"

var timeFields = map[string]bool{
	"time":          true,
	"time_created":  true,
	"time_modified": true,
}

const timeLayout = time.RFC3339

//ds.Session.Copy()

func NewContext(dialStr string) (*Context, error) {
	session, err := mgo.Dial(dialStr)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return &Context{
		Session: session,
	}, nil
}
func find(s *mgo.Session, field string, start, stop *time.Time) ([]map[string]interface{}, error) {
	glog.V(4).Infof("find %+v %+v %+v", field, start, stop)
	var rslt []map[string]interface{}
	pipe := matchPipe(field, start, stop)
	project := map[string]interface{}{
		"_id": 0,
		"time": 1,
	}
	project[field] = 1
	pipe = append(pipe, map[string]interface{}{"$project": project})
	pipe = append(pipe, map[string]interface{}{"$sort": map[string]interface{}{"time": -1}})
	if err := collection(s).Pipe(pipe).All(&rslt); err != nil {
		return nil, err
	}
	if len(rslt) == 0 {
		return nil, mgo.ErrNotFound
	}
	return rslt, nil
}
func matchPipe(field string, start, stop *time.Time) []bson.M {
	q := make(map[string]interface{})
	q[field+qcSuffix] = 0
	var tm map[string]interface{}
	if start != nil || stop != nil {
		tm = make(map[string]interface{})
		q["time"] = tm
	}
	if start != nil {
		tm["$gte"] = start
	}
	if stop != nil {
		tm["$lte"] = stop
	}
	pipe := []bson.M{{"$match": q}}
	return pipe

}
func aggregate(s *mgo.Session, field, action string, start, stop *time.Time) (*result, error) {
	var gr = make(map[string]interface{})
	gr["$"+action] = "$" + field
	pipe := matchPipe(field, start, stop)
	pipe = append(pipe, bson.M{"$group": bson.M{"_id": nil, "result": gr}})
	glog.V(4).Infof("pipe %+v", pipe)
	var rs result
	if err := collection(s).Pipe(pipe).One(&rs); err != nil {
		return nil, err
	}
	glog.V(4).Infof("result %+v", rs)
	return &rs, nil
}

type result struct {
	Result float64 `json:"result",bson:"result"`
}

type Context struct {
	Session *mgo.Session
}

func collection(s *mgo.Session) *mgo.Collection {
	return s.DB(dbName).C(collectionName)
}

func fillDBIfEmpty(s *mgo.Session, url string) error {
	n, err := collection(s).Count()
	if err != nil {
		glog.Error(err)
		return err
	}
	if n > 0 {
		glog.Infof("db is not empty: %+v, skip filling", n)
		return nil
	}
	return fillDB(s, url)
}

func fillDB(s *mgo.Session, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("invalid status code %+v for requesting url %+v", resp.Status, url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respD respData
	if err := json.Unmarshal(b, &respD); err != nil {
		return err
	}
	columns := respD.Table.ColumnNames
	var itemsToInsert []interface{}
	for _, r := range respD.Table.Rows {
		if len(columns) != len(r) {
			err := fmt.Errorf("coluns count doesnt match row item count: columns: %+v, row: %+v", columns, r)
			return err
		}
		d := make(map[string]interface{})
		for i, v := range r {
			c := columns[i]
			if _, isTimeField := timeFields[c]; !isTimeField {
				d[c] = v
				continue
			}
			t, err := time.Parse(timeLayout, fmt.Sprintf("%+v", v))
			if err != nil {
				glog.Error(err)
				return err
			}
			d[c] = t
		}
		itemsToInsert = append(itemsToInsert, d)
	}
	if err := collection(s).Insert(itemsToInsert...); err != nil {
		glog.Error(err)
		return err
	}
	glog.Infof("inserted %+v items to db", len(itemsToInsert))
	return nil
}

type respData struct {
	Table table `json:"table"`
}

type table struct {
	ColumnNames []string        `json:"columnNames"`
	Rows        [][]interface{} `json:"rows"`
}

func getBounds(s *mgo.Session) (*bounds, error) {
	var bnds bounds
	if err := collection(s).Pipe([]interface{}{map[string]map[string]map[string]string{
		"$group": map[string]map[string]string{
			"_id": nil,
			"min": map[string]string{
				"$min": "$time",
			},
			"max": map[string]string{
				"$max": "$time",
			},
		},
	}}).One(&bnds); err != nil {
		glog.Error(err)
		return nil, err
	}
	return &bnds, nil
}

type bounds struct {
	Min time.Time `bson:"min"`
	Max time.Time `bson:"max"`
}
