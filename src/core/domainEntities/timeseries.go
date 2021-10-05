package domainEntities

// HeaderOnly creates a shallow copy of the time series with only headers (i.e. without observations).
// This is useful when creating a new timeseries, as observations are not processed by the endpoint.
func (ts TimeSeries) HeadersOnly() TimeSeries {
	seriesHeaders := []SeriesEntry{}
	for _, series := range ts.TimeSeriesEntry {
		seriesHeaders = append(seriesHeaders, SeriesEntry{
			Header:       series.Header,
			Observations: nil,
		})
	}

	return TimeSeries{
		TimeSeriesType:  ts.TimeSeriesType,
		TimeSeriesEntry: seriesHeaders,
	}
}

type TimeSeries struct {
	TimeSeriesType  string        `json:"tstype"`
	TimeSeriesEntry []SeriesEntry `json:"tseries"`
}

type SeriesEntry struct {
	Header       SeriesHeader        `json:"header"`
	Observations []SeriesObservation `json:"observations"`
}

type SeriesHeader struct {
	Id    HeaderId    `json:"id"`
	Extra HeaderExtra `json:"extra"`
}

type HeaderId struct {
	GliderId  string `json:"gliderID"`  // unique id for research vessel (e.g. association initials + _ + vessel name)
	Parameter string `json:"parameter"` // what has been measured (e.g. "temperature")
}

type HeaderExtra struct {
	Source string `json:"source"` // name of association contributing data
	Name   string `json:"name"`   // name of vessel
}

type SeriesObservation struct {
	Time string          `json:"time"` // TODO: switch to date datatype??
	Body ObservationBody `json:"body"`
}

type ObservationBody struct {
	Pos    ObservationPosition `json:"pos"`
	Value  string              `json:"value"` // the measured value (TODO data type?)
	QcFlag string              `json:"qc_flag"`
}

type ObservationPosition struct {
	Lon    string `json:"lon"` // TODO figure out correct datatype
	Lat    string `json:"lat"`
	Depth  string `json:"depth"` // depth of measurement expressed in meters
	QcFlag string `json:"qc_flag"`
}
