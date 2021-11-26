package http

import (
	"embed"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/geoloc"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"io/fs"
	"net/http"
	"os"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "web")
)

type WebResources struct {
	Res   embed.FS
	Files map[string]string
}

func StartWebServer(seedConfig seednode.Config, webResources fs.FS, ips *[]geoloc.GeolocalizedPeers) {
	// serve static assets
	http.Handle("/", http.FileServer(http.FS(webResources)))

	// serve endpoint
	http.HandleFunc("/api/peers", writeFakePeers)

	// start web server in non-blocking
	// go func() {
	err := http.ListenAndServe(":"+seedConfig.HttpPort, nil)
	logger.Info("HTTP Server started", "port", seedConfig.HttpPort)
	if err != nil {
		panic(err)
	}
	// }()
}

func writePeers(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(&geoloc.ResolvedPeers)
	if err != nil {
		logger.Info("Failed to marshal peers list")
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}

func writeFakePeers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[{\"moniker\":\"2hfay3\",\"ip\":\"85.195.100.119\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"HE\",\"region_name\":\"\",\"city\":\"Frankfurt am Main\",\"zip\":\"60313\",\"lat\":50.1109,\"lon\":8.68213,\"timezone\":\"Europe/Berlin\",\"isp\":\"Host Europe GmbH\",\"org\":\"Shenzhen Topdata Technology Co., Ltd\",\"as\":\"AS20773 Host Europe GmbH\",\"query\":\"85.195.100.119\"},{\"moniker\":\"KalpaTech\",\"ip\":\"178.63.21.33\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"SN\",\"region_name\":\"\",\"city\":\"Falkenstein\",\"zip\":\"08223\",\"lat\":50.475,\"lon\":12.365,\"timezone\":\"Europe/Berlin\",\"isp\":\"Hetzner Online GmbH\",\"org\":\"Hetzner\",\"as\":\"AS24940 Hetzner Online GmbH\",\"query\":\"178.63.21.33\"},{\"moniker\":\"cc-sentry3\",\"ip\":\"46.101.202.54\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"HE\",\"region_name\":\"\",\"city\":\"Frankfurt am Main\",\"zip\":\"60313\",\"lat\":50.1188,\"lon\":8.6843,\"timezone\":\"Europe/Berlin\",\"isp\":\"DigitalOcean, LLC\",\"org\":\"Digital Ocean\",\"as\":\"AS14061 DigitalOcean, LLC\",\"query\":\"46.101.202.54\"},{\"moniker\":\"validatus-archive\",\"ip\":\"193.26.156.221\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"BW\",\"region_name\":\"\",\"city\":\"Karlsruhe\",\"zip\":\"76185\",\"lat\":49.0291,\"lon\":8.35695,\"timezone\":\"Europe/Berlin\",\"isp\":\"netcup GmbH\",\"org\":\"netcup GmbH\",\"as\":\"AS197540 netcup GmbH\",\"query\":\"193.26.156.221\"},{\"moniker\":\"validatus\",\"ip\":\"45.142.178.239\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"BW\",\"region_name\":\"\",\"city\":\"Karlsruhe\",\"zip\":\"76185\",\"lat\":49.0291,\"lon\":8.35695,\"timezone\":\"Europe/Berlin\",\"isp\":\"netcup GmbH\",\"org\":\"netcup GmbH\",\"as\":\"AS197540 netcup GmbH\",\"query\":\"45.142.178.239\"},{\"moniker\":\"Kalpatech\",\"ip\":\"194.163.181.100\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"NW\",\"region_name\":\"\",\"city\":\"Düsseldorf\",\"zip\":\"40599\",\"lat\":51.1878,\"lon\":6.8607,\"timezone\":\"Europe/Berlin\",\"isp\":\"Contabo GmbH\",\"org\":\"Contabo GmbH\",\"as\":\"AS51167 Contabo GmbH\",\"query\":\"194.163.181.100\"},{\"moniker\":\"BlockNgine-IBC\",\"ip\":\"91.109.29.72\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"HE\",\"region_name\":\"\",\"city\":\"Frankfurt am Main\",\"zip\":\"60326\",\"lat\":50.0996,\"lon\":8.63573,\"timezone\":\"Europe/Berlin\",\"isp\":\"LeaseWeb DE\",\"org\":\"Leaseweb Deutschland GmbH\",\"as\":\"AS28753 Leaseweb Deutschland GmbH\",\"query\":\"91.109.29.72\"},{\"moniker\":\"fullnode\",\"ip\":\"157.245.213.0\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"NJ\",\"region_name\":\"\",\"city\":\"Clifton\",\"zip\":\"07014\",\"lat\":40.8364,\"lon\":-74.1403,\"timezone\":\"America/New_York\",\"isp\":\"DigitalOcean, LLC\",\"org\":\"DigitalOcean, LLC\",\"as\":\"AS14061 DigitalOcean, LLC\",\"query\":\"157.245.213.0\"},{\"moniker\":\"Another-Stargaze\",\"ip\":\"44.197.170.144\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"VA\",\"region_name\":\"\",\"city\":\"Ashburn\",\"zip\":\"20149\",\"lat\":39.0438,\"lon\":-77.4874,\"timezone\":\"America/New_York\",\"isp\":\"Amazon.com\",\"org\":\"AWS EC2 (us-east-1)\",\"as\":\"AS14618 Amazon.com, Inc.\",\"query\":\"44.197.170.144\"},{\"moniker\":\"aafc1cf2-d37b-4f63-b186-eaf9eae99a2e\",\"ip\":\"3.145.19.19\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"OH\",\"region_name\":\"\",\"city\":\"Dublin\",\"zip\":\"43017\",\"lat\":40.0992,\"lon\":-83.1141,\"timezone\":\"America/New_York\",\"isp\":\"Amazon.com, Inc.\",\"org\":\"AWS EC2 (us-east-2)\",\"as\":\"AS16509 Amazon.com, Inc.\",\"query\":\"3.145.19.19\"},{\"moniker\":\"witval\",\"ip\":\"138.68.20.186\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"CA\",\"region_name\":\"\",\"city\":\"Santa Clara\",\"zip\":\"95051\",\"lat\":37.3417,\"lon\":-121.9753,\"timezone\":\"America/Los_Angeles\",\"isp\":\"DigitalOcean, LLC\",\"org\":\"Digital Ocean\",\"as\":\"AS14061 DigitalOcean, LLC\",\"query\":\"138.68.20.186\"},{\"moniker\":\"Slaanesh\",\"ip\":\"147.182.254.3\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"CA\",\"region_name\":\"\",\"city\":\"Santa Clara\",\"zip\":\"95054\",\"lat\":37.3931,\"lon\":-121.962,\"timezone\":\"America/Los_Angeles\",\"isp\":\"DigitalOcean, LLC\",\"org\":\"DigitalOcean, LLC\",\"as\":\"AS14061 DigitalOcean, LLC\",\"query\":\"147.182.254.3\"},{\"moniker\":\"ythsfg53gg1\",\"ip\":\"52.9.118.174\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"CA\",\"region_name\":\"\",\"city\":\"San Jose\",\"zip\":\"95141\",\"lat\":37.3394,\"lon\":-121.895,\"timezone\":\"America/Los_Angeles\",\"isp\":\"Amazon.com, Inc.\",\"org\":\"AWS EC2 (us-west-1)\",\"as\":\"AS16509 Amazon.com, Inc.\",\"query\":\"52.9.118.174\"},{\"moniker\":\"aafc1cf2-d37b-4f63-b186-eaf9eae99a2e\",\"ip\":\"3.133.83.17\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"OH\",\"region_name\":\"\",\"city\":\"Dublin\",\"zip\":\"43017\",\"lat\":40.0992,\"lon\":-83.1141,\"timezone\":\"America/New_York\",\"isp\":\"Amazon.com, Inc.\",\"org\":\"AWS EC2 (us-east-2)\",\"as\":\"AS16509 Amazon.com, Inc.\",\"query\":\"3.133.83.17\"},{\"moniker\":\"node\",\"ip\":\"34.229.83.207\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"VA\",\"region_name\":\"\",\"city\":\"Ashburn\",\"zip\":\"20149\",\"lat\":39.0438,\"lon\":-77.4874,\"timezone\":\"America/New_York\",\"isp\":\"Amazon.com, Inc.\",\"org\":\"AWS EC2 (us-east-1)\",\"as\":\"AS14618 Amazon.com, Inc.\",\"query\":\"34.229.83.207\"},{\"moniker\":\"btinn0\",\"ip\":\"116.203.77.250\",\"status\":\"success\",\"country\":\"Germany\",\"country_code\":\"\",\"region\":\"BY\",\"region_name\":\"\",\"city\":\"Nuremberg\",\"zip\":\"90403\",\"lat\":49.4521,\"lon\":11.0767,\"timezone\":\"Europe/Berlin\",\"isp\":\"Hetzner Online GmbH\",\"org\":\"Hetzner\",\"as\":\"AS24940 Hetzner Online GmbH\",\"query\":\"116.203.77.250\"},{\"moniker\":\"osmosis\",\"ip\":\"159.223.125.28\",\"status\":\"success\",\"country\":\"United States\",\"country_code\":\"\",\"region\":\"NJ\",\"region_name\":\"\",\"city\":\"North Bergen\",\"zip\":\"07047\",\"lat\":40.793,\"lon\":-74.0247,\"timezone\":\"America/New_York\",\"isp\":\"DigitalOcean, LLC\",\"org\":\"DigitalOcean, LLC\",\"as\":\"AS14061 DigitalOcean, LLC\",\"query\":\"159.223.125.28\"}]"))
}
