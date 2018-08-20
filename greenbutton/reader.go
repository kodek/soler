package greenbutton

import (
	"time"

	"errors"
	"fmt"
	"strconv"

	"strings"

	"github.com/gocarina/gocsv"
)

type GBPoint struct {
	T        time.Time
	UsageKwh float64
}

type Reader struct {
}

const HEADER = "TYPE,DATE,START TIME,END TIME,USAGE,UNITS,NOTES\n"

type csvGbRow struct {
	Type      string `csv:"TYPE"`
	Date      string `csv:"DATE"`
	StartTime string `csv:"START TIME"`
	EndTime   string `csv:"END TIME"`
	Usage     string `csv:"USAGE"`
	Units     string `csv:"UNITS"`
	Notes     string `csv:"NOTES"`
}

func stripToHeader(in string) (string, error) {
	i := strings.Index(in, HEADER)
	if i == -1 {
		return "", errors.New("CSV header not found")
	}
	return in[i:], nil
}

func Read(in string) ([]GBPoint, error) {
	in, err := stripToHeader(in)
	if err != nil {
		return nil, err
	}

	rows := make([]*csvGbRow, 0)
	gocsv.UnmarshalString(in, &rows)

	out := make([]GBPoint, 0)
	for _, row := range rows {
		beginTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", row.Date, row.StartTime), time.Now().Location())
		if err != nil {
			return nil, err
		}
		endTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", row.Date, row.EndTime), time.Now().Location())
		if err != nil {
			return nil, err
		}
		duration := endTime.Sub(beginTime)
		if duration.Minutes() != 59 {
			return nil, errors.New(fmt.Sprintf("Expected duration of 59 seconds, but got '%s'", duration.String()))
		}
		kwhUsed, err := strconv.ParseFloat(row.Usage, 64)
		if err != nil {
			return nil, err
		}
		if row.Units != "kWh" {
			return nil, errors.New(fmt.Sprintf("Expected units of 'kWh', but got '%s'", duration.String()))
		}

		out = append(out, GBPoint{
			T:        beginTime,
			UsageKwh: kwhUsed,
		})
	}
	return out, nil
}
