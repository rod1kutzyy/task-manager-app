package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// ════════════════════════════════════════════════════════════════════════════════
//  CLI
// ════════════════════════════════════════════════════════════════════════════════

var (
	baseURL         = flag.String("base-url", "http://127.0.0.1:5050/api/v1", "Base API URL")
	numUsers        = flag.Int("users", 5, "Number of test users")
	tasksPerUser    = flag.Int("tasks-per-user", 20, "Tasks per user (seed)")
	concurrency     = flag.Int("concurrency", 50, "Goroutines per phase")
	phaseDuration   = flag.Duration("phase-duration", 20*time.Second, "Duration of each phase")
	readBurstCount  = flag.Int("read-burst", 10, "Consecutive GETs in a read-heavy iteration")
	mixedReadCount  = flag.Int("mixed-reads", 3, "GETs between mutations in mixed")
	mixedWriteCount = flag.Int("mixed-writes", 3, "Mutations (POST+PATCH+DELETE) in mixed")
	reportFile      = flag.String("report", "load_test_report.txt", "Path for the text report")
)

// ════════════════════════════════════════════════════════════════════════════════
//  DTO (swagger)
// ════════════════════════════════════════════════════════════════════════════════

type CreateUserRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
type UserResponse struct {
	ID uuid.UUID `json:"id"`
}

type CreateTaskRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	AuthorUserID uuid.UUID `json:"author_user_id"`
}
type TaskResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}
type PatchTaskRequest struct {
	Title     *string `json:"title,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

// ════════════════════════════════════════════════════════════════════════════════
//  HTTP
// ════════════════════════════════════════════════════════════════════════════════

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 500,
		IdleConnTimeout:     90 * time.Second,
	},
}

func doJSON(method, url string, payload any) (*http.Response, []byte, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, b, nil
}

// ════════════════════════════════════════════════════════════════════════════════
//  Metrics
// ════════════════════════════════════════════════════════════════════════════════

type Stats struct {
	mu        sync.Mutex
	name      string
	latencies []time.Duration
	ok, fail  int64
}

func NewStats(name string) *Stats { return &Stats{name: name} }

func (s *Stats) Record(d time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latencies = append(s.latencies, d)
	if success {
		s.ok++
	} else {
		s.fail++
	}
}

type Snap struct {
	Name                    string
	OK, Fail                int64
	Avg, P50, P95, P99, Max time.Duration
}

func (s *Stats) Snapshot() Snap {
	s.mu.Lock()
	defer s.mu.Unlock()
	sn := Snap{Name: s.name, OK: s.ok, Fail: s.fail}
	if len(s.latencies) == 0 {
		return sn
	}
	sorted := make([]time.Duration, len(s.latencies))
	copy(sorted, s.latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	var sum time.Duration
	for _, l := range sorted {
		sum += l
	}
	sn.Avg = sum / time.Duration(len(sorted))
	sn.P50 = pctl(sorted, 50)
	sn.P95 = pctl(sorted, 95)
	sn.P99 = pctl(sorted, 99)
	sn.Max = sorted[len(sorted)-1]
	return sn
}

func (s Snap) Total() int64 { return s.OK + s.Fail }
func (s Snap) ErrPct() float64 {
	if s.Total() == 0 {
		return 0
	}
	return float64(s.Fail) / float64(s.Total()) * 100
}

func pctl(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// ════════════════════════════════════════════════════════════════════════════════
//  Helpers
// ════════════════════════════════════════════════════════════════════════════════

func randomString(n int) string {
	const c = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = c[rand.Intn(len(c))]
	}
	return string(b)
}
func randomPhone() string     { return fmt.Sprintf("+7%010d", rand.Int63n(1e10)) }
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func fmtDur(d time.Duration) string {
	us := d.Microseconds()
	if us < 1000 {
		return fmt.Sprintf("%dµs", us)
	}
	return fmt.Sprintf("%.2fms", float64(us)/1000)
}

func fmtDurSigned(d time.Duration) string {
	if d >= 0 {
		return "+" + fmtDur(d)
	}
	return "-" + fmtDur(-d)
}

var stdinScanner = bufio.NewScanner(os.Stdin)

func waitForUser(prompt string) bool {
	fmt.Println()
	fmt.Printf("  ⏸  %s\n", prompt)
	fmt.Print("  ➜  Press Enter when ready (or 'q' to cancel): ")
	stdinScanner.Scan()
	text := strings.TrimSpace(strings.ToLower(stdinScanner.Text()))
	return text != "q" && text != "quit" && text != "n" && text != "no"
}

// ════════════════════════════════════════════════════════════════════════════════
//  Setup / Teardown
// ════════════════════════════════════════════════════════════════════════════════

func setupUsers(base string, n int) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, n)
	for i := 0; i < n; i++ {
		resp, body, err := doJSON(http.MethodPost, base+"/users", CreateUserRequest{
			FullName:    fmt.Sprintf("LT-User-%d-%s", i+1, randomString(4)),
			PhoneNumber: randomPhone(),
		})
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 201 {
			return nil, fmt.Errorf("POST /users → %d: %s", resp.StatusCode, body)
		}
		var u UserResponse
		json.Unmarshal(body, &u)
		ids = append(ids, u.ID)
	}
	return ids, nil
}

func seedTasks(base string, userIDs []uuid.UUID, perUser int) (map[uuid.UUID][]uuid.UUID, error) {
	m := make(map[uuid.UUID][]uuid.UUID, len(userIDs))
	for _, uid := range userIDs {
		for j := 0; j < perUser; j++ {
			resp, body, err := doJSON(http.MethodPost, base+"/tasks", CreateTaskRequest{
				Title:        fmt.Sprintf("Seed task %s-%d", uid, j),
				Description:  fmt.Sprintf("Seeded task for user %s", uid),
				AuthorUserID: uid,
			})
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != 201 {
				return nil, fmt.Errorf("seed POST /tasks → %d: %s", resp.StatusCode, body)
			}
			var t TaskResponse
			json.Unmarshal(body, &t)
			m[uid] = append(m[uid], t.ID)
		}
	}
	return m, nil
}

func teardownTasks(base string, m map[uuid.UUID][]uuid.UUID) {
	for _, tasks := range m {
		for _, tid := range tasks {
			doJSON(http.MethodDelete, fmt.Sprintf("%s/tasks/%s", base, tid), nil)
		}
	}
}

func teardownUsers(base string, ids []uuid.UUID) {
	for _, id := range ids {
		doJSON(http.MethodDelete, fmt.Sprintf("%s/users/%s", base, id), nil)
	}
}

func warmUpCache(base string, userIDs []uuid.UUID) {
	for _, uid := range userIDs {
		doJSON(http.MethodGet, fmt.Sprintf("%s/tasks?user_id=%s&limit=10", base, uid), nil)
	}
}

func checkServer(base string) error {
	fmt.Print("  Checking server... ")
	_, _, err := doJSON(http.MethodGet, base+"/tasks?limit=1", nil)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return err
	}
	fmt.Println("OK")
	return nil
}

// ════════════════════════════════════════════════════════════════════════════════
//  Workers
// ════════════════════════════════════════════════════════════════════════════════

func readHeavyWorker(
	stop <-chan struct{}, wg *sync.WaitGroup,
	base string, myUID uuid.UUID, allUIDs []uuid.UUID, tasksByUser map[uuid.UUID][]uuid.UUID,
	sListSelf, sListOther, sGetByID *Stats, iters *atomic.Int64,
) {
	defer wg.Done()
	myTasks := tasksByUser[myUID]
	for {
		select {
		case <-stop:
			return
		default:
		}
		otherUID := pickOther(myUID, allUIDs)
		otherTasks := tasksByUser[otherUID]
		for i := 0; i < *readBurstCount; i++ {
			select {
			case <-stop:
				return
			default:
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%s&limit=10&offset=%d",
					base, myUID, rand.Intn(maxInt(1, len(myTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListSelf.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.Intn(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, fmt.Sprintf("%s/tasks/%s", base, tid), nil)
				sGetByID.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%s&limit=10&offset=%d",
					base, otherUID, rand.Intn(maxInt(1, len(otherTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListOther.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
		}
		iters.Add(1)
	}
}

func mixedWorker(
	stop <-chan struct{}, wg *sync.WaitGroup,
	base string, myUID uuid.UUID, allUIDs []uuid.UUID, tasksByUser map[uuid.UUID][]uuid.UUID,
	sListSelf, sListOther, sGetByID, sPost, sPatch, sDelete *Stats,
	iters *atomic.Int64,
) {
	defer wg.Done()
	myTasks := make([]uuid.UUID, len(tasksByUser[myUID]))
	copy(myTasks, tasksByUser[myUID])
	for {
		select {
		case <-stop:
			return
		default:
		}
		otherUID := pickOther(myUID, allUIDs)
		otherTasks := tasksByUser[otherUID]

		for i := 0; i < *mixedReadCount; i++ {
			select {
			case <-stop:
				return
			default:
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%s&limit=10&offset=%d",
					base, myUID, rand.Intn(maxInt(1, len(myTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListSelf.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.Intn(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, fmt.Sprintf("%s/tasks/%s", base, tid), nil)
				sGetByID.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
			{
				url := fmt.Sprintf("%s/tasks?user_id=%s&limit=10&offset=%d",
					base, otherUID, rand.Intn(maxInt(1, len(otherTasks)-10)))
				start := time.Now()
				resp, _, err := doJSON(http.MethodGet, url, nil)
				sListOther.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
		}

		for i := 0; i < *mixedWriteCount; i++ {
			select {
			case <-stop:
				return
			default:
			}
			{
				start := time.Now()
				resp, body, err := doJSON(http.MethodPost, base+"/tasks", CreateTaskRequest{
					Title:        fmt.Sprintf("Mixed %s [u:%s]", randomString(6), myUID),
					Description:  fmt.Sprintf("Mixed load task user %s", myUID),
					AuthorUserID: myUID,
				})
				sPost.Record(time.Since(start), err == nil && resp.StatusCode == 201)
				if err == nil && resp.StatusCode == 201 {
					var t TaskResponse
					json.Unmarshal(body, &t)
					if t.ID != uuid.Nil {
						myTasks = append(myTasks, t.ID)
					}
				}
			}
			if len(myTasks) > 0 {
				tid := myTasks[rand.Intn(len(myTasks))]
				start := time.Now()
				resp, _, err := doJSON(http.MethodPatch, fmt.Sprintf("%s/tasks/%s", base, tid), PatchTaskRequest{
					Title: strPtr(fmt.Sprintf("Upd %s", randomString(4))), Completed: boolPtr(rand.Intn(2) == 1),
				})
				sPatch.Record(time.Since(start), err == nil && resp.StatusCode == 200)
			}
			if len(myTasks) > 1 {
				idx := rand.Intn(len(myTasks))
				tid := myTasks[idx]
				myTasks[idx] = myTasks[len(myTasks)-1]
				myTasks = myTasks[:len(myTasks)-1]
				start := time.Now()
				resp, _, err := doJSON(http.MethodDelete, fmt.Sprintf("%s/tasks/%s", base, tid), nil)
				sDelete.Record(time.Since(start), err == nil && resp.StatusCode == 204)
			}
		}
		iters.Add(1)
	}
}

func pickOther(my uuid.UUID, all []uuid.UUID) uuid.UUID {
	if len(all) <= 1 {
		return my
	}
	o := my
	for o == my {
		o = all[rand.Intn(len(all))]
	}
	return o
}

// ════════════════════════════════════════════════════════════════════════════════
//  Progress (stdout only)
// ════════════════════════════════════════════════════════════════════════════════

func showProgress(stop <-chan struct{}, phase string, dur time.Duration, iters *atomic.Int64) {
	start := time.Now()
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-stop:
			fmt.Println()
			return
		case <-tick.C:
			el := time.Since(start).Round(time.Second)
			pct := float64(el) / float64(dur) * 100
			if pct > 100 {
				pct = 100
			}
			fmt.Printf("\r    [%s] %s / %s (%4.0f%%)  iterations: %d  ",
				phase, el, dur, pct, iters.Load())
		}
	}
}

// ════════════════════════════════════════════════════════════════════════════════
//  Phase run
// ════════════════════════════════════════════════════════════════════════════════

func runPhase(
	name string, dur time.Duration, conc int,
	base string, userIDs []uuid.UUID, tasksByUser map[uuid.UUID][]uuid.UUID,
	launch func(stop <-chan struct{}, wg *sync.WaitGroup, uid uuid.UUID, iters *atomic.Int64),
) {
	fmt.Printf("\n  ▶ Phase: %s (%s, %d goroutines)\n", name, dur, conc)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var iters atomic.Int64
	pstop := make(chan struct{})
	go showProgress(pstop, name, dur, &iters)
	for i := 0; i < conc; i++ {
		uid := userIDs[i%len(userIDs)]
		wg.Add(1)
		go func(u uuid.UUID) { launch(stop, &wg, u, &iters) }(uid)
	}
	time.Sleep(dur)
	close(stop)
	wg.Wait()
	close(pstop)
	fmt.Printf("    Completed: %d iterations\n", iters.Load())
}

// ════════════════════════════════════════════════════════════════════════════════
//  runFullTest
// ════════════════════════════════════════════════════════════════════════════════

type TestResult struct {
	ReadSnaps       []Snap
	MixedReadSnaps  []Snap
	MixedWriteSnaps []Snap
}

func runFullTest(base string, userIDs []uuid.UUID, tasksByUser map[uuid.UUID][]uuid.UUID) TestResult {
	warmUpCache(base, userIDs)

	rSelf := NewStats("GET /tasks?user_id=self")
	rOther := NewStats("GET /tasks?user_id=other")
	rByID := NewStats("GET /tasks/{id}")
	runPhase("READ-HEAVY", *phaseDuration, *concurrency, base, userIDs, tasksByUser,
		func(stop <-chan struct{}, wg *sync.WaitGroup, uid uuid.UUID, iters *atomic.Int64) {
			readHeavyWorker(stop, wg, base, uid, userIDs, tasksByUser, rSelf, rOther, rByID, iters)
		})

	warmUpCache(base, userIDs)

	mSelf := NewStats("GET /tasks?user_id=self")
	mOther := NewStats("GET /tasks?user_id=other")
	mByID := NewStats("GET /tasks/{id}")
	mPost := NewStats("POST /tasks")
	mPatch := NewStats("PATCH /tasks/{id}")
	mDel := NewStats("DELETE /tasks/{id}")
	runPhase("MIXED", *phaseDuration, *concurrency, base, userIDs, tasksByUser,
		func(stop <-chan struct{}, wg *sync.WaitGroup, uid uuid.UUID, iters *atomic.Int64) {
			mixedWorker(stop, wg, base, uid, userIDs, tasksByUser,
				mSelf, mOther, mByID, mPost, mPatch, mDel, iters)
		})

	return TestResult{
		ReadSnaps:       []Snap{rSelf.Snapshot(), rOther.Snapshot(), rByID.Snapshot()},
		MixedReadSnaps:  []Snap{mSelf.Snapshot(), mOther.Snapshot(), mByID.Snapshot()},
		MixedWriteSnaps: []Snap{mPost.Snapshot(), mPatch.Snapshot(), mDel.Snapshot()},
	}
}

// ════════════════════════════════════════════════════════════════════════════════
//  Table formatting (accept io.Writer)
// ════════════════════════════════════════════════════════════════════════════════

func fprintSnapTable(w io.Writer, title string, rows []Snap) {
	fmt.Fprintf(w, "\n  ┌─ %s\n", title)
	fmt.Fprintf(w, "  │ %-32s │ %7s │ %6s │ %5s │ %9s │ %9s │ %9s │ %9s\n",
		"Endpoint", "OK", "FAIL", "Err%", "Avg", "p50", "p95", "p99")
	fmt.Fprintln(w, "  │ "+
		"─────────────────────────────────┼─────────┼────────┼───────┼───────────┼───────────┼───────────┼───────────")
	for _, r := range rows {
		if r.Total() == 0 {
			fmt.Fprintf(w, "  │ %-32s │     —   │   —    │  —    │     —     │     —     │     —     │     —\n", r.Name)
			continue
		}
		fmt.Fprintf(w, "  │ %-32s │ %7d │ %6d │ %4.1f%% │ %9s │ %9s │ %9s │ %9s\n",
			r.Name, r.OK, r.Fail, r.ErrPct(),
			fmtDur(r.Avg), fmtDur(r.P50), fmtDur(r.P95), fmtDur(r.P99))
	}
	fmt.Fprintln(w, "  └──")
}

type compareSection struct {
	title string
	a, b  []Snap
}

func fprintCrossCompare(w io.Writer, labelA, labelB string, sections []compareSection) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  ╔════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "  ║  COMPARISON:  A = %-20s  vs  B = %-20s                            ║\n", labelA, labelB)
	fmt.Fprintln(w, "  ╠════════════════════════════════════════════════════════════════════════════════════════════════════════╣")

	for _, sec := range sections {
		fmt.Fprintf(w, "\n  ┌─ %s\n", sec.title)
		fmt.Fprintf(w, "  │ %-32s │ %10s │ %10s │ %10s │ %10s │ %10s │ %10s │ %8s\n",
			"Endpoint", labelA+" p50", labelB+" p50", "Δ p50",
			labelA+" p95", labelB+" p95", "Δ p95", "Speedup")
		fmt.Fprintln(w, "  │ "+
			"─────────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼──────────")

		n := len(sec.a)
		if len(sec.b) < n {
			n = len(sec.b)
		}
		for i := 0; i < n; i++ {
			a, b := sec.a[i], sec.b[i]
			dp50 := b.P50 - a.P50
			dp95 := b.P95 - a.P95
			sp := float64(0)
			if b.P50 > 0 {
				sp = float64(a.P50) / float64(b.P50)
			}
			fmt.Fprintf(w, "  │ %-32s │ %10s │ %10s │ %10s │ %10s │ %10s │ %10s │ %7.2fx\n",
				a.Name,
				fmtDur(a.P50), fmtDur(b.P50), fmtDurSigned(dp50),
				fmtDur(a.P95), fmtDur(b.P95), fmtDurSigned(dp95), sp)
		}
		fmt.Fprintln(w, "  └──")
	}
}

func fprintHeader(w io.Writer) {
	fmt.Fprintln(w, "╔═══════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(w, "║    Load test /tasks - comparison: no cache vs with cache          ║")
	fmt.Fprintln(w, "╚═══════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Date:             %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "  Base URL:         %s\n", *baseURL)
	fmt.Fprintf(w, "  Users:            %d\n", *numUsers)
	fmt.Fprintf(w, "  Tasks/user seed:  %d\n", *tasksPerUser)
	fmt.Fprintf(w, "  Concurrency:      %d\n", *concurrency)
	fmt.Fprintf(w, "  Phase duration:   %s\n", *phaseDuration)
	fmt.Fprintf(w, "  Read burst:       %d\n", *readBurstCount)
	fmt.Fprintf(w, "  Mixed r/w:        %d reads + %d writes\n", *mixedReadCount, *mixedWriteCount)
	fmt.Fprintln(w)
}

func main() {
	flag.Parse()

	// Open the report file immediately - write to it as we go
	if dir := filepath.Dir(*reportFile); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("  Error creating directory %s: %v\n", dir, err)
			return
		}
	}
	f, err := os.Create(*reportFile)
	if err != nil {
		fmt.Printf("  Error creating report file %s: %v\n", *reportFile, err)
		return
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()

	// out writes to stdout and file simultaneously
	out := io.MultiWriter(os.Stdout, bw)

	fprintHeader(out)

	fmt.Fprintf(out, "  Report:           %s\n", *reportFile)
	fmt.Fprintln(out)

	// ══════════════════════════════════════════════════════════════════════
	//  STEP 1: server WITHOUT cache
	// ══════════════════════════════════════════════════════════════════════

	if !waitForUser("Start the server WITHOUT cache and press Enter...") {
		fmt.Println("  Cancelled.")
		return
	}

	if err := checkServer(*baseURL); err != nil {
		return
	}

	// Setup
	fmt.Fprintln(out, "\n  ── SETUP ──")
	fmt.Fprintf(out, "  Creating %d users... ", *numUsers)
	userIDs, err := setupUsers(*baseURL, *numUsers)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	fmt.Fprintf(out, "OK (IDs: %v)\n", userIDs)

	fmt.Fprintf(out, "  Creating %d tasks per user... ", *tasksPerUser)
	tasksByUser, err := seedTasks(*baseURL, userIDs, *tasksPerUser)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		teardownUsers(*baseURL, userIDs)
		return
	}
	total := 0
	for _, t := range tasksByUser {
		total += len(t)
	}
	fmt.Fprintf(out, "OK (%d tasks)\n", total)

	// Run A
	fmt.Fprintln(out, "\n  ════════════════════════════════════════════")
	fmt.Fprintln(out, "    RUN A: WITHOUT CACHE")
	fmt.Fprintln(out, "  ════════════════════════════════════════════")

	resultA := runFullTest(*baseURL, userIDs, tasksByUser)

	fprintSnapTable(out, "A / READ-HEAVY (read only)", resultA.ReadSnaps)
	fprintSnapTable(out, "A / MIXED - reads", resultA.MixedReadSnaps)
	fprintSnapTable(out, "A / MIXED - mutations", resultA.MixedWriteSnaps)

	// Flush so the file has everything before pause
	bw.Flush()

	// ══════════════════════════════════════════════════════════════════════
	//  STEP 2: restart server WITH cache
	// ══════════════════════════════════════════════════════════════════════

	fmt.Fprintln(out)
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")
	fmt.Fprintln(out, "  Run A completed. The script will recreate data for Run B to ensure a fair comparison.")
	fmt.Fprintln(out, "  ────────────────────────────────────────────────────────────")

	if !waitForUser("Stop the server, enable cache, restart it, and press Enter...") {
		fmt.Fprintln(out, "\n  Cancelled.")
		teardownTasks(*baseURL, tasksByUser)
		teardownUsers(*baseURL, userIDs)
		return
	}

	if err := checkServer(*baseURL); err != nil {
		return
	}

	// Recreate data for a fair A/B comparison.
	fmt.Fprintln(out, "\n  ── RESET DATA BEFORE RUN B ──")
	fmt.Fprintln(out, "  Deleting Run A tasks/users...")
	teardownTasks(*baseURL, tasksByUser)
	teardownUsers(*baseURL, userIDs)
	fmt.Fprintln(out, "  OK")

	fmt.Fprintln(out, "\n  ── SETUP (RUN B) ──")
	fmt.Fprintf(out, "  Creating %d users... ", *numUsers)
	userIDsB, err := setupUsers(*baseURL, *numUsers)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		return
	}
	fmt.Fprintf(out, "OK (IDs: %v)\n", userIDsB)

	fmt.Fprintf(out, "  Creating %d tasks per user... ", *tasksPerUser)
	tasksByUserB, err := seedTasks(*baseURL, userIDsB, *tasksPerUser)
	if err != nil {
		fmt.Fprintf(out, "ERROR: %v\n", err)
		teardownUsers(*baseURL, userIDsB)
		return
	}
	totalB := 0
	for _, t := range tasksByUserB {
		totalB += len(t)
	}
	fmt.Fprintf(out, "OK (%d tasks)\n", totalB)

	// Run B
	fmt.Fprintln(out, "\n  ════════════════════════════════════════════")
	fmt.Fprintln(out, "    RUN B: WITH CACHE")
	fmt.Fprintln(out, "  ════════════════════════════════════════════")

	resultB := runFullTest(*baseURL, userIDsB, tasksByUserB)

	fprintSnapTable(out, "B / READ-HEAVY (read only)", resultB.ReadSnaps)
	fprintSnapTable(out, "B / MIXED - reads", resultB.MixedReadSnaps)
	fprintSnapTable(out, "B / MIXED - mutations", resultB.MixedWriteSnaps)

	// ══════════════════════════════════════════════════════════════════════
	//  Comparison A vs B
	// ══════════════════════════════════════════════════════════════════════

	fprintCrossCompare(out, "no-cache", "with-cache", []compareSection{
		{"READ-HEAVY (read only)", resultA.ReadSnaps, resultB.ReadSnaps},
		{"MIXED - reads (cache invalidated)", resultA.MixedReadSnaps, resultB.MixedReadSnaps},
		{"MIXED - mutations", resultA.MixedWriteSnaps, resultB.MixedWriteSnaps},
	})

	// Flush file before teardown
	bw.Flush()

	fmt.Printf("\n  Report saved: %s\n", *reportFile)

	// ══════════════════════════════════════════════════════════════════════
	//  Teardown
	// ══════════════════════════════════════════════════════════════════════

	fmt.Println("\n  ── TEARDOWN ──")
	fmt.Print("  Deleting tasks... ")
	teardownTasks(*baseURL, tasksByUserB)
	fmt.Println("OK")
	fmt.Print("  Deleting users... ")
	teardownUsers(*baseURL, userIDsB)
	fmt.Println("OK")

	fmt.Println("\n  Done.")
}
