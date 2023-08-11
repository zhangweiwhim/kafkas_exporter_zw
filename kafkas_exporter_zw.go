package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	plog "github.com/prometheus/common/promlog"
	plogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/tidwall/gjson"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	_ "net/http"
	"os"
)

func main() {

	var (
		configPath = toFlagString("config.path", "U need input the config path", "/Users/zhangwei/Documents/golang_src/kafkas_exporter_zw/kafka_info.json")
		//configPath = "/Users/zhangwei/Documents/golang_src/kafkas_exporter_zw/kafka_info.json"
		topicFilter  = toFlagString("topic.filter", "Regex that determines which topics to collect.", ".*")
		topicExclude = toFlagString("topic.exclude", "Regex that determines which topics to exclude.", "^$")
		groupFilter  = toFlagString("group.filter", "Regex that determines which consumer groups to collect.", ".*")
		groupExclude = toFlagString("group.exclude", "Regex that determines which consumer groups to exclude.", "^$")
		logSarama    = toFlagBool("log.enable-sarama", "Turn on Sarama logging, default is false.", false, "false")
	)

	klog.InitFlags(flag.CommandLine)
	plConfig := plog.Config{}
	plogflag.AddFlags(kingpin.CommandLine, &plConfig)
	kingpin.Version(version.Print("kafka_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	data, errFile := os.ReadFile(*configPath)
	if nil != errFile {
		klog.Fatalln(errFile)
	} else {
		klog.Infoln("Config json file was read success \r \n ")
	}

	var resp []map[string]interface{}
	clusters := gjson.Parse(string(data)).Get("clusters")
	tmpWLA := gjson.Parse(string(data)).Get("web_listen_address").String()
	tmpWTP := gjson.Parse(string(data)).Get("web_telemetry_path").String()
	tmpWLH := gjson.Parse(string(data)).Get("web_listen_host").String()
	tmpCC := gjson.Parse(string(data)).Get("component_code").String()

	r := gin.Default()
	for _, cluster := range clusters.Array() {
		tmpEnable := cluster.Get("enable").Bool()
		if tmpEnable {
			tmpIN := cluster.Get("instance_name").String()
			tmpCE := cluster.Get("component_env").String()
			tmpCII := cluster.Get("component_instance_id").String()
			tmpII := cluster.Get("instance_id").String()
			tmpURI := []string{cluster.Get("instance_host").String()}
			owner := cluster.Get("owner").String()
			importance := cluster.Get("importance").String()
			targets := []string{tmpWLH + tmpWLA}
			tmpLabels := map[string]interface{}{
				"instance_name":         tmpIN,
				"owner":                 owner,
				"importance":            importance,
				"component_instance_id": tmpCII,
				"component_code":        tmpCC,
				"component_env":         tmpCE,
				"instance_id":           tmpII,
				"instance_host":         tmpURI[0],
			}
			resp = append(resp, map[string]interface{}{
				"targets": targets,
				"labels":  tmpLabels,
			})
			opts := &kafkaOpts{
				uri:                      tmpURI,
				kafkaVersion:             cluster.Get("kafka_version").String(),
				useSASL:                  cluster.Get("sasl_enabled").Bool(),
				useSASLHandshake:         true,
				saslMechanism:            cluster.Get("sasl_mechanism").String(),
				saslUsername:             cluster.Get("sasl_username").String(),
				saslPassword:             cluster.Get("sasl_password").String(),
				labels:                   map[string]string{},
				serviceName:              "",
				kerberosConfigPath:       "",
				realm:                    "",
				kerberosAuthType:         "",
				keyTabPath:               "",
				saslDisablePAFXFast:      false,
				useTLS:                   false,
				tlsServerName:            "",
				tlsCAFile:                "",
				tlsCertFile:              "",
				tlsKeyFile:               "",
				serverUseTLS:             false,
				serverMutualAuthEnabled:  false,
				serverTlsCAFile:          "",
				serverTlsCertFile:        "",
				serverTlsKeyFile:         "",
				tlsInsecureSkipTLSVerify: false,
				useZooKeeperLag:          false,
				uriZookeeper:             []string{"localhost:2181"},
				metadataRefreshInterval:  "30s",
				offsetShowAll:            true,
				allowConcurrent:          false,
				topicWorkers:             100,
				allowAutoTopicCreation:   false,
				verbosityLogLevel:        0,
			}
			reg := setup(*topicFilter, *topicExclude, *groupFilter, *groupExclude, *logSarama, *opts)
			r.GET(fmt.Sprintf("%s/%s", tmpWTP, cluster.Get("instance_name").String()), gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
		}
	}

	r.GET("/metrics")
	r.GET("kafka_sd_targets", func(context *gin.Context) {
		context.JSON(http.StatusOK, resp)
	})
	r.Run(tmpWLA)
}
