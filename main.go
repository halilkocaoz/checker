package main

import (
	"fmt"
	"os"
	"time"

	"github.com/halilkocaoz/upsmo-checker/model"
	"github.com/halilkocaoz/upsmo-checker/storage"
)

var (
	region                string = os.Getenv("REGION")
	pureMonitorsStatement string = `SELECT "ID", 
	"Host", 
	"Method", 
	"Region", 
	"IntervalMs", 
	"TimeoutMs", 
	"CreatedAt" 
	FROM "Monitors" 
	%s 
	ORDER BY "CreatedAt"`
	byRegion                        string = fmt.Sprintf(pureMonitorsStatement, `WHERE ("DeletedAt" IS NULL AND "Region" = '%s')`)
	byRegionAndGreaterThanCreatedAt string = fmt.Sprintf(pureMonitorsStatement, `WHERE ("DeletedAt" IS NULL AND "Region" = '%s' AND "CreatedAt" > '%s')`)
)

func main() {
	if !(len(region) > 0) {
		panic("Program cannot be executed without region")
	}

existMonitors: // get exist monitors and put them into process
	statement := fmt.Sprintf(byRegion, region)
	monitors, _ := getMonitorsByStatement(statement)
	if len(monitors) == 0 {
		fmt.Printf("There is no monitor with region: %v\n", region)
		time.Sleep(10 * time.Second)
		goto existMonitors
	}
	processMonitors(monitors)
	lastMonitor := monitors[len(monitors)-1]

newMonitors: // get new monitors(according to lastMonitor.CreatedAt) and put them into process
	statement = fmt.Sprintf(byRegionAndGreaterThanCreatedAt, region, lastMonitor.CreatedAt)
	monitors, _ = getMonitorsByStatement(statement)
	if len(monitors) > 0 {
		processMonitors(monitors)
		lastMonitor = monitors[len(monitors)-1]
	}
	time.Sleep(30 * time.Second)
	goto newMonitors
}

func processMonitors(monitors []*model.Monitor) {
	for _, mn := range monitors {
		go mn.Process()
	}
}

func getMonitorsByStatement(statement string) ([]*model.Monitor, error) {
	monitors := make([]*model.Monitor, 0)
	db, _ := storage.UpsMoDBConn()
	defer db.Close()

	monitorRows, err := db.Query(statement)
	if err != nil {
		return monitors, err
	}
	defer monitorRows.Close()

	for monitorRows.Next() {
		mn := new(model.Monitor)

		err = monitorRows.Scan(&mn.ID,
			&mn.Host,
			&mn.Method,
			&mn.Region,
			&mn.IntervalMs,
			&mn.TimeoutMs,
			&mn.CreatedAt)
		monitors = append(monitors, mn)
	}

	return monitors, err
}
