package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sMailund/boatload/src/core/applicationServices"
	"github.com/sMailund/boatload/src/core/domainEntities"
)

var UploadService applicationServices.UploadService

const QC_UNCERTAIN = "2" // quality value: uncertain

const UploadRoute = "/api/upload"

func RegisterRoutes(us applicationServices.UploadService, mux *http.ServeMux) {
	UploadService = us
	mux.HandleFunc(UploadRoute, UploadTimeSeries)
}

func UploadTimeSeries(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	timeSeries, err := readTimeSeriesFromRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = UploadService.UploadTimeSeries(timeSeries)
	if err != nil {
		// TODO: improve error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func readTimeSeriesFromRequest(req *http.Request) (domainEntities.TimeSeries, error) {
	keys := []string{"temperature", "conductivity"}
	req.ParseMultipartForm(32 << 20) // limit your max input length!
	var buf bytes.Buffer
	// in your case file would be fileupload
	file, _, err := req.FormFile("file")
	if err != nil {
		return domainEntities.TimeSeries{}, err
	}
	defer file.Close()
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	observations, err := MapToSeriesObservation(keys, &buf)

	seriesEntries := []domainEntities.SeriesEntry{}
	for _, key := range keys {
		entry := domainEntities.SeriesEntry{
			Header: domainEntities.SeriesHeader{
				Id: domainEntities.HeaderId{
					GliderId:  "TODO",
					Parameter: "TODO",
				},
				Extra: domainEntities.HeaderExtra{
					Source: "TODO",
					Name:   "TODO",
				},
			},
			Observations: observations[key],
		}
		seriesEntries = append(seriesEntries, entry) // TODO fortsett her
	}

	timeSeries := domainEntities.TimeSeries{
		TimeSeriesType:  "TODO",
		TimeSeriesEntry: seriesEntries,
	}

	err = json.NewDecoder(&buf).Decode(&timeSeries)
	if err != nil {
		return domainEntities.TimeSeries{}, err
	}

	return timeSeries, nil
}

func MapToSeriesObservation(keys []string, data io.Reader) (map[string][]domainEntities.SeriesObservation, error) {
	observations := map[string][]domainEntities.SeriesObservation{}
	reader := csv.NewReader(data)

	header, err := reader.Read()
	if err != nil {
		return observations, fmt.Errorf("failed to get reader: %v", err)
	}

	for _, key := range keys {
		observations[key] = []domainEntities.SeriesObservation{}
	}

	keyIndexes := map[string]int{}

	for _, key := range keys {
		keyIndex, err := getKeyIndex(header, key) // TODO better validation
		if err != nil {
			return observations, err
		}
		keyIndexes[key] = keyIndex
	}

	timeIndex, err := getKeyIndex(header, "timestamp")
	if err != nil {
		return observations, err
	}
	lonIndex, err := getKeyIndex(header, "lat")
	if err != nil {
		return observations, err
	}
	latIndex, err := getKeyIndex(header, "lon")
	if err != nil {
		return observations, err
	}
	depthIndex, err := getKeyIndex(header, "depth")
	if err != nil {
		return observations, err
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return observations, fmt.Errorf("failed to read csv rows: %v", err)
	}

	for _, row := range rows {
		timestamp, err := strconv.Atoi(row[timeIndex])
		if err != nil {
			return observations, fmt.Errorf("failed to get timestamp: %v", err)
		}
		position := domainEntities.ObservationPosition{Lat: row[latIndex], Lon: row[lonIndex], Depth: row[depthIndex], QcFlag: QC_UNCERTAIN}
		for _, key := range keys {
			observations[key] = append(observations[key], domainEntities.SeriesObservation{Time: time.Unix(int64(timestamp), 0).Format(time.RFC3339), Body: domainEntities.ObservationBody{Pos: position, Value: row[keyIndexes[key]], QcFlag: QC_UNCERTAIN}})
		}
	}

	return observations, nil
}

func getKeyIndex(row []string, key string) (int, error) {
	for i, v := range row {
		if v == key {
			return i, nil
		}
	}
	return -1, fmt.Errorf("csv header missing column key %v", key)
}
