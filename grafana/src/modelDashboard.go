package main

import ()

// Dashboard is a model for Grafana Dashboard model
type Dashboard struct {
	LastPanelId   int
	currentY      int
	currentX      int
	currentHeight int
	Annotations   struct {
		List []DashboardAnnotation `json:"list"`
	} `json:"annotations"`
	Editable      bool              `json:"editable"`
	GnetID        interface{}       `json:"gnetId"`
	GraphTooltip  int               `json:"graphTooltip"`
	ID            int               `json:"id"`
	Links         []interface{}     `json:"links"`
	Panels        []*DashboardPanel `json:"panels"`
	Refresh       bool              `json:"refresh"`
	SchemaVersion int               `json:"schemaVersion"`
	Style         string            `json:"style"`
	Tags          []interface{}     `json:"tags"`
	Templating    struct {
		List []interface{} `json:"list"`
	} `json:"templating"`
	Time struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time"`
	Timepicker struct {
		RefreshIntervals []string `json:"refresh_intervals"`
		TimeOptions      []string `json:"time_options"`
	} `json:"timepicker"`
	Timezone string `json:"timezone"`
	Title    string `json:"title"`
	UID      string `json:"uid"`
	Version  int    `json:"version"`
}

// DashboardAnnotation TBD
type DashboardAnnotation struct {
	BuiltIn    int    `json:"builtIn"`
	Datasource string `json:"datasource"`
	Enable     bool   `json:"enable"`
	Hide       bool   `json:"hide"`
	IconColor  string `json:"iconColor"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type DashboardPanel struct {
	Collapsed bool `json:"collapsed,omitempty"`
	GridPos   struct {
		H int `json:"h"`
		W int `json:"w"`
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"gridPos"`
	ID          int           `json:"id"`
	Panels      []interface{} `json:"panels,omitempty"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	AliasColors struct {
	} `json:"aliasColors,omitempty"`
	Bars       bool        `json:"bars,omitempty"`
	DashLength int         `json:"dashLength,omitempty"`
	Dashes     bool        `json:"dashes,omitempty"`
	Datasource interface{} `json:"datasource,omitempty"`
	Fill       int         `json:"fill,omitempty"`
	Legend     struct {
		Avg     bool `json:"avg"`
		Current bool `json:"current"`
		Max     bool `json:"max"`
		Min     bool `json:"min"`
		Show    bool `json:"show"`
		Total   bool `json:"total"`
		Values  bool `json:"values"`
	} `json:"legend,omitempty"`
	Lines           bool               `json:"lines,omitempty"`
	Linewidth       int                `json:"linewidth,omitempty"`
	Links           []interface{}      `json:"links,omitempty"`
	NullPointMode   string             `json:"nullPointMode,omitempty"`
	Percentage      bool               `json:"percentage,omitempty"`
	Pointradius     int                `json:"pointradius,omitempty"`
	Points          bool               `json:"points,omitempty"`
	Renderer        string             `json:"renderer,omitempty"`
	SeriesOverrides []interface{}      `json:"seriesOverrides,omitempty"`
	SpaceLength     int                `json:"spaceLength,omitempty"`
	Stack           bool               `json:"stack,omitempty"`
	SteppedLine     bool               `json:"steppedLine,omitempty"`
	Targets         []*DashboardTarget `json:"targets,omitempty"`
	Thresholds      []interface{}      `json:"thresholds,omitempty"`
	TimeFrom        interface{}        `json:"timeFrom,omitempty"`
	TimeShift       interface{}        `json:"timeShift,omitempty"`
	Tooltip         struct {
		Shared    bool   `json:"shared"`
		Sort      int    `json:"sort"`
		ValueType string `json:"value_type"`
	} `json:"tooltip,omitempty"`
	Xaxis struct {
		Buckets interface{}   `json:"buckets"`
		Mode    string        `json:"mode"`
		Name    interface{}   `json:"name"`
		Show    bool          `json:"show"`
		Values  []interface{} `json:"values"`
	} `json:"xaxis,omitempty"`
	Yaxes []*DashboardYAxes `json:"yaxes,omitempty"`
	Yaxis struct {
		Align      bool        `json:"align"`
		AlignLevel interface{} `json:"alignLevel"`
	} `json:"yaxis,omitempty"`
}

// DashboardTarget TBD
type DashboardTarget struct {
	Expr           string `json:"expr"`
	Format         string `json:"format"`
	Hide           bool   `json:"hide"`
	Instant        bool   `json:"instant"`
	Interval       string `json:"interval"`
	IntervalFactor int    `json:"intervalFactor"`
	LegendFormat   string `json:"legendFormat"`
	RefID          string `json:"refId"`
}

// DashboardYAxes tbd
type DashboardYAxes struct {
	Format  string      `json:"format"`
	Label   interface{} `json:"label"`
	LogBase int         `json:"logBase"`
	Max     interface{} `json:"max"`
	Min     interface{} `json:"min"`
	Show    bool        `json:"show"`
}
