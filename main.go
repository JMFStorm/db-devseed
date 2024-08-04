package main

import (
        "database/sql"
        "fmt"
        "os"
        "bufio"
        "log"
        "strings"
        "path/filepath"

        _ "github.com/lib/pq"
)

const (
        errorGeneral = 1
        errorCommandNotFound = 27
        errorCommandLine = 64
        errorDataFormat = 65
        errorSystem = 71
        configFileName = "dbds.cfg"
        scriptsDirName = "dbds_scripts"
        dropTablesFile = "schemas_drop.sql"
        createTablesFile = "schemas_create.sql"
        generateIndexesFile = "schemas_indexes.sql"
        populateTablesFile = "schemas_populate.sql"
        defaultConfigContent = "connectionString:\"\"\ndbType:\"postgres\""
)

func fileExists(path string) bool {
        _, err := os.Stat(path)
        if os.IsNotExist(err) {
            return false
        }
        return err == nil
}

func searchSqlScrip(filename string, scriptsPath string) (sqlPath string) {
        sqlPath = filepath.Join(scriptsPath, filename)
        if !fileExists(sqlPath) {
                fmt.Printf("Could not find script from path: %s\n", sqlPath)
                os.Exit(1)
        }
        return sqlPath
}

func executeSqlQueryFromFile(db *sql.DB, sqlPath string) {
        sqlBytes, err := os.ReadFile(sqlPath)
        if err != nil {
                log.Fatal(err)
        }
        queryStr := string(sqlBytes)
        _, err = db.Query(queryStr)
        if err != nil {
                log.Fatal(err)
        }
}

func rebuildDatabases(db *sql.DB, scriptsPath string) {
        dropTablesPath := searchSqlScrip(dropTablesFile, scriptsPath)
        createTablesPath := searchSqlScrip(createTablesFile, scriptsPath)
        generateIndexesPath := searchSqlScrip(generateIndexesFile, scriptsPath)
        populateTablesPath := searchSqlScrip(populateTablesFile, scriptsPath)

        executeSqlQueryFromFile(db, dropTablesPath)
        fmt.Println("Tables dropped")

        executeSqlQueryFromFile(db, createTablesPath)
        fmt.Println("Tables created")

        executeSqlQueryFromFile(db, generateIndexesPath)
        fmt.Println("Indexes created")

        executeSqlQueryFromFile(db, populateTablesPath)
        fmt.Println("Tables populated")
}

func splitFirst(s string, delimiter rune) (string, string) {
    index := strings.IndexRune(s, delimiter)
    if index == -1 {
        return s, ""
    }
    left := s[:index]
    right := s[index+1:]
    return left, right
}

func createDbdsScript(scriptFilepath string, fileText string) {
        scriptFile, err := os.Create(scriptFilepath)
        if err != nil {
                fmt.Printf("Error creating file %s: %v\n", scriptFilepath, err)
                os.Exit(errorSystem)
        }
        defer scriptFile.Close()
        fmt.Printf("Sql file '%s' created\n", scriptFilepath)

        _, err = scriptFile.WriteString(fileText)
        if err != nil {
                fmt.Printf("Error writing to file %s: %v\n", scriptFilepath, err)
                os.Exit(errorSystem)
        } 
        fmt.Printf("Text written to file '%s'\n", fileText)
}

func deleteFileIfExists(filepath string) error {
    _, err := os.Stat(filepath)
    if os.IsNotExist(err) {
        return nil
    }
    err = os.Remove(filepath)
    if err != nil {
        return fmt.Errorf("Failed to delete file: %w", err)
    }
    return nil
}

func initConfig() {
        dropTablesFilepath := scriptsDirName + "/" + dropTablesFile
        createTablesFilepath := scriptsDirName + "/" + createTablesFile
        generateIndexesFilepath := scriptsDirName + "/" + generateIndexesFile
        populateTablesFilepath := scriptsDirName + "/" + populateTablesFile

        if fileExists(configFileName) {
                reader := bufio.NewReader(os.Stdin)
                for {
                        fmt.Printf("Configuration file '%s' already exists.\n", configFileName)
                        fmt.Printf("Force delete existing dbds config and sql scripts under dbds_script? (y/n)\n")
                        input, err := reader.ReadString('\n')
                        if err != nil {
                                fmt.Println("Error reading input:", err)
                                os.Exit(errorGeneral)
                        }
                        input = strings.TrimSpace(strings.ToLower(input))
                        if input == "y" {
                                deleteFileIfExists(dropTablesFilepath)
                                deleteFileIfExists(createTablesFilepath)
                                deleteFileIfExists(generateIndexesFilepath)
                                deleteFileIfExists(populateTablesFilepath)
                                break
                        } else if input == "n" {
                                fmt.Printf("dbds terminated\n")
                                os.Exit(0)
                        } else {
                                fmt.Println("Invalid input. Please enter 'y' or 'n'.")
                        }
                }
        }
        err := os.Mkdir(scriptsDirName, 0755) // 0755 is the permission mode (rwxr-xr-x)
        if err != nil && !os.IsExist(err) {
                fmt.Printf("Error creating directory: %v\n", err)
                os.Exit(errorSystem)
        }

        createDbdsScript(dropTablesFilepath, "-- Insert your 'DROP TABLE' statements here")
        createDbdsScript(createTablesFilepath, "-- Insert your 'CREATE TABLE' statements here")
        createDbdsScript(generateIndexesFilepath, "-- Insert your 'CREATE INDEX' statements here")
        createDbdsScript(populateTablesFilepath, "-- Insert your 'INSERT INTO' statements here")

        file, err := os.Create(configFileName)
        if err != nil {
                fmt.Println("Error creating file:", err)
                os.Exit(errorSystem)
        }
        defer file.Close()
        writer := bufio.NewWriter(file)
        _, err = writer.WriteString(defaultConfigContent)
        if err != nil {
                fmt.Println("Error writing to file:", err)
                os.Exit(errorSystem)
        }
        writer.Flush()
        fmt.Println("Config file 'dbds.cfg' created.")
        fmt.Println("Please fill out the required information to the configuraton file and the sql scripts.")
}

func isEmptyOrWhitespace(s string) bool {
    trimmed := strings.TrimSpace(s)
    return len(trimmed) == 0
}

func trimQuotes(s string) string {
    if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
        return s[1 : len(s)-1]
    }
    return s
}

func readConfig() (connectionString string, dbType string) {
        file, err := os.Open(configFileName)
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                line := scanner.Text()
                key, value := splitFirst(line, ':')
                if key == "connectionString" {
                        connectionString = trimQuotes(value)
                } else if key == "dbType" {
                        dbType = trimQuotes(value)
                } else if !isEmptyOrWhitespace(line) {
                        if !strings.Contains(line, ":") {
                                fmt.Printf("Invalid config format. Use key:\"value\". (From line: '%s')\n", line)
                                os.Exit(errorDataFormat)
                        } 
                        fmt.Printf("Invalid config key '%s' found.\n", key)
                        os.Exit(errorDataFormat)
                }
        }
        if err := scanner.Err(); err != nil {
                fmt.Println("Error reading file:", err)
                os.Exit(errorSystem)
        }
        if connectionString == "" {
                fmt.Println("Did not find configuration value for 'connectionString'.")
                os.Exit(errorDataFormat)
        }
        if dbType == "" {
                fmt.Println("Did not find configuration value for 'dbType'.")
                os.Exit(errorDataFormat)
        }
        return connectionString, dbType
}

func main() {
        if len(os.Args) < 2 {
                fmt.Println("dbds usage:")
                fmt.Println("[dbds init]    -> Initialize new configuration and create base sql scripts. Fill the required data in them.")
                fmt.Println("[dbds rebuild] -> Start database seed with existing config and sql scripts.")
                os.Exit(1)
        }

        command := os.Args[1]
        switch command {
        case "init":
                initConfig()
        case "rebuild":
                connStr, dbType := readConfig()
                if dbType != "postgres" {
                        fmt.Printf("Found config dbType: '%s'. However, postgres is the only supported type.", dbType)
                        os.Exit(1)
                }
                db, err := sql.Open(dbType, connStr)
                if err != nil {
                        log.Fatal(err)
                }
                defer db.Close()
                rebuildDatabases(db, scriptsDirName)
                fmt.Println("dbds rebuild completed")
        default:
                fmt.Println("Unknown command:", command)
                os.Exit(errorCommandNotFound)
        }
}
