package main

import (
	"fmt"
	"chicha/packages/db"
)


func main() {
  Db.CreateDB()
	Db.UpdateDB()

	fmt.Println("Chicha worked!")


}

