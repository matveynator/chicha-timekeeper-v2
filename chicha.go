package main

import (
	"fmt"
	"chicha/packages/database"
	"chicha/packages/collector"
)


func main() {
  Database.CreateDB()
	Database.UpdateDB()
  Collector.StartDataCollector()
	fmt.Println("Chicha worked!")


}

