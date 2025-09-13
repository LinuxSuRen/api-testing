# Database Connection Leak Fix

## Problem Description
The ORM database store in api-testing has a connection leak issue where MySQL client connections keep increasing over time, especially when switching between databases. This eventually leads to "too many connections" errors.

## Root Cause Analysis

### Issue Location
The problem is in the `atest-ext-store-orm` extension, specifically in the `getClientWithDatabase` method in `pkg/server.go`.

### Current Problematic Code
```go
var dbCache = make(map[string]*gorm.DB)
var dbNameCache = make(map[string]string)

func (s *dbserver) getClientWithDatabase(ctx context.Context, dbName string) (dbQuery DataQuery, err error) {
    // ... store and database logic ...
    
    var ok bool
    var db *gorm.DB
    if db, ok = dbCache[store.Name]; (ok && db != nil && dbNameCache[store.Name] != database) || !ok {
        if db, err = createDB(store.Username, store.Password, store.URL, database, driver); err == nil {
            dbCache[store.Name] = db  // <- PROBLEM: Overwrites without closing old connection
            dbNameCache[store.Name] = database
        } else {
            return
        }
    }
    // ...
}
```

### Issues Identified
1. **Cache Key Problem**: Using only `store.Name` as cache key means switching databases creates new connections
2. **No Connection Cleanup**: Old connections are overwritten without being properly closed
3. **Missing Connection Pool Configuration**: No limits set on connection pool size
4. **Improper Cache Logic**: The condition creates new connections instead of reusing existing ones per database

## Solution

### 1. Fix Cache Key Strategy
Use composite keys that include both store name and database name:

```go
// Generate a unique cache key for store + database combination
func generateCacheKey(storeName, database string) string {
    return fmt.Sprintf("%s:%s", storeName, database)
}
```

### 2. Proper Connection Management
```go
var dbCache = make(map[string]*gorm.DB)
var cacheMutex = sync.RWMutex{}

func (s *dbserver) getClientWithDatabase(ctx context.Context, dbName string) (dbQuery DataQuery, err error) {
    store := remote.GetStoreFromContext(ctx)
    if store == nil {
        err = errors.New("no connect to database")
        return
    }

    database := dbName
    if database == "" {
        if v, ok := store.Properties["database"]; ok && v != "" {
            database = v
        }
    }

    driver := DialectorMySQL
    if v, ok := store.Properties["driver"]; ok && v != "" {
        driver = v
    }

    // Use composite cache key
    cacheKey := generateCacheKey(store.Name, database)
    
    cacheMutex.RLock()
    db, exists := dbCache[cacheKey]
    cacheMutex.RUnlock()

    if !exists || db == nil {
        log.Printf("Creating new connection for store[%s] database[%s]", store.Name, database)
        
        if db, err = createDBWithPool(store.Username, store.Password, store.URL, database, driver); err == nil {
            cacheMutex.Lock()
            dbCache[cacheKey] = db
            cacheMutex.Unlock()
        } else {
            return
        }
    }

    dbQuery = NewCommonDataQuery(GetInnerSQL(driver), db)
    return
}
```

### 3. Enhanced createDB Function with Connection Pool
```go
func createDBWithPool(user, password, address, database, driver string) (db *gorm.DB, err error) {
    var dialector gorm.Dialector
    var dsn string
    
    switch driver {
    case DialectorMySQL, "", "greptime":
        if !strings.Contains(address, ":") {
            address = fmt.Sprintf("%s:%d", address, 3306)
        }
        dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", user, password, address, database)
        dialector = mysql.Open(dsn)
    case "sqlite":
        dsn = fmt.Sprintf("%s.db", database)
        dialector = sqlite.Open(dsn)
    case DialectorPostgres:
        obj := strings.Split(address, ":")
        host, port := obj[0], "5432"
        if len(obj) > 1 {
            port = obj[1]
        }
        dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, database, port)
        dialector = postgres.Open(dsn)
    case "tdengine":
        dsn = fmt.Sprintf("%s:%s@ws(%s)/%s", user, password, address, database)
        dialector = NewTDengineDialector(dsn)
    default:
        err = fmt.Errorf("invalid database driver %q", driver)
        return
    }

    log.Printf("try to connect to %q", dsn)
    db, err = gorm.Open(dialector, &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        err = fmt.Errorf("failed to connect to %q %v", dsn, err)
        return
    }

    // Configure connection pool to prevent connection leaks
    if sqlDB, sqlErr := db.DB(); sqlErr == nil {
        // Set maximum number of open connections
        sqlDB.SetMaxOpenConns(25)
        // Set maximum number of idle connections  
        sqlDB.SetMaxIdleConns(10)
        // Set maximum connection lifetime
        sqlDB.SetConnMaxLifetime(time.Hour)
        // Set maximum connection idle time
        sqlDB.SetConnMaxIdleTime(10 * time.Minute)
        
        log.Printf("Database connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v", 
            25, 10, time.Hour)
    }

    if driver != "tdengine" && driver != "greptime" {
        err = errors.Join(err, db.AutoMigrate(&TestSuite{}))
        err = errors.Join(err, db.AutoMigrate(&TestCase{}))
        err = errors.Join(err, db.AutoMigrate(&HistoryTestResult{}))
    }
    return
}
```

### 4. Add Connection Cleanup
```go
// Add cleanup function for graceful shutdown
func (s *dbserver) CleanupConnections() error {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    
    var errs []error
    for key, db := range dbCache {
        if sqlDB, err := db.DB(); err == nil {
            log.Printf("Closing database connection: %s", key)
            if closeErr := sqlDB.Close(); closeErr != nil {
                errs = append(errs, fmt.Errorf("failed to close connection %s: %w", key, closeErr))
            }
        }
        delete(dbCache, key)
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
```

### 5. Add Connection Monitoring
```go
// Add method to monitor connection stats
func (s *dbserver) GetConnectionStats() map[string]sql.DBStats {
    cacheMutex.RLock()
    defer cacheMutex.RUnlock()
    
    stats := make(map[string]sql.DBStats)
    for key, db := range dbCache {
        if sqlDB, err := db.DB(); err == nil {
            stats[key] = sqlDB.Stats()
        }
    }
    return stats
}
```

## Implementation Steps

1. **Backup Current Code**: Ensure you have a backup of the current `atest-ext-store-orm` code
2. **Apply Connection Pool Fix**: Update the `createDB` function with connection pool configuration
3. **Fix Cache Logic**: Implement the composite cache key strategy
4. **Add Cleanup**: Implement proper connection cleanup
5. **Add Monitoring**: Add connection statistics monitoring
6. **Test**: Thoroughly test database switching scenarios

## Testing the Fix

### Test Scenario
```go
// Test script to verify fix
func TestConnectionLeak(t *testing.T) {
    server := NewRemoteServer(10)
    
    // Create contexts for different databases
    ctx1 := remote.WithIncomingStoreContext(context.TODO(), &Store{
        Name: "test-store",
        URL: "localhost:3306",
        Username: "root", 
        Password: "root",
        Properties: map[string]string{
            "driver": "mysql",
            "database": "db1",
        },
    })
    
    ctx2 := remote.WithIncomingStoreContext(context.TODO(), &Store{
        Name: "test-store",
        URL: "localhost:3306", 
        Username: "root",
        Password: "root",
        Properties: map[string]string{
            "driver": "mysql",
            "database": "db2",
        },
    })
    
    // Alternate between databases multiple times
    for i := 0; i < 100; i++ {
        _, err1 := server.Query(ctx1, &server.DataQuery{Sql: "SELECT 1"})
        assert.NoError(t, err1)
        
        _, err2 := server.Query(ctx2, &server.DataQuery{Sql: "SELECT 1"}) 
        assert.NoError(t, err2)
        
        // Check connection stats
        if i%10 == 0 {
            stats := server.GetConnectionStats()
            for key, stat := range stats {
                t.Logf("Connection %s: Open=%d, InUse=%d", key, stat.OpenConnections, stat.InUse)
                // Ensure connections don't keep growing
                assert.True(t, stat.OpenConnections <= 50, "Too many open connections")
            }
        }
    }
}
```

## Benefits

1. **Prevents Connection Leaks**: Properly manages connection lifecycle
2. **Improved Performance**: Connection reuse reduces overhead
3. **Resource Control**: Connection pool limits prevent resource exhaustion
4. **Better Monitoring**: Connection statistics help diagnose issues
5. **Thread Safety**: Proper mutex usage prevents race conditions

## Configuration Options

Add these properties to store configuration for fine-tuning:

```yaml
properties:
  maxOpenConns: 25      # Maximum open connections
  maxIdleConns: 10      # Maximum idle connections  
  connMaxLifetime: 1h   # Maximum connection lifetime
  connMaxIdleTime: 10m  # Maximum connection idle time
```

This fix addresses the root cause of the database connection leak and provides a robust solution for connection management in the ORM store extension.
