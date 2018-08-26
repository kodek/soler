package main

import (
	"fmt"

	"io/ioutil"

	"github.com/kodek/soler/greenbutton"
)

func ReadFile(path string) ([]greenbutton.GBPoint, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return greenbutton.Read(string(b))
}

func main() {
	data, err := ReadFile("~/Downloads/sce_interval_data.csv")
	if err != nil {
		panic(err)
	}
	for _, dp := range data {
		fmt.Printf("Time: %s usage: %0.2f\n", dp.T.Format("2006-01-02 15:04"), dp.UsageKwh)
	}
}
