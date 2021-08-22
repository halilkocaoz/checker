package main

import (
	"github.com/halilkocaoz/upsmo-checker/model"
	"github.com/halilkocaoz/upsmo-checker/storage"
)

func main() {

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
