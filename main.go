package main

import (
	"flag"
	"log"
	"net"
	"strings"

	realip "github.com/Ferluci/fast-realip"
	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

type GeoIpInfo struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

var geoIp *maxminddb.Reader

func init() {
	db, err := maxminddb.Open("data/country.mmdb")
	if err != nil {
		panic(err)
	}
	geoIp = db
}

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %+v", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	cors(ctx)
	ip := net.ParseIP(realip.FromRequest(ctx))
	ctx.SetContentType("application/json; charset=utf8")
	if ip == nil {
		ctx.WriteString("{\"country_iso\":\"UNKNOWN\"}")
		return
	}
	var record GeoIpInfo
	err := geoIp.Lookup(ip, &record)
	if err != nil || record.Country.ISOCode == "" {
		log.Printf("Fail to load geo ip info from %s, %+v, %+v", ip.String(), err, record)
		ctx.WriteString("{\"country_iso\":\"UNKNOWN\"}")
		return
	}
	ctx.WriteString("{\"country_iso\":\"" + record.Country.ISOCode + "\"}")
}

func cors(ctx *fasthttp.RequestCtx) {
	originHeader := string(ctx.Request.Header.Peek("Origin"))
	method := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))

	headers := []string{}
	if len(ctx.Request.Header.Peek("Access-Control-Request-Headers")) > 0 {
		headers = strings.Split(string(ctx.Request.Header.Peek("Access-Control-Request-Headers")), ",")
	}

	ctx.Response.Header.Set("Access-Control-Allow-Origin", originHeader)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", method)
	if len(headers) > 0 {
		ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
	}
}
