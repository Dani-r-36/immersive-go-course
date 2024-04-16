package main

import (
	"reflect"
	"testing"
	"bytes"
)

func TestWriterReturn(t *testing.T){
	jsonData := `[
  {"name": "Aya", "high_score": 2},
  {"name": "Prisha", "high_score": 30},
  {"name": "Charlie", "high_score": -1},
  {"name": "Margot", "high_score": 50}
]`

	t.Run("json data", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString(jsonData)
		record, _ := jsonText(&buffer)
		expectedRecords := []Record{
		{Name: "Aya", Score: 2},
		{Name: "Prisha", Score: 30},
		{Name: "Charlie", Score: -1},
		{Name: "Margot", Score: 50},
		}
		reflect.DeepEqual(record,expectedRecords)
	})
	repeatJsonData := `# This file contains lines of data stored in JSON format. Each line contains exactly one record stored as an object.
# Lines starting with # are comments and should be ignored.
{"name": "Aya", "high_score": 50}
{"name": "Prisha", "high_score": 60}
# Charlie didn't do fantastically :(
{"name": "Charlie", "high_score": -1}
{"name": "Margot", "high_score": 50}
`

	t.Run("repeat json data", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString(repeatJsonData)
		record, _ := repeatedJson(&buffer)
		expectedRecords := []Record{
		{Name: "Aya", Score: 50},
		{Name: "Prisha", Score: 60},
		{Name: "Charlie", Score: -1},
		{Name: "Margot", Score: 50},
		}
		reflect.DeepEqual(record,expectedRecords)
	})
	csvData := `name,high score
				Aya,100
				Prisha,3
				Charlie,-1
				Margot,25`
	t.Run("csv data", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString(csvData)
		record, _ := csvFile(&buffer)
		expectedRecords := []Record{
		{Name: "Aya", Score: 100},
		{Name: "Prisha", Score: 3},
		{Name: "Charlie", Score: -1},
		{Name: "Margot", Score: 25},
		}
		reflect.DeepEqual(record,expectedRecords)
	})

	
	

}