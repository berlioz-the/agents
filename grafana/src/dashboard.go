package main

import (
	"github.com/berlioz-the/connector-go"
)

func constructDashboard(consumes []berlioz.ConsumesModel) *Dashboard {

	dashboard := &Dashboard{
		Editable:      true,
		GraphTooltip:  0,
		ID:            2,
		Refresh:       false,
		SchemaVersion: 16,
		Style:         "dark",
		Title:         "Berlioz Dashboard",
		Version:       1,
		LastPanelId:   1,
	}

	dashboard.Time.From = "now-30m"
	dashboard.Time.To = "now"

	dashboard.addDashboardRow("Summary")

	chartSummaryCPU := dashboard.addGraph("Service CPU", 12, 9)
	chartSummaryCPU.addGraphYAxes("percent")
	chartSummaryCPU.addGraphYAxes("short")
	for _, consumed := range consumes {
		if consumed.Meta != nil && *consumed.Meta && consumed.Kind == "service" {
			serviceFullName := consumed.Cluster + "-" + consumed.Sector + "-" + consumed.Name
			query := "quantile(0.5, quantile_over_time(0.5, berlioz_cpu_usage{service=\"" + serviceFullName + "\"}[1m]))"
			chartSummaryCPU.addGraphTarget(query, serviceFullName)
		}
	}

	chartSummaryMemory := dashboard.addGraph("Service Memory", 12, 9)
	chartSummaryMemory.addGraphYAxes("decbytes")
	chartSummaryMemory.addGraphYAxes("short")
	for _, consumed := range consumes {
		if consumed.Meta != nil && *consumed.Meta && consumed.Kind == "service" {
			serviceFullName := consumed.Cluster + "-" + consumed.Sector + "-" + consumed.Name
			query := "quantile(0.5, quantile_over_time(0.5, berlioz_memory_usage{service=\"" + serviceFullName + "\"}[1m]))"
			chartSummaryMemory.addGraphTarget(query, serviceFullName)
		}
	}

	dashboard.newLine()

	for _, consumed := range consumes {
		if consumed.Meta != nil && *consumed.Meta && consumed.Kind == "service" {
			addServiceBlock(dashboard, consumed.Cluster, consumed.Sector, consumed.Name)
		}
	}

	return dashboard
}

func addServiceBlock(dashboard *Dashboard, cluster string, sector string, service string) {
	serviceFullName := cluster + "-" + sector + "-" + service
	title := "Service: " + serviceFullName
	// row :=
	dashboard.addDashboardRow(title)

	serviceFilter := "{service=\"" + serviceFullName + "\"}"

	chartTitle := title + " CPU"
	q1 := "berlioz_cpu_usage" + serviceFilter
	chart1 := dashboard.addGraph(chartTitle, 12, 9)
	chart1.addGraphTarget(q1, "{{service}}-{{identity}}")
	chart1.addGraphYAxes("percent")
	chart1.addGraphYAxes("short")

	chartTitle = title + " Memory"
	q2 := "berlioz_memory_usage" + serviceFilter
	chart2 := dashboard.addGraph(chartTitle, 12, 9)
	chart2.addGraphTarget(q2, "{{service}}-{{identity}}")
	chart2.addGraphYAxes("decbytes")
	chart2.addGraphYAxes("short")

	dashboard.newLine()
}

func (dashboard *Dashboard) addGraph(title string, width int, height int) *DashboardPanel {
	panel := dashboard.addPanel(title, width, height)
	panel.Type = "graph"
	panel.Legend.Show = true
	panel.Xaxis.Show = true
	panel.Xaxis.Mode = "time"

	return panel
}

func (panel *DashboardPanel) addGraphTarget(expr string, legend string) *DashboardPanel {
	target := &DashboardTarget{
		Expr:           expr,
		Format:         "time_series",
		IntervalFactor: 1,
		LegendFormat:   legend,
	}
	panel.Targets = append(panel.Targets, target)

	asciiNum := 65 + len(panel.Targets)
	target.RefID = string(asciiNum)

	return panel
}

func (panel *DashboardPanel) addGraphYAxes(format string) *DashboardPanel {
	axes := &DashboardYAxes{
		Format:  format,
		LogBase: 1,
		Show:    true,
	}
	panel.Yaxes = append(panel.Yaxes, axes)
	return panel
}

func (dashboard *Dashboard) addDashboardRow(title string) *DashboardPanel {
	panel := dashboard.addPanel(title, 24, 1)
	panel.Type = "row"
	panel.Collapsed = false

	dashboard.newLine()
	return panel
}

func (dashboard *Dashboard) addPanel(title string, width int, height int) *DashboardPanel {
	dashboard.LastPanelId++
	panel := &DashboardPanel{
		ID: dashboard.LastPanelId,
	}
	panel.GridPos.W = width
	panel.GridPos.H = height
	panel.GridPos.X = dashboard.currentX
	panel.GridPos.Y = dashboard.currentY
	panel.Title = title

	dashboard.currentX += width
	if height > dashboard.currentHeight {
		dashboard.currentHeight = height
	}

	dashboard.Panels = append(dashboard.Panels, panel)
	return panel
}

func (dashboard *Dashboard) newLine() {
	dashboard.currentY += dashboard.currentHeight
	dashboard.currentHeight = 0
	dashboard.currentX = 0
}
