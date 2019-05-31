
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
)

// Constants
const (
	command = "Target"
	tblDef = "Definition"
	tblData = "Data"
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


func (a *MyActivity) callPost(sURL string) {
	var post slPost
	post.addCol("time_stamp", "date")
	post.addCol("connection", "string")
	post.addCol("objectExtIdTableSize", "int")
	post.addCol("objectTableSize", "int")
	post.addCol("expired", "boolean")
	
	t2 := makeTimestamp()
	

	var d map[string]interface{}
	d = make(map[string]interface{})
	d["time_stamp"] = t2
	d["connection"] = "new51Cache"
	d["objectExtIdTableSize"] = 10000
	d["objectTableSize"] = 1000
	d["expired"] = false
	var d3 map[string]interface{}
	d3 = make(map[string]interface{})
	d3["time_stamp"] = t2
	d3["connection"] = "newbe4cache"
	d3["objectExtIdTableSize"] = 20000
	d3["objectTableSize"] = 2000
	d3["expired"] = false

	var d2 []map[string]interface{}
	d2 = make([]map[string]interface{}, 2,6)
	d2[0] = d
	d2[1] = d3
	
	post.Data = d2
	bolD, _ := json.Marshal(post)
//	test("http://localhost:5000/rtview/json/data/TbeObjectTable")
//	test("http://localhost:8085/rtview/json/data/TbeObjectTable")

//	post.postRowData("http://10.0.0.20:5000/rtview/write/cache/TbeObjectTable")
	post.postRowData(sURL)
//	post.postRowData("http://localhost:8085/rtview/write/cache/TbeSlop")
	fmt.Println(string(bolD))	
}

func (a *MyActivity) updatePerformance(sURL string, sFlowName string, iMillis int) {
	var post slPost
	post.addCol("time_stamp", "date")
	post.addCol("flowName", "string")
	post.addCol("duration", "int")
	post.addCol("expired", "boolean")
	
	t2 := makeTimestamp()
	

	var d map[string]interface{}
	d = make(map[string]interface{})
	d["time_stamp"] = t2
	d["flowName"] = sFlowName
	d["duration"] = iMillis
	d["expired"] = false

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

	//ivMsg := `{"cmd":"./test.sh","params":"aaa bbb ccc}`
	// put input varable into a slice (note not order guarenteed)
	sTarget := context.GetInput(command).(string)
	fmt.Println(sTarget)
	sTableDef := context.GetInput(tblDef).(string)
	fmt.Println(sTableDef)
	sData := context.GetInput(tblData).(string)
	fmt.Println(sData)
	sName := context.ActivityHost().Name()
		
    a.updatePerformance(sTarget, sName, 100)


	// get the command to execute including path
//	cmd, ok := ivCmdParams[command].(string) // this should be the command or script to execute
//	fmt.Println(cmd);
//	if ok == false {
//		// no Command to execute
//		log.Infof("No Command to execute - check input syntax: [%s]", err)
//		context.SetOutput(result, err.Error())
//		return true, err
//	}
	// We should have the command to execture - check its there
//	_, err = os.Stat(cmd)

//	if err != nil {
	// If the file doesn't exist return error
//		context.SetOutput("result", err.Error())
//		log.Infof("File [%s] does not exist", cmd)
//		return true, err
//	}

	// Get the Commands Params
//	var paramsArray [20]string                    // FIX THIS: make dynamic but ordered
//	cmdParams, ok := ivCmdParams[params].(string) // this is a string containg space separated parameters
//	if ok == false {
		// no params
//		log.Infof("No params provided")
//	} else {
		// Put into array for exec.Command to use, space separated
		// put command arguments into an array in the order they are entered.  Order is important.
//		i := 0
//		for _, field := range split(cmdParams, ' ') {
//			paramsArray[i] = field
//			i++
//		}
//	}

	// launch the command
//	var cmdOut []byte
//	if cmdOut, err = exec.Command(cmd, paramsArray[0:]...).Output(); err != nil {
//		log.Infof("Error running Flogo setQoS activity: [%s]", err)
//		context.SetOutput(result, err.Error())
//		return true, err
//	}
//	rslt := string(cmdOut)
	// Set the result as part of the context
//	context.SetOutput(result, rslt)

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
