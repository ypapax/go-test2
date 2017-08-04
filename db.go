package go_test2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

const dbName = "go-test2"
const collectionName = "neracoos"

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

type Context struct {
	Session *mgo.Session
}

func collection(s *mgo.Session) *mgo.Collection {
	return s.DB(dbName).C(collectionName)
}

func FillDBIfEmpty(s *mgo.Session, url string) error {
	n, err := collection(s).Count()
	if err != nil {
		glog.Error(err)
		return err
	}
	if n > 0 {
		glog.Infof("db is not empty, skip filling")
		return nil
	}
	return FillDB(s, url)
}

func FillDB(s *mgo.Session, url string) error {
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
			d[columns[i]] = v
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
