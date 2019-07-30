package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zorkian/go-datadog-api"
)

func TestDashboardCreateAndDeleteAdvancesTimeseries(t *testing.T) {
	expected := getTestDashboard(createAdvancedTimeseriesGraph)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}

	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// now try to fetch it freshly and compare it again
	actual, err = client.GetDashboard(*actual.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actual, expected)

}

func TestDashboardUpdate(t *testing.T) {
	expected := getTestDashboard(createGraph)
	board, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}

	defer cleanUpDashboard(t, *board.Id)
	board.Title = datadog.String("___New-Test-Board___")

	if err := client.UpdateDashboard(board); err != nil {
		t.Fatalf("Updating a dashboard failed when it shouldn't: %s", err)
	}

	actual, err := client.GetDashboard(*board.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't: %s", err)
	}

	assertDashboardEquals(t, actual, board)
}

func TestDashboardGet(t *testing.T) {
	boards, err := client.GetDashboards()
	if err != nil {
		t.Fatalf("Retrieving dashboards failed when it shouldn't: %s", err)
	}

	num := len(boards)
	board := createTestDashboard(t)
	defer cleanUpDashboard(t, *board.Id)

	boards, err = client.GetDashboards()
	if err != nil {
		t.Fatalf("Retrieving dashboards failed when it shouldn't: %s", err)
	}

	if num+1 != len(boards) {
		t.Fatalf("Number of dashboards didn't match expected: %d != %d", len(boards), num+1)
	}
}

func TestDashboardCreateWithCustomGraph(t *testing.T) {
	expected := getTestDashboard(createCustomGraph)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}

	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// now try to fetch it freshly and compare it again
	actual, err = client.GetDashboard(*actual.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actual, expected)
}

// Create a dashboard with a metrics query graph.
func TestDashboardCreateWithQuery(t *testing.T) {
	expected := getTestDashboard(createGraph)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}
	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// now try to fetch it freshly and compare it again
	actual, err = client.GetDashboard(*actual.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actual, expected)
}

// Create a dashboard with a log query graph.
func TestDashboardCreateWithLogOrApmQuery(t *testing.T) {
	expected := getTestDashboard(createGraphWithLogQuery)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}
	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// now try to fetch it freshly and compare it again
	actual, err = client.GetDashboard(*actual.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actual, expected)
}

// Create a dashboard with a process query graph.
func TestDashboardCreateWithProcessQuery(t *testing.T) {
	expected := getTestDashboard(createGraphWithProcessQuery)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}
	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// now try to fetch it freshly and compare it again
	actual, err = client.GetDashboard(*actual.Id)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actual, expected)
}

func TestDashboardGetWithNewId(t *testing.T) {
	expected := getTestDashboard(createGraph)
	// create the dashboard and compare it
	actual, err := client.CreateDashboard(expected)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't. (%s)", err)
	}

	defer cleanUpDashboard(t, *actual.Id)

	assertDashboardEquals(t, actual, expected)

	// try to fetch it freshly using the new id format and compare it again
	actualWithNewId, err := client.GetDashboard(*actual.NewId)
	if err != nil {
		t.Fatalf("Retrieving a dashboard failed when it shouldn't. (%s)", err)
	}
	assertDashboardEquals(t, actualWithNewId, expected)

	// the ids are equal whether fetching using the old or the new id
	assert.Equal(t, *actualWithNewId.Id, *actual.Id)

	// try to fetch it freshly using a string, but with a wrong value
	actual, err = client.GetDashboard("random_string")
	if assert.NotNil(t, err) {
		// it should not fail because of the id format
		assert.NotContains(t, err.Error(), "unsupported id type")
		assert.Contains(t, err.Error(), "404")
	}

	// try to fetch it freshly using a boolean
	actual, err = client.GetDashboard(true)
	if assert.NotNil(t, err) {
		// it should fail because of the id format
		assert.Contains(t, err.Error(), "unsupported id type")
	}

	// try to fetch it freshly using a float64
	actual, err = client.GetDashboard(5.5)
	if assert.NotNil(t, err) {
		// it should fail because of the id format
		assert.Contains(t, err.Error(), "unsupported id type")
	}
}
func createTestDashboard(t *testing.T) *datadog.Dashboard {
	board := getTestDashboard(createGraph)
	board, err := client.CreateDashboard(board)
	if err != nil {
		t.Fatalf("Creating a dashboard failed when it shouldn't: %s", err)
	}
	return board
}

func getTestDashboard(createGraph func() []datadog.Graph) *datadog.Dashboard {
	return &datadog.Dashboard{
		Title:             datadog.String("___Test-Board___"),
		Description:       datadog.String("Testboard description"),
		TemplateVariables: []datadog.TemplateVariable{},
		Graphs:            createGraph(),
	}
}

func cleanUpDashboard(t *testing.T, id int) {
	if err := client.DeleteDashboard(id); err != nil {
		t.Fatalf("Deleting a dashboard failed when it shouldn't. Manual cleanup needed. (%s)", err)
	}

	deletedBoard, err := client.GetDashboard(id)
	if deletedBoard != nil {
		t.Fatal("Dashboard hasn't been deleted when it should have been. Manual cleanup needed.")
	}

	if err == nil {
		t.Fatal("Fetching deleted dashboard didn't lead to an error. Manual cleanup needed.")
	}
}

func createGraph() []datadog.Graph {
	gd := &datadog.GraphDefinition{}
	gd.SetViz("timeseries")

	r := gd.Requests
	gd.Requests = append(r, datadog.GraphDefinitionRequest{
		Query:   datadog.String("avg:system.mem.free{*}"),
		Stacked: datadog.Bool(false),
	})

	graph := datadog.Graph{
		Title:      datadog.String("Mandatory graph"),
		Definition: gd,
	}

	graphs := []datadog.Graph{}
	graphs = append(graphs, graph)
	return graphs
}

func createAdvancedTimeseriesGraph() []datadog.Graph {
	gd := &datadog.GraphDefinition{}
	gd.SetViz("timeseries")

	r := gd.Requests
	gd.Requests = append(r, datadog.GraphDefinitionRequest{
		Query:   datadog.String("avg:system.mem.free{*}"),
		Stacked: datadog.Bool(false),
		Type:    datadog.String("bars"),
		Style:   &datadog.GraphDefinitionRequestStyle{Palette: datadog.String("warm")},
	})
	graph := datadog.Graph{Title: datadog.String("Custom type and style graph"), Definition: gd}

	graphs := []datadog.Graph{}
	graphs = append(graphs, graph)
	return graphs
}

func createCustomGraph() []datadog.Graph {
	gd := &datadog.GraphDefinition{}
	gd.SetViz("timeseries")

	r := gd.Requests
	gd.Requests = append(r, datadog.GraphDefinitionRequest{
		Query:      datadog.String("( sum:system.mem.used{*} / sum:system.mem.free{*} ) * 100"),
		Stacked:    datadog.Bool(false),
		Aggregator: datadog.String("avg"),
		Metadata: map[string]datadog.GraphDefinitionMetadata{
			"(sum:system.mem.used{*}/sum:system.mem.free{*})*100": {
				Alias: datadog.String("mem_used_ratio"),
			},
		},
		ConditionalFormats: []datadog.DashboardConditionalFormat{
			{
				Comparator: datadog.String(">"),
				Value:      datadog.JsonNumber(json.Number("99.9")),
				Palette:    datadog.String("white_on_green")},
			{
				Comparator: datadog.String(">="),
				Value:      datadog.JsonNumber(json.Number("99")),
				Palette:    datadog.String("white_on_yellow")},
			{
				Comparator: datadog.String("<"),
				Value:      datadog.JsonNumber(json.Number("99")),
				Palette:    datadog.String("white_on_red")}}})

	graph := datadog.Graph{
		Title:      datadog.String("Mandatory graph 2"),
		Definition: gd,
	}

	graphs := []datadog.Graph{}
	graphs = append(graphs, graph)
	return graphs
}

func createGraphWithLogQuery() []datadog.Graph {
	gd := &datadog.GraphDefinition{}
	gd.SetViz("timeseries")

	r := gd.Requests
	gd.Requests = append(r, datadog.GraphDefinitionRequest{
		Aggregator: datadog.String("avg"),
		Type:       datadog.String("line"),
		LogQuery: &datadog.GraphApmOrLogQuery{
			Index: datadog.String("avg"),
			Compute: &datadog.GraphApmOrLogQueryCompute{
				Aggregation: datadog.String("count"),
				Facet:       datadog.String("host"),
				Interval:    datadog.Int(300000),
			},
			Search:  &datadog.GraphApmOrLogQuerySearch{},
			GroupBy: []datadog.GraphApmOrLogQueryGroupBy{},
		},
	})

	graph := datadog.Graph{
		Title:      datadog.String("Mandatory graph 3"),
		Definition: gd,
	}

	graphs := []datadog.Graph{}
	graphs = append(graphs, graph)
	return graphs
}

func createGraphWithProcessQuery() []datadog.Graph {
	gd := &datadog.GraphDefinition{}
	gd.SetViz("timeseries")

	r := gd.Requests
	gd.Requests = append(r, datadog.GraphDefinitionRequest{
		Aggregator: datadog.String("avg"),
		Type:       datadog.String("line"),
		ProcessQuery: &datadog.GraphProcessQuery{
			Metric: datadog.String("process.stat.cpu.total_pct"),
			Limit:  datadog.Int(10),
		},
	})

	graph := datadog.Graph{
		Title:      datadog.String("Mandatory graph 4"),
		Definition: gd,
	}

	graphs := []datadog.Graph{}
	graphs = append(graphs, graph)
	return graphs
}

func assertDashboardEquals(t *testing.T, actual, expected *datadog.Dashboard) {
	if *actual.Title != *expected.Title {
		t.Errorf("Dashboard title does not match: %s != %s", *actual.Title, *expected.Title)
	}
	if *actual.Description != *expected.Description {
		t.Errorf("Dashboard description does not match: %s != %s", *actual.Description, *expected.Description)
	}
	if len(actual.Graphs) != len(expected.Graphs) {
		t.Errorf("Number of Dashboard graphs does not match: %d != %d", len(actual.Graphs), len(expected.Graphs))
	}
	if len(actual.TemplateVariables) != len(expected.TemplateVariables) {
		t.Errorf("Number of Dashboard template variables does not match: %d != %d", len(actual.TemplateVariables), len(expected.TemplateVariables))
	}
}
