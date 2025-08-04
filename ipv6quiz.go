package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type QuizResult struct {
	Score      int
	Total      int
	Percentage float64
}

var (
	questions      []Question
	quizTemplate   *template.Template
	currentQuiz    []Question
	userAnswers    []int
	correctAnswers []int
)

func main() {
	daemon := flag.Bool("d", false, "Run as a daemon")
	flag.Parse()

	if *daemon {
		daemonize()
	}

	initializeQuestions()
	initializeTemplates()
	startWebServer()
}

// Add bits and pieces to daemonize this using -d

func daemonize() {
	// Check if we're already a daemon
	if os.Getppid() == 1 {
		return
	}

	// Create a new command with the same arguments
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	// Start the process
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to daemonize: %v", err)
	}

	// Write PID file
	pid := cmd.Process.Pid
	err = os.WriteFile("ipv6quiz.pid", []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		log.Printf("Warning: Could not write PID file: %v", err)
	}

	fmt.Printf("Daemon started with PID %d\n", pid)
	os.Exit(0)
}

func initializeQuestions() {
	questions = []Question{
		// IPv6 quiz questions, add more in the correct format and recompile to update.
		{
			Text: "Which RFC defines the newest address range was originally reserved for documentation?",
			Options: []string{
				"RFC 6724",
				"RFC 8200",
				"RFC 9637",
				"RFC 2549",
			},
			Answer: 2,
		},
		{
			Text: "Which IPv6 address range was originally reserved for documentation?",
			Options: []string{
				"2001:db8::/32",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 equivalent of IPv4's 169.254.0.0/16 (APIPA)?",
			Options: []string{
				"::1/128",
				"fe80::/10",
				"ff02::1",
				"2000::/3",
			},
			Answer: 1,
		},
		{
			Text: "How many hexadecimal digits are in an uncompressed IPv6 address?",
			Options: []string{
				"16",
				"24",
				"32",
				"48",
			},
			Answer: 2,
		},
		{
			Text: "Which multicast address represents all IPv6 routers?",
			Options: []string{
				"ff02::1",
				"ff02::2",
				"ff05::1",
				"ff05::2",
			},
			Answer: 1,
		},
		{
			Text: "What is the maximum number of consecutive zero groups that can be compressed with '::'?",
			Options: []string{
				"1",
				"2",
				"Any number",
				"None, only single zeros can be compressed",
			},
			Answer: 2,
		},
		{
			Text: "Which ICMPv6 message type is used for Neighbor Solicitation?",
			Options: []string{
				"Type 133",
				"Type 134",
				"Type 135",
				"Type 136",
			},
			Answer: 2,
		},
		{
			Text: "What is the purpose of the Hop Limit field in the IPv6 header?",
			Options: []string{
				"To identify the packet flow",
				"To prevent infinite packet looping",
				"To specify the next header type",
				"To indicate packet priority",
			},
			Answer: 1,
		},
		{
			Text: "Which IPv6 extension header must be processed by every router along the path?",
			Options: []string{
				"Hop-by-Hop Options",
				"Routing",
				"Fragment",
				"Destination Options",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for global unicast addresses?",
			Options: []string{
				"2000::/3",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 0,
		},
		{
			Text: "Which protocol number is used for ICMPv6?",
			Options: []string{
				"1",
				"6",
				"17",
				"58",
			},
			Answer: 3,
		},
		{
			Text: "What is the IPv6 site-local address range (deprecated)?",
			Options: []string{
				"fec0::/10",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 0,
		},
		{
			Text: "Which DNS record type is used for IPv6 addresses?",
			Options: []string{
				"A",
				"AAAA",
				"PTR",
				"MX",
			},
			Answer: 1,
		},
		{
			Text: "What is the default Teredo server IPv4 address?",
			Options: []string{
				"65.54.227.120",
				"74.125.43.99",
				"192.88.99.1",
				"208.67.222.222",
			},
			Answer: 0,
		},
		{
			Text: "Which IPv6 transition technology uses IPv4-mapped addresses?",
			Options: []string{
				"6to4",
				"Teredo",
				"ISATAP",
				"SIIT",
			},
			Answer: 3,
		},
		{
			Text: "What is the IPv6 prefix for 6to4 addresses?",
			Options: []string{
				"2001::/32",
				"2002::/16",
				"fc00::/7",
				"ff00::/8",
			},
			Answer: 1,
		},
		{
			Text: "Which multicast scope represents link-local scope?",
			Options: []string{
				"1",
				"2",
				"5",
				"8",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for Unique Local Addresses (ULA)?",
			Options: []string{
				"2000::/3",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 1,
		},
		{
			Text: "Which field in the IPv6 header replaces the IPv4 TTL field?",
			Options: []string{
				"Version",
				"Traffic Class",
				"Flow Label",
				"Hop Limit",
			},
			Answer: 3,
		},
		{
			Text: "What is the IPv6 prefix for Teredo addresses?",
			Options: []string{
				"2001::/32",
				"2002::/16",
				"fc00::/7",
				"ff00::/8",
			},
			Answer: 0,
		},
		{
			Text: "Which ICMPv6 message type is used for Router Advertisement?",
			Options: []string{
				"Type 133",
				"Type 134",
				"Type 135",
				"Type 136",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for ORCHID addresses?",
			Options: []string{
				"2001:10::/28",
				"2001:20::/28",
				"2001:db8::/32",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which IPv6 extension header is used for Mobile IPv6?",
			Options: []string{
				"Hop-by-Hop Options",
				"Routing",
				"Fragment",
				"Destination Options",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for discard-only addresses?",
			Options: []string{
				"100::/64",
				"2001::/32",
				"2002::/16",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all DHCPv6 servers?",
			Options: []string{
				"ff02::1",
				"ff02::2",
				"ff02::1:2",
				"ff05::1:3",
			},
			Answer: 2,
		},
		{
			Text: "What is the IPv6 prefix for 6bone testing addresses?",
			Options: []string{
				"3ffe::/16",
				"2001::/32",
				"2002::/16",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which ICMPv6 message type is used for Neighbor Advertisement?",
			Options: []string{
				"Type 133",
				"Type 134",
				"Type 135",
				"Type 136",
			},
			Answer: 3,
		},
		{
			Text: "What is the IPv6 prefix for benchmark testing addresses?",
			Options: []string{
				"2001:2::/48",
				"2001:db8::/32",
				"3ffe::/16",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all NTP servers?",
			Options: []string{
				"ff02::101",
				"ff02::102",
				"ff05::101",
				"ff05::102",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for automatic tunnel addresses?",
			Options: []string{
				"::/96",
				"::ffff:0:0/96",
				"2001::/32",
				"2002::/16",
			},
			Answer: 0,
		},
		{
			Text: "Which IPv6 extension header is used for ESP?",
			Options: []string{
				"50",
				"51",
				"58",
				"59",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for IPv4-compatible addresses (deprecated)?",
			Options: []string{
				"::/96",
				"::ffff:0:0/96",
				"2001::/32",
				"2002::/16",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all MLDv2-capable routers?",
			Options: []string{
				"ff02::16",
				"ff02::17",
				"ff05::16",
				"ff05::17",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for IPv4-mapped addresses?",
			Options: []string{
				"::/96",
				"::ffff:0:0/96",
				"2001::/32",
				"2002::/16",
			},
			Answer: 1,
		},
		{
			Text: "Which IPv6 extension header is used for AH?",
			Options: []string{
				"50",
				"51",
				"58",
				"59",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for ORCHIDv2 addresses?",
			Options: []string{
				"2001:20::/28",
				"2001:30::/28",
				"2001:db8::/32",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all site-local routers?",
			Options: []string{
				"ff02::1",
				"ff02::2",
				"ff05::1",
				"ff05::2",
			},
			Answer: 3,
		},
		{
			Text: "What is the IPv6 prefix for 6rd addresses?",
			Options: []string{
				"2001::/32",
				"2002::/16",
				"fc00::/7",
				"ff00::/8",
			},
			Answer: 0,
		},
		{
			Text: "Which ICMPv6 message type is used for Redirect?",
			Options: []string{
				"Type 137",
				"Type 138",
				"Type 139",
				"Type 140",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for LISP addresses?",
			Options: []string{
				"2001:5::/32",
				"2001:10::/28",
				"2001:20::/28",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all PIM routers?",
			Options: []string{
				"ff02::d",
				"ff02::13",
				"ff05::d",
				"ff05::13",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for AS112 addresses?",
			Options: []string{
				"2001:4:112::/48",
				"2001:10::/28",
				"2001:20::/28",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which IPv6 extension header is used for mobility?",
			Options: []string{
				"43",
				"44",
				"50",
				"51",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for DNS name server addresses?",
			Options: []string{
				"2001:500::/32",
				"2001:501::/32",
				"2001:502::/32",
				"2001:503::/32",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all EIGRP routers?",
			Options: []string{
				"ff02::a",
				"ff02::b",
				"ff05::a",
				"ff05::b",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for BGP router addresses?",
			Options: []string{
				"2001:5::/32",
				"2001:10::/28",
				"2001:20::/28",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which ICMPv6 message type is used for Echo Request?",
			Options: []string{
				"Type 128",
				"Type 129",
				"Type 130",
				"Type 131",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for root server addresses?",
			Options: []string{
				"2001:7f8::/29",
				"2001:500::/32",
				"2001:501::/32",
				"2001:502::/32",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all OSPF routers?",
			Options: []string{
				"ff02::5",
				"ff02::6",
				"ff05::5",
				"ff05::6",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 prefix for TLA assignments?",
			Options: []string{
				"2001::/16",
				"2002::/16",
				"2003::/16",
				"fc00::/7",
			},
			Answer: 0,
		},
		{
			Text: "Which ICMPv6 message type is used for Echo Reply?",
			Options: []string{
				"Type 128",
				"Type 129",
				"Type 130",
				"Type 131",
			},
			Answer: 1,
		},
		{
			Text: "What is the IPv6 prefix for NLA assignments?",
			Options: []string{
				"2001:0300::/23",
				"2001:0400::/23",
				"2001:0500::/23",
				"fc00::/7",
			},
			Answer: 0,
		},

		{
			Text: "What is the length of an IPv6 address in bits?",
			Options: []string{
				"32 bits",
				"64 bits",
				"128 bits",
				"256 bits",
			},
			Answer: 2,
		},
		{
			Text: "Which of these is a valid IPv6 address?",
			Options: []string{
				"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				"2001:db8:85a3::8a2e:370:7334",
				"::1",
				"All of the above",
			},
			Answer: 3,
		},
		{
			Text: "What does the '::' notation represent in an IPv6 address?",
			Options: []string{
				"An error in the address",
				"A sequence of one or more groups of 16 zero bits",
				"A multicast address",
				"A reserved address for documentation",
			},
			Answer: 1,
		},
		{
			Text: "Which prefix is used for link-local IPv6 addresses?",
			Options: []string{
				"2000::/3",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 2,
		},
		{
			Text: "What is the IPv6 loopback address?",
			Options: []string{
				"127.0.0.1",
				"::1",
				"::",
				"fe80::1",
			},
			Answer: 1,
		},
		{
			Text: "Which protocol is used for IPv6 neighbor discovery?",
			Options: []string{
				"ARP",
				"NDP",
				"DHCPv6",
				"ICMPv6",
			},
			Answer: 1,
		},
		{
			Text: "What is the purpose of the Unique Local Address (ULA) in IPv6?",
			Options: []string{
				"Global unicast addressing",
				"Private addressing within a site",
				"Multicast addressing",
				"Link-local communication",
			},
			Answer: 1,
		},
		{
			Text: "Which field is NOT present in the IPv6 header?",
			Options: []string{
				"Version",
				"Traffic Class",
				"Header Checksum",
				"Hop Limit",
			},
			Answer: 2,
		},
		{
			Text: "What is the minimum MTU for IPv6?",
			Options: []string{
				"576 bytes",
				"1280 bytes",
				"1500 bytes",
				"9000 bytes",
			},
			Answer: 1,
		},
		{
			Text: "Which ICMPv6 message type is used for router solicitation?",
			Options: []string{
				"Type 133",
				"Type 134",
				"Type 135",
				"Type 136",
			},
			Answer: 0,
		},
		{
			Text: "What is the IPv6 equivalent of IPv4's 0.0.0.0?",
			Options: []string{
				"::",
				"::1",
				"ff02::1",
				"fe80::1",
			},
			Answer: 0,
		},
		{
			Text: "Which multicast address represents all IPv6 nodes?",
			Options: []string{
				"ff02::1",
				"ff02::2",
				"ff05::1",
				"ff05::2",
			},
			Answer: 0,
		},
		{
			Text: "How many bits are used for the interface identifier in an IPv6 address?",
			Options: []string{
				"32 bits",
				"48 bits",
				"64 bits",
				"128 bits",
			},
			Answer: 2,
		},
		{
			Text: "Which protocol is used for automatic address configuration in IPv6?",
			Options: []string{
				"DHCPv6",
				"SLAAC",
				"NDP",
				"Both B and C",
			},
			Answer: 3,
		},
		{
			Text: "What is the purpose of the Flow Label field in the IPv6 header?",
			Options: []string{
				"To identify packets that belong to the same flow",
				"To indicate the priority of the packet",
				"To specify the next header type",
				"To prevent packet fragmentation",
			},
			Answer: 0,
		},
		{
			Text: "Which extension header is used for fragmentation in IPv6?",
			Options: []string{
				"Hop-by-Hop Options",
				"Routing",
				"Fragment",
				"Destination Options",
			},
			Answer: 2,
		},
		{
			Text: "What is the IPv6 prefix for multicast addresses?",
			Options: []string{
				"2000::/3",
				"fc00::/7",
				"fe80::/10",
				"ff00::/8",
			},
			Answer: 3,
		},
		{
			Text: "Which IPv6 address type is used for one-to-nearest communication?",
			Options: []string{
				"Unicast",
				"Anycast",
				"Multicast",
				"Broadcast",
			},
			Answer: 1,
		},
		{
			Text: "What is the purpose of the Router Advertisement message in IPv6?",
			Options: []string{
				"To inform hosts of available routers",
				"To assign IPv6 addresses to hosts",
				"To resolve IPv6 addresses to MAC addresses",
				"To manage multicast group membership",
			},
			Answer: 0,
		},
	}
}

func initializeTemplates() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	var err error
	quizTemplate, err = template.New("quiz").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>IPv6 Quiz</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .question { margin-bottom: 20px; padding: 15px; background-color: #f5f5f5; border-radius: 5px; }
        .options { margin-left: 20px; }
        .correct { color: green; font-weight: bold; }
        .incorrect { color: red; }
        .result { font-size: 1.2em; margin: 20px 0; padding: 15px; background-color: #e6f7ff; border-radius: 5px; }
        button { padding: 8px 15px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #45a049; }
    </style>
</head>
<body>
    <h1>IPv6 Knowledge Quiz</h1>
    {{if .ShowResults}}
        <div class="result">
            <h2>Quiz Results</h2>
            <p>Your score: {{.Score}}/{{.Total}} ({{printf "%.1f" .Percentage}}%)</p>
        </div>
        <h3>Review:</h3>
        {{range $index, $question := .Questions}}
        <div class="question">
            <p><strong>Question {{add $index 1}}:</strong> {{$question.Text}}</p>
            <div class="options">
                {{range $i, $option := $question.Options}}
                    {{if eq $i (index $.CorrectAnswers $index)}}
                        <p class="correct">{{add $i 1}}. {{$option}} (Correct Answer)</p>
                    {{else if eq $i (index $.UserAnswers $index)}}
                        <p class="incorrect">{{add $i 1}}. {{$option}} (Your Answer)</p>
                    {{else}}
                        <p>{{add $i 1}}. {{$option}}</p>
                    {{end}}
                {{end}}
            </div>
        </div>
        {{end}}
        <button onclick="window.location.href='/'">Start New Quiz</button>
    {{else}}
        <form method="POST" action="/submit">
            {{range $index, $question := .Questions}}
            <div class="question">
                <p><strong>Question {{add $index 1}}:</strong> {{$question.Text}}</p>
                <div class="options">
                    {{range $i, $option := $question.Options}}
                        <input type="radio" id="q{{$index}}o{{$i}}" name="q{{$index}}" value="{{$i}}" required>
                        <label for="q{{$index}}o{{$i}}">{{$option}}</label><br>
                    {{end}}
                </div>
            </div>
            {{end}}
            <button type="submit">Submit Quiz</button>
        </form>
    {{end}}
</body>
</html>
`)
	if err != nil {
		panic("Failed to parse template: " + err.Error())
	}
}

func startWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleQuiz)
	mux.HandleFunc("/submit", handleSubmit)

	server := &http.Server{
		Addr:    "[::1]:5000",
		Handler: mux,
	}

	fmt.Println("Starting IPv6 web server on [::1]:5000")
	err := server.ListenAndServe()
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

func handleQuiz(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Select 20 random questions
	currentQuiz = selectRandomQuestions(questions, 20)
	userAnswers = make([]int, len(currentQuiz))
	correctAnswers = make([]int, len(currentQuiz))
	for i, q := range currentQuiz {
		correctAnswers[i] = q.Answer
	}

	renderQuiz(w, false, QuizResult{}, currentQuiz, userAnswers, correctAnswers)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	score := 0
	for i := range currentQuiz {
		ans, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("q%d", i)))
		userAnswers[i] = ans
		if ans == currentQuiz[i].Answer {
			score++
		}
	}

	total := len(currentQuiz)
	percentage := float64(score) / float64(total) * 100
	result := QuizResult{
		Score:      score,
		Total:      total,
		Percentage: percentage,
	}

	renderQuiz(w, true, result, currentQuiz, userAnswers, correctAnswers)
}
func renderQuiz(w http.ResponseWriter, showResults bool, result QuizResult, questions []Question, userAnswers, correctAnswers []int) {
	data := struct {
		ShowResults    bool
		Result         QuizResult
		Questions      []Question
		UserAnswers    []int
		CorrectAnswers []int
	}{
		ShowResults:    showResults,
		Result:         result,
		Questions:      questions,
		UserAnswers:    userAnswers,
		CorrectAnswers: correctAnswers,
	}

	err := quizTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func selectRandomQuestions(allQuestions []Question, count int) []Question {
	if count > len(allQuestions) {
		count = len(allQuestions)
	}

	indexes := rand.Perm(len(allQuestions))
	selected := make([]Question, count)
	for i := 0; i < count; i++ {
		selected[i] = allQuestions[indexes[i]]
	}

	return selected
}
