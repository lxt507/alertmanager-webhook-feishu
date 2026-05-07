package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/alertmanager-webhook-feishu/feishu"
	"github.com/xujiahua/alertmanager-webhook-feishu/model"
)

type Server struct {
	bots          map[string]feishu.IBot
	splitByStatus bool
}

func New(bots map[string]feishu.IBot, splitByStatus bool) *Server {
	s := &Server{
		bots:          bots,
		splitByStatus: splitByStatus,
	}
	return s
}

func (s Server) hook(w http.ResponseWriter, r *http.Request) {
	// get path param
	vars := mux.Vars(r)
	group := vars["group"]
	bot, ok := s.bots[group]
	if !ok {
		logrus.Errorf("group not found: %s", group)
		http.Error(w, "group not found", http.StatusBadRequest)
		return
	}

	// verify signature
	if err := s.verifySignature(r, bot); err != nil {
		logrus.Errorf("signature verification failed for group %s: %s", group, err)
		http.Error(w, "signature verification failed", http.StatusUnauthorized)
		return
	}

	// get body param
	var alerts model.WebhookMessage
	err := json.NewDecoder(r.Body).Decode(&alerts)
	if err != nil {
		logrus.Errorf("cannot parse content, %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		spew.Dump(alerts)
	}

	// get query string
	meta := make(map[string]string)
	for key, values := range r.URL.Query() {
		meta[key] = strings.Join(values, ",")
	}
	// also include path param
	meta["group"] = group

	var alertsGroups []model.WebhookMessage
	if s.splitByStatus {
		alertsGroups = split(alerts)
	} else {
		alertsGroups = []model.WebhookMessage{alerts}
	}

	for _, alerts := range alertsGroups {
		alerts.Meta = meta
		err = bot.Send(&alerts)
		if err != nil {
			logrus.Errorf("cannot send alerts, %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, _ = fmt.Fprintf(w, "ok")
}

func split(alerts model.WebhookMessage) []model.WebhookMessage {
	var groups []model.WebhookMessage
	if len(alerts.Alerts.Firing()) != 0 {
		alertsClone := alerts
		alertsClone.Alerts = alerts.Alerts.Firing()
		groups = append(groups, alertsClone)
	}
	if len(alerts.Alerts.Resolved()) != 0 {
		alertsClone := alerts
		alertsClone.Alerts = alerts.Alerts.Resolved()
		groups = append(groups, alertsClone)
	}
	return groups
}

func (s Server) verifySignature(r *http.Request, bot feishu.IBot) error {
	secret := bot.GetSecret()
	if secret == "" {
		// 未配置签名密钥，跳过验证
		return nil
	}

	// 从请求头获取签名和时间戳
	timestamp := r.Header.Get("X-Lark-Request-Timestamp")
	sign := r.Header.Get("X-Lark-Request-Signature")

	if timestamp == "" || sign == "" {
		return fmt.Errorf("missing signature headers")
	}

	// 解析时间戳
	var ts int64
	_, err := fmt.Sscanf(timestamp, "%d", &ts)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %s", timestamp)
	}

	// 验证签名
	return feishu.VerifySign(secret, ts, sign)
}

func (s Server) health(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s Server) reload(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s Server) Start(address string) error {
	r := mux.NewRouter()
	r.HandleFunc("/hook/{group}", s.hook).Methods("POST")

	// management etc...
	sr := r.PathPrefix("/-").Subrouter()
	sr.HandleFunc("/healthz", s.health).Methods("GET")
	sr.HandleFunc("/reload", s.health).Methods("GET")

	// prometheus
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}
