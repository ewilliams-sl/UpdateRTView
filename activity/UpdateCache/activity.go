
package UpdateCache

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
	"time"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"os"
)

// Constants
const (
	command = "Target"
	tblDef = "Definition"
	tblData = "Data"
	spStartTime = "StartTime"
	params  = "params"
	result  = "result"
)

// log is the default package logger which we'll use to log
var log = logger.GetLogger("activity-setQoS")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

type slPost struct {
	Metadata []slColDef `json:"metadata"`
	Data []map[string]interface{} `json:"data"`
	cols int
}

type slMetadata struct {
	Metadata []slColDef `json:"metadata"`
}

type slColDef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (d slColDef) string() string {
	return "\"name\":\"" + d.Name + "\",\"type:\""+d.Type
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func  (m *slPost) addCol (n string, t string) {
	var colDef slColDef
	colDef.Name = n
	colDef.Type = t
	iLast := m.cols;
	if (cap(m.Metadata) == 0) {
		m.Metadata = make([]slColDef, 5, 100)
	}
	m.Metadata[iLast] = colDef
	m.cols += 1
}

// Post info contained in row data to rtview via the url
func (rowData *slPost) postRowData(url string) error {
  json, _ := json.Marshal(rowData)
  
  fmt.Println("Sending: ", url)
  fmt.Println("Post:\n", string(json))
  
  resp, err := http.Post(url, "plain/text", bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("response:\n", string(body))
	return nil
}



func (a *MyActivity) updatePerformance(sURL string,  sFlowName string, iMillis int64) {
	var post slPost
//	post.addCol("time_stamp", "date")
	post.addCol("hostName", "string")
	post.addCol("flowName", "string")
	post.addCol("duration", "int")
	
//	t2 := makeTimestamp()
	
	sHostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	
	var d map[string]interface{}
	d = make(map[string]interface{})
//	d["time_stamp"] = t2
	d["hostName"] = sHostName
	d["flowName"] = sFlowName
	d["duration"] = iMillis

	var d2 []map[string]interface{}
	d2 = make([]map[string]interface{}, 1,1)
	d2[0] = d
	
	post.Data = d2
	bolD, _ := json.Marshal(post)

	post.postRowData(sURL)
	fmt.Println(string(bolD))	
}




// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	sTarget := context.GetInput(command).(string)
	fmt.Println(sTarget)
	
	iVal := context.GetInput(spStartTime).(int64)
	sName := context.ActivityHost().Name()
	iEndVal := makeTimestamp() -iVal
	
    a.updatePerformance(sTarget, sName, iEndVal)

	// Signal to the Flogo engine that the activity is completed
	return true, nil
}

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		m[word] += 1
	}
	return m
}

func split(tosplit string, sep rune) []string {
	var fields []string

	last := 0
	for i, c := range tosplit {
		if c == sep {
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1
		}
	}

	// Don't forget the last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}
