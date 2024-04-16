package main 

import (
	"io"
	"os"
	"fmt"
	"encoding/json"
	"encoding/csv"
	"io/ioutil" //need as json in array
	"strings"
	"flag"
)

type Record struct {
   Name string `json:"name"`
   Score int `json:"high_score"`
}

func main(){
	format := flag.String("format", "", "Format the file is serialised in. Accepted values: json,repeated-json,csv,binary")
	file := flag.String("file", "", "Path to the file to read data from")
	flag.Parse()

	if *file == ""{
		fmt.Fprintf(os.Stderr,"No file passed\n")
	}

	fileContent, err := os.Open(*file)
   	if err != nil {
		fmt.Fprintf(os.Stderr,"Error opening: %v", err)
      	return
   	}
	defer fileContent.Close()
	var records []Record
	switch *format{
		case "json":
			records, err = jsonText(fileContent)
		case "repeated-json":
			records, err = repeatedJson(fileContent)
		case "csv":
			records, err = csvFile(fileContent)
		// case "bin":
		// 	parser = &binary.Parser{}
		case "":
			fmt.Fprintf(os.Stderr, "format is a required argument")
		default:
			fmt.Fprintf(os.Stderr,"Didn't know how to parse format %q", *format)
	}
	if err !=nil{return}
	highestScore := records[0].Score
	highestPerson := records[0].Name

	for _, r := range records {
		if r.Score > highestScore{
			highestScore = r.Score
			highestPerson = r.Name
		}
	}
	fmt.Printf("%s has the highest score of %d\n\n", highestPerson, highestScore)
}


func jsonText(fileContent io.Reader) ([]Record, error) {
   //read file using ioutil as in array 
   data, err := ioutil.ReadAll(fileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error reading file: %v", err)
		return nil, err
	}
   var records []Record

   //unmarshal converts the json into govalues, here thats a slice of the struct
   err = json.Unmarshal(data, &records)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error decoding JSON: %v", err)
		return nil, err
	}
	return records, nil
   
}

func repeatedJson(fileContent io.Reader) ([]Record, error){
   data, err := ioutil.ReadAll(fileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error reading file: %v", err)
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var recordsrepeated []Record
	// Iterate over each line in the file
	for _, line := range lines {
		// Ignore lines starting with "#"
		if strings.HasPrefix(line, "#") {
			continue
		}
   		var record Record
		//unmarshal converts the json into govalues, here thats a slice of the struct
		err := json.Unmarshal([]byte(line), &record)
			if err != nil {
				fmt.Fprintf(os.Stderr,"Error decoding JSON: %v\nProbably EOF\n", err)
				continue
			}
		recordsrepeated = append(recordsrepeated, record)
		}
   return recordsrepeated, nil
}

func csvFile(fileContent io.Reader) ([]Record, error){
    csvReader := csv.NewReader(fileContent)
    header, err := csvReader.Read()
    if err != nil {
        fmt.Fprintf(os.Stderr,"Unable to parse file as CSV for %v" ,err)
		return nil, err
    }
	if len(header) !=2 || header[0] != "name" || header[1] != "high score" {
		fmt.Fprintf(os.Stderr,"invalid CSV header\n")
		return nil, err
	}
	var recordscsv []Record
	for {
		record, err := csvReader.Read()
		if err != nil{
			break
		}
		if len(record) != 2 {
				continue
			}
		score := 0
		fmt.Sscanf(record[1], "%d", &score)
		recordscsv = append(recordscsv, Record{Name: record[0], Score: score})
	}
	return recordscsv, nil
}