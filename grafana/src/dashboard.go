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

	chartTitle = title + " Memory"
	q2 := "berlioz_memory_usage" + serviceFilter
	chart2 := dashboard.addGraph(chartTitle, 12, 9)
	chart2.addGraphTarget(q2, "{{service}}-{{identity}}")

	dashboard.newLine()
}

func (dashboard *Dashboard) addGraph(title string, width int, height int) *DashboardPanel {
	panel := dashboard.addPanel(title, width, height)
	panel.Type = "graph"
	panel.Legend.Show = true

	// dashboard.newLine()
	return panel
}

func (panel *DashboardPanel) addGraphTarget(expr string, legend string) *DashboardPanel {
	target := &DashboardTarget{
		Expr:           expr,
		Format:         "time_series",
		IntervalFactor: 1,
		LegendFormat:   legend,
		RefID:          "A",
	}
	panel.Targets = append(panel.Targets, target)
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
