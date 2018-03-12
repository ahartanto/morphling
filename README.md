# Morphling
Morphling is a SQL wrapper for database replication topology. Morphling helps you to determine which connection to be
used when execute query to the sql database. **It's not an ORM**



#### Example

```go
package main

import (
    "fmt"
    "github.com/ahartanto/morphling"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

    var urlFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local"

    // Prepare data source for master connection
    masterUser := "user"
    masterPassword := "password"
    masterHost := "masterhost"
    masterPort := "3306"
    masterDBName := "dbname"
    masterCharset := "utf8"
    dataSourceMaster := fmt.Sprintf(urlFormat, masterUser, masterPassword,masterHost, masterPort, masterDBName, masterCharset)

    // Prepare data source for slave connection
    slaveUser := "user"
    slavePassword := "password"
    slaveHost := "slavehost"
    slavePort := "3306"
    slaveDBName := "dbname"
    slaveCharset := "utf8"
    dataSourceSlave := fmt.Sprintf(urlFormat, slaveUser, slavePassword,slaveHost, slavePort, slaveDBName, slaveCharset)

    // Get database handler
    db, err := morphling.OpenConnection(morphling.MySQLDriver, dataSourceMaster, dataSourceSlave)
    if err != nil {
        panic(fmt.Errorf("failed dial database slave : %v", err))
    }

    // Then you can do the query as golang sql package does
    var id int
    q := fmt.Sprintf("SELECT id FROM accounts LIMIT 1")
    err = db.QueryRow(q).Scan(&id)
    fmt.Printf(id)

}
```