package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
)

type jsonData struct {
	Data string  `json:"Data"`
	Ans  float64 `json:"Ans"`
}

func main() {
	fs := http.StripPrefix("/", http.FileServer(http.Dir("./static")))
	http.Handle("/", fs)
	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/clearhistory", clearHistory)

	fmt.Println("Listening")
	http.ListenAndServe(":8080", nil)
}
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func calculate(w http.ResponseWriter, r *http.Request) {
	var fetchdata jsonData
	err := json.NewDecoder(r.Body).Decode(&fetchdata)
	// fmt.Println(fetchdata)
	if err != nil {
		fmt.Println(err)
	} else {
		var ans []string
		var temp string
		for i := 0; i < len(fetchdata.Data); i++ {
			char := string(fetchdata.Data[i])
			if char == "*" || char == "/" || char == "%" || char == "+" || char == "-" {
				ans = append(ans, temp)
				temp = ""
				ans = append(ans, char)
			} else {
				temp += char
			}
		}
		ans = append(ans, temp)
		// fmt.Println(ans)
		// var ans2 []string
		// j := 0
		for i := 0; i < len(ans); {
			char := ans[i]
			// fmt.Println(char)
			if char == "*" || char == "/" || char == "%" {
				valA, _ := strconv.ParseFloat(ans[i-1], 64)
				valB, _ := strconv.ParseFloat(ans[i+1], 64)
				var res float64
				switch char {
				case "*":
					res = valA * valB
					break
				case "/":
					if valB == 0 {
						res, _ := json.Marshal(struct {
							IsError bool
							Msg     string
						}{true, "In division, zero is not allowed"})
						w.Write(res)
						return
					}
					res = valA / valB
					break
				case "%":
					if valB == 0 {
						res, _ := json.Marshal(struct {
							IsError bool
							Msg     string
						}{true, "In modulo, zero is not allowed"})
						w.Write(res)
						return
					}
					res = math.Mod(valA, valB)
					break
				}
				strres := strconv.FormatFloat(res, 'f', -1, 64)
				ans = RemoveIndex(ans, i)
				ans = RemoveIndex(ans, i)
				ans[i-1] = strres
				// ans2[j-1] = strres
			} else {
				// ans2 = append(ans2, ans[i])
				// j++
				i++
			}
		}
		i := 0
		// fmt.Println(ans)
		for len(ans) != 1 {
			char := ans[i]
			if char == "+" || char == "-" {
				valA, _ := strconv.ParseFloat(ans[i-1], 64)
				valB, _ := strconv.ParseFloat(ans[i+1], 64)
				var res float64
				switch char {
				case "+":
					res = valA + valB
					break
				case "-":
					res = valA - valB
					break
				}
				strres := strconv.FormatFloat(res, 'f', -1, 64)
				ans = RemoveIndex(ans, i)
				ans = RemoveIndex(ans, i)
				ans[i-1] = strres
				// ans2[j-1] = strres
			} else {
				i++
			}
		}
		fetchdata.Ans, _ = strconv.ParseFloat(ans[0], 64)
		writeHistory(fetchdata)
		res, _ := json.Marshal(struct {
			IsError bool
			Msg     string
			Data    jsonData
		}{false, "Modulo Performed Successully", fetchdata})
		w.Write(res)
		// fmt.Println(ans)
	}
}

func clearHistory(w http.ResponseWriter, r *http.Request) {
	if err := os.Remove("history.json"); err != nil {
		fmt.Println(err)
	}
}

func writeHistory(fetchdata jsonData) {
	DataJSONFile, err := os.Open("history.json")
	if err != nil {
		fmt.Println(err)
	}
	defer DataJSONFile.Close()
	DataByteValue, _ := ioutil.ReadAll(DataJSONFile)
	var DataArray []jsonData
	json.Unmarshal(DataByteValue, &DataArray)

	DataArray = append(DataArray, fetchdata)

	data, err := json.MarshalIndent(DataArray, "", "	")

	if err = ioutil.WriteFile("history.json", data, 0644); err != nil {
		fmt.Println(err)
	}
}
