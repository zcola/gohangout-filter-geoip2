package main

import (
	"net"
	"github.com/oschwald/geoip2-golang"
	"github.com/childe/gohangout/value_render"
	"github.com/golang/glog"
)


type GeoIP2Filter struct {
	config map[interface{}]interface{}
	src string
	srcVR value_render.ValueRender
	target string
	language string
	dbPath string
	db *geoip2.Reader
}


func New(config map[interface{}]interface{}) interface{} {
	plugin := &GeoIP2Filter{
		config: config,
		target: "geoip",
	}
	if src, ok := config["src"]; ok {
		plugin.src = src.(string)
		plugin.srcVR = value_render.GetValueRender2(plugin.src)
	} else {
		glog.Fatal("src must be set in GeoIP2 filter plugin")
	}
	if language, ok := config["language"]; ok {
		plugin.language = language.(string)
	} else {
		plugin.language = "en"
	}
	if target, ok := config["target"]; ok {
		plugin.target = target.(string)
	}
	if dbPath, ok := config["dbPath"]; ok {
		plugin.dbPath = dbPath.(string)
	} else {
		glog.Fatal("dbPath must be set in GeoIP2 filter plugin")
	}
	db, err := geoip2.Open(plugin.dbPath)
	if err != nil {
		glog.Fatalf("Failed to open GeoIP2 database: %v", err)
	}
	plugin.db = db
	return plugin
}

func (plugin *GeoIP2Filter) Filter(event map[string]interface{}) (map[string]interface{}, bool) {
	ipAddress := plugin.srcVR.Render(event)
	if ipAddressStr, ok := ipAddress.(string); ok {
		ip := net.ParseIP(ipAddressStr)
		if ip == nil {
			glog.V(10).Infof("Invalid IP address: %s", ipAddressStr)
			return event, false
		}
		record, err := plugin.db.City(ip)
		if err != nil {
			glog.V(10).Infof("Failed to lookup IP address %s: %v using City method", ipAddressStr, err)
			ispRecord, ispErr := plugin.db.ISP(ip)
			if ispErr != nil {
				glog.V(10).Infof("Failed to lookup IP address %s using ISP method: %v", ipAddressStr, ispErr)
				return event, false
			}
			geoData := make(map[string]interface{})
			geoData["isp"] = ispRecord.ISP
			if existingGeoIP, ok := event[plugin.target]; ok {
				if geoIPMap, ok := existingGeoIP.(map[string]interface{}); ok {
					for key, value := range geoIPMap {
						geoData[key] = value
					}
				}
			}
			event[plugin.target] = geoData
			return event, true
		}
		geoData := make(map[string]interface{})
		geoData["timezone"] = record.Location.TimeZone
		geoData["city_name"] = record.City.Names[plugin.language]
		if len(record.Subdivisions) > 0 {
			geoData["region_name"] = record.Subdivisions[0].Names[plugin.language]
		}
		geoData["country_name"] = record.Country.Names[plugin.language]
		type Coordinates struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}
		coordinates := Coordinates{Lon: record.Location.Longitude, Lat: record.Location.Latitude}
		geoData["location"] = map[string]interface{}{
    			"lon": coordinates.Lon,
    			"lat": coordinates.Lat,
		}

		event[plugin.target] = geoData
		return event, true
	} else {
		glog.V(10).Infof("Invalid IP address: %v", ipAddress)
		return event, false
	}
}
