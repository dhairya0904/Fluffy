package monitor

import (
	"fluffy/domain"
	"time"
)

var line = "--------------------------------------------------------------------\n"

var DOG string = "         (__)\n         (oo)\n   /------\\/\n  / |    ||\n *  /\\---/\\\n    ~~   ~~\n"
var REPORT_MAX_HITS string = line + "\t\tReport for DATE: %s\n" + line + "PATH\t\t\tCOUNT\n" + line
var ROW = "%s\t\t\t%d\n"

var SyntheticData = []domain.Event{
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/report/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/login/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/report/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/login/user", Status: 502, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/report/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/logout/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/login/user", Status: 503, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/report/user", Status: 501, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/logout/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
	domain.Event{Host: "127.0.0.1", User: "james", Time: time.Now(), URL: "/math/user", Status: 200, Method: "GET", Protocol: "HTTP/1.0"},
}
