package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

type data struct {
	ValA   float64 `json:"ValA"`
	ValB   float64 `json:"ValB"`
	Action string  `json:"Action"`
}

type jsonData struct {
	ValA   float64 `json:"ValA"`
	ValB   float64 `json:"ValB"`
	Ans    float64 `json:"Ans"`
	Action string  `json:"Action"`
}

func main() {
	fs := http.StripPrefix("/", http.FileServer(http.Dir("./static")))
	http.Handle("/", fs)
	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/clearhistory", clearHistory)

	fmt.Println("Listening")
	http.ListenAndServe(":8000", nil)
}

func calculate(w http.ResponseWriter, r *http.Request) {
	var fetchdata data
	err := json.NewDecoder(r.Body).Decode(&fetchdata)
	// fmt.Println(fetchdata)
	if err != nil {
		fmt.Println(err)
	} else {
		switch fetchdata.Action {
		case "*":
			ans := fetchdata.ValA * fetchdata.ValB
			writeHistory(fetchdata, ans)
			res, _ := json.Marshal(struct {
				IsError bool
				Msg     string
				Ans     float64
				Data    data
			}{false, "Multiplication Performed Successully", ans, fetchdata})
			w.Write(res)
		case "+":
			ans := fetchdata.ValA + fetchdata.ValB
			writeHistory(fetchdata, ans)
			res, _ := json.Marshal(struct {
				IsError bool
				Msg     string
				Ans     float64
				Data    data
			}{false, "Addition Performed Successully", ans, fetchdata})
			w.Write(res)
			break
		case "-":
			ans := fetchdata.ValA - fetchdata.ValB
			writeHistory(fetchdata, ans)
			res, _ := json.Marshal(struct {
				IsError bool
				Msg     string
				Ans     float64
				Data    data
			}{false, "Subtraction Performed Successully", ans, fetchdata})
			w.Write(res)
			break
		case "/":
			if fetchdata.ValB == 0 {
				res, _ := json.Marshal(struct {
					IsError bool
					Msg     string
				}{true, "In denominator, zero is not allowed"})
				w.Write(res)
				break
			}
			ans := fetchdata.ValA / fetchdata.ValB
			writeHistory(fetchdata, ans)
			res, _ := json.Marshal(struct {
				IsError bool
				Msg     string
				Ans     float64
				Data    data
			}{false, "Divison Performed Successully", ans, fetchdata})
			w.Write(res)
			break
		case "%":
			if fetchdata.ValB == 0 {
				res, _ := json.Marshal(struct {
					IsError bool
					Msg     string
				}{true, "In modulo, zero is not allowed"})
				w.Write(res)
				break
			}

			ans := math.Mod(fetchdata.ValA, fetchdata.ValB)
			writeHistory(fetchdata, ans)
			res, _ := json.Marshal(struct {
				IsError bool
				Msg     string
				Ans     float64
				Data    data
			}{false, "Modulo Performed Successully", ans, fetchdata})
			w.Write(res)
			break
		}
	}
}

func clearHistory(w http.ResponseWriter, r *http.Request) {
	if err := os.Remove("history.json"); err != nil {
		fmt.Println(err)
	}
}

func writeHistory(fetchdata data, ans float64) {
	var jsonFetchData jsonData
	jsonFetchData = convert(fetchdata, ans)
	DataJSONFile, err := os.Open("history.json")
	if err != nil {
		fmt.Println(err)
	}
	defer DataJSONFile.Close()
	DataByteValue, _ := ioutil.ReadAll(DataJSONFile)
	var DataArray []jsonData
	json.Unmarshal(DataByteValue, &DataArray)

	DataArray = append(DataArray, jsonFetchData)

	data, err := json.MarshalIndent(DataArray, "", "	")

	if err = ioutil.WriteFile("history.json", data, 0644); err != nil {
		fmt.Println(err)
	}
}

func convert(fetchdata data, ans float64) jsonData {
	var temp jsonData
	temp.ValA = fetchdata.ValA
	temp.ValB = fetchdata.ValB
	temp.Action = fetchdata.Action
	temp.Ans = ans
	return temp
}
