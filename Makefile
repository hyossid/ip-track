
GEO_LITE_VERSION=2022.01.10
databases := GeoLite2-City GeoLite2-Country GeoLite2-ASN

geoip-download:
	mkdir -p data
	@curl -fSL -o data/country.mmdb https://github.com/P3TERX/GeoLite.mmdb/releases/download/${GEO_LITE_VERSION}/GeoLite2-Country.mmdb

prepare: geoip-download
	go mod tidy

build: prepare
	go build
