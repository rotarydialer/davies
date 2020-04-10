package main

import (
    _ "github.com/denisenkom/go-mssqldb"
    "database/sql"
    "context"
    "log"
    "fmt"
    "os"
    "flag"
    "path/filepath"
    "strings"
    "strconv"
)

var db *sql.DB

func Usage() {
    fmt.Println("This command requires the following arguments to initialize it for use with davies:")
    flag.PrintDefaults()
}

func main() {
    dirPtr := flag.String("dir", ".", "Working directory to initialize (root of davies repo)")
    urlPtr := flag.String("url", "", "Hostname of MS SQL server in format: sqlserver://<hostname>:<port>")
    dbPtr := flag.String("db", "", "Name of database to initialize")
    userPtr := flag.String("user", "sa", "Database user to connect with")
    
    // SETUP
    // get db password from ~/.mssqlpass file, if user matches
        // ! - only read this file if permissions are correct: -rw-------
        // $ stat -f%p ~/.mssqlpass
        //   100600
    
    // INITIALIZATION
    // creates a ./scripts directory
    // connects to the db and creates a davies_evo schema
    // git commits similar to SEM: https://github.com/mbryzek/schema-evolution-manager/blob/master/bin/sem-init#L66
        // ! - see go-git: https://git-scm.com/book/en/v2/Appendix-B%3A-Embedding-Git-in-your-Applications-go-git

    flag.Parse()
    flag.Usage = Usage

    CheckArgs(*dirPtr, *urlPtr, *dbPtr, *userPtr)

    host, port := ParseDbUrl(*urlPtr)

    var user = *userPtr
    var password = "" // TODO: get this from file
    var database = *dbPtr

    fmt.Println("host, port, user, database:", host, port, user, database)

    var err error

    // Create connection string
    connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
        host, user, password, port, database)

    fmt.Println(connString)

    // Create connection pool
    db, err = sql.Open("sqlserver", connString)
    if err != nil {
        log.Fatal("Error creating connection pool: " + err.Error())
    }
    log.Printf("Connected!\n")
    
    SelectVersion()

    InitDb()

    // Close the database connection pool after program executes
    defer db.Close()
    log.Printf("Db connection closed.\n")

    os.Exit(0)

}

// check all arguments for valid parameters
func CheckArgs(dir, url, db, user string) {

    if len(os.Args) == 1 {
        Usage()
        os.Exit(1)
    }

    // validate directory
    dir,err := filepath.Abs(dir)
    if err == nil {
      // fmt.Println(" -> dir:", dir)
    }
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        // path/to/whatever does not exist
        msg := fmt.Sprintf("Directory '%s' does not exist!", dir)
        fmt.Println(msg)
        os.Exit(1)
    }

    // TODO: check/validate these
    fmt.Println(" dir:", dir)
    fmt.Println(" url:", url)
    fmt.Println(" db:", db)
    fmt.Println(" user:", user)
}

func ParseDbUrl(url string) (string, int) {
    urlSplit := strings.Split(strings.TrimPrefix(url, "sqlserver://"), ":")

    server := urlSplit[0]
    port, err := strconv.Atoi(urlSplit[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }

    return server, port
}

// Gets and prints SQL Server version
func SelectVersion() {
    // Use background context
    ctx := context.Background()

    // Ping database to see if it's still alive.
    // Important for handling network issues and long queries.
    err := db.PingContext(ctx)
    if err != nil {
        log.Fatal("Error pinging database: " + err.Error())
    }

    var result string

    // Run query and scan for result
    err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
    if err != nil {
        log.Fatal("Scan failed:", err.Error())
    }
    fmt.Printf("%s\n", result)
}

func InitDb(){
    // TODO: dogfood this: move to sql scripts
    
    ctx := context.Background()

    err := db.PingContext(ctx)
    if err != nil {
        log.Fatal("Error pinging database: " + err.Error())
    }

    // Run query and scan for result
    rows, err := db.QueryContext(ctx, "SELECT name, schema_id FROM sys.schemas WHERE name = 'davies_evo'")
    // TODO: fix this check; it's not working
    // if err == nil {
    //     log.Fatal("Schmea 'davies_evo' already exists. Has this db already been initialized?")
    //     os.Exit(1)
    // }
    fmt.Println("rows --- ", rows)

    fmt.Println("Creating 'davies_evo' schema...")

    _, schErr := db.ExecContext(ctx, "CREATE SCHEMA davies_evo;")

    if schErr != nil {
        log.Fatal("Create schema failed: ", schErr.Error())
        os.Exit(1)
    }
    fmt.Println("Schema created.")

    _, tblErr := db.ExecContext(ctx, "CREATE TABLE davies_evo.scripts;")
    if tblErr != nil {
        log.Fatal("Error creating scripts table: ", tblErr.Error())
    }

}