package main

import (
        "database/sql"
        "fmt"
)

import (
        "bufio"
        "flag"
        _ "github.com/go-sql-driver/mysql"
        "log"
        "net"
        "net/http"
        "os"
        "os/user"
        "regexp"
        "strconv"
        "sync"
        "time"
)

const WSREPLOCALSTATE = "wsrep_local_state"
const DEFAULTPORT = 4793
const DEFAULT_POLL_INTERVAL = 60
const DB_USER = "root"
const DB_PASSWD = "theD@t@isM1n3"
const DB_HOST = "localhost"
const DB_MYCNF = ".my.cnf"
const HOST_REGEX = `^host\s+=\s+([a-zA-Z0-9-_.]+)`
const USER_REGEX = `^user\s+=\s+([a-zA-Z0-9-_]+)`
const PASSWORD_REGEX = `^password\s+=\s+(.+)$`

type Metric struct {
        Name      string
        Value     string
        Timestamp int64
}

type MetricStore struct {
        store map[string]Metric
        Mutex *sync.Mutex
}

func (m *MetricStore) init() {
        m.store = make(map[string]Metric)
        m.Mutex = &sync.Mutex{}
}

func (m *MetricStore) get(s string) Metric {
        m.Mutex.Lock()
        v, ok := m.store[s]
        m.Mutex.Unlock()

        if ok == false {
                return Metric{"", "", 0}
        }

        return v
}

func (m *MetricStore) put(metric Metric) {

        m.Mutex.Lock()
        m.store[metric.Name] = metric
        m.Mutex.Unlock()
}

type dbConnectionParameters struct {
        user     string
        password string
        host     string
}

func ReadMyCnfFileforPasswd() dbConnectionParameters {
        user_regex, err := regexp.Compile(USER_REGEX)

        if err != nil {
                fmt.Println("Can not compile USER_REGEX regular expression.")
        }
        passwd_regex, err := regexp.Compile(PASSWORD_REGEX)
        if err != nil {
                fmt.Println("Can not compile PASSWORD_REGEX regular expression.")
        }
        host_regex, err := regexp.Compile(HOST_REGEX)
        if err != nil {
                fmt.Println("Can not compile HOST_REGEX regular expression.")
        }

        _user := getUser()

        file, err := os.Open(_user.HomeDir + "/" + DB_MYCNF)

        if err != nil {
                return dbConnectionParameters{}
        }

        defer file.Close()

        connParameters := dbConnectionParameters{}

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                line := scanner.Text()

                if len(line) > 0 {
                        user_match := user_regex.FindSubmatch([]byte(line))

                        if user_match != nil {
                                connParameters.user = string(user_match[1])
                        }

                        password_match := passwd_regex.FindSubmatch([]byte(line))

                        if password_match != nil {
                                connParameters.password = string(password_match[1])
                        }

                        host_match := host_regex.FindSubmatch([]byte(line))

                        if host_match != nil {
                                connParameters.host = string(host_match[1])
                        }
                }
        }

        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }

        return connParameters
}

func getUser() user.User {
        var _user *user.User
        _user, err := user.Current()

        if err != nil {
                return user.User{}
        }

        return *_user
}

func assembleConnectionString(user string, passwd string, host string) string {
        return user + ":" + passwd + "@(" + host + ")/mysql"
}

func get_wsrep_local_state(metricsCh chan Metric) {
        connParam := ReadMyCnfFileforPasswd()
        fmt.Printf("%v\n", connParam)
        db, err := sql.Open("mysql", assembleConnectionString(
                connParam.user, connParam.password, connParam.host))

        if err != nil {
                fmt.Printf("DB Error: %s", err)
        }
        defer db.Close()

        rows, err := db.Query("show status")
        if err != nil {
                fmt.Println("Query Error")
        }

        defer rows.Close()

        for rows.Next() {
                var name string
                var value string
                if err := rows.Scan(&name, &value); err != nil {
                        fmt.Printf("Scan Error: %s", err)
                }

                now := time.Now()
                m := Metric{name, value, now.Unix()}
                metricsCh <- m
        }
}

func QueryFullProcessList(metricChan chan Metric) {
        connParam := ReadMyCnfFileforPasswd()
        db, err := sql.Open("mysql", assembleConnectionString(
                connParam.user, connParam.password, connParam.host))

        if err != nil {
                fmt.Printf("DB Error: %s", err)
        }
        defer db.Close()

        rows, err := db.Query("show full processlist")
        if err != nil {
                fmt.Println("Query Error")
        }

        defer rows.Close()

        userDbCommand := UserDbCommand{}
        userDbCommand.init()

        for rows.Next() {
                var Id sql.NullString
                var User sql.NullString
                var Host sql.NullString
                var db sql.NullString
                var Command sql.NullString
                var Time sql.NullString
                var State sql.NullString
                var Info sql.NullString
                var Progress sql.NullString
                if err := rows.Scan(&Id, &User, &Host, &db, &Command, &Time, &State, &Info, &Progress); err != nil {
                        fmt.Printf("Scan Error: %s", err)
                }

                process := Process{Id: Id.String,
                        User: User.String,
                        Host: Host.String,
                        db: db.String,
                        Command: Command.String,
                        Time: Time.String,
                        State: State.String,
                        Info: Info.String,
                        Progress: Progress.String}

                userDbCommand.generate(process)
        }

        userDbCommand.dump(metricChan)
}

func IsNumeric(s string) bool {
        _, err := strconv.ParseFloat(s, 64)
        return err == nil
}

func send(m Metric) {
        hostname, err := os.Hostname()

        if err != nil {
                fmt.Println("Error getting hostname.")
                return
        }

        host_addr := "graphite:2003"

        conn, err := net.DialTimeout("tcp", host_addr, time.Duration(time.Second*30))

        defer conn.Close()

        if err != nil {
                fmt.Println("Error with DialTCP")
                return
        }

        metric := fmt.Sprintf("metrics.mariadb.%s.%s %s %v\n", hostname, m.Name, m.Value, m.Timestamp)
        fmt.Printf("metrics.mariadb.%s.%s.%s %v\n", hostname, m.Name, m.Value, m.Timestamp)

        conn.Write([]byte(metric))
}

func process_metric(metricsCh chan Metric, metricStore MetricStore) {
        for {
                m := <-metricsCh
                //fmt.Printf("%s  %s %d\n", m.Name, m.Value, m.Timestamp)
                metricStore.put(m)

                if IsNumeric(m.Value) {
                        send(m)
                }
        }
}

func setError(w http.ResponseWriter) {
        http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func status(metricStore MetricStore) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                s := metricStore.get(WSREPLOCALSTATE)

                state, err := strconv.Atoi(s.Value)

                if err != nil {
                        setError(w)
                        return
                }

                if state != 4 {
                        setError(w)
                        return
                }

                now := time.Now()
                epoch := now.Unix()
                delta := epoch - s.Timestamp
                if delta > 180 {
                        setError(w)
                        return
                }

                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                fmt.Fprint(w, "OK")

        })
}

func get_metric(metricStore MetricStore) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                params, ok := r.URL.Query()["m"]

                if !ok || len(params[0]) < 1 {
                        fmt.Fprintf(w, "[]")
                        return
                }

                m := metricStore.get(params[0])

                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                fmt.Fprintf(w, "%s %s", params[0], m.Value)
        })
}

func put_metric(metricStore MetricStore) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                params_k, ok := r.URL.Query()["k"]

                if !ok || len(params_k[0]) < 1 {
                        fmt.Fprintf(w, "[]")
                        return
                }

                params_v, ok := r.URL.Query()["v"]

                if !ok || len(params_v[0]) < 1 {
                        fmt.Fprintf(w, "[]")
                        return
                }

                now := time.Now()

                metric := Metric{params_k[0], params_v[0], now.Unix()}

                metricStore.put(metric)

                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                fmt.Fprintf(w, "%s %s", metric.Name, metric.Value)
        })
}

func registerHandlers(metricStore MetricStore) {
        http.HandleFunc("/", status(metricStore))
        http.HandleFunc("/get", get_metric(metricStore))
        http.HandleFunc("/put", put_metric(metricStore))
}

func webserver(port int) {

        portString := strconv.Itoa(port)
        for {
                http.ListenAndServe(":"+portString, nil)
        }
}

func main() {

        pollIntervalPtr := flag.Int("i", DEFAULT_POLL_INTERVAL, "Connection polling interval")
        emulatePtr := flag.Bool("e", false, "Emulation flag")
        portPtr := flag.Int("p", DEFAULTPORT, "Connection polling interval")
        flag.Parse()

        var metricStore MetricStore
        metricStore.init()

        metricsCh := make(chan Metric)
        go process_metric(metricsCh, metricStore)

        go QueryFullProcessList(metricsCh)

        if !*emulatePtr {
                go get_wsrep_local_state(metricsCh)
        }
        registerHandlers(metricStore)
        go webserver(*portPtr)

        var tickerInterval = time.Second * time.Duration(*pollIntervalPtr)

        ticker := time.NewTicker(tickerInterval)
        for range ticker.C {
                go QueryFullProcessList(metricsCh)
                if !*emulatePtr {
                        go get_wsrep_local_state(metricsCh)
                }
        }
}
