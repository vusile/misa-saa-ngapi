package application

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/vusile/misa-saa-ngapi/model"
)

var ParokiaIndexName = "parokia"

func CreateESIndex(client *elasticsearch.Client) {
	_, err := esapi.IndicesExistsRequest{
		Index: []string{ParokiaIndexName},
	}.Do(context.Background(), client)

	if err != nil {
		_, err := client.Indices.Create(ParokiaIndexName)

		if err != nil {
			fmt.Println("Create Index Error", err)
		}
	}
}

func IndexParokia(app *App) {
	var parokia []model.Parokia
	tx := app.gorm.Preload("Jimbo").Find(&parokia)

	if tx.Error != nil || tx.RowsAffected == 0 {
		fmt.Println("Maybe No records to index ", tx.Error)
	}

	for _, p := range parokia {
		p.AddToIndex(app.esClient)
	}
}

func SearchAsYouType(client *elasticsearch.Client) {
	mappings := map[string]interface{}{
		// "mappings": map[string]interface{}{
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "search_as_you_type",
			},
			"jimbo": map[string]interface{}{
				"type": "search_as_you_type",
			},
			"location": map[string]interface{}{
				"type": "search_as_you_type",
			},
		},
		// },
	}

	mappingJson, err := json.Marshal(mappings)

	if err != nil {
		fmt.Println("failed to marshal mapping:", err)
	}

	mappingReq := esapi.IndicesPutMappingRequest{
		Index: []string{ParokiaIndexName},
		Body:  bytes.NewReader(mappingJson),
	}

	resp, err := mappingReq.Do(context.Background(), client)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("Mapping Error", err)
	}
}

func ConnectToESClient(config Config) (*elasticsearch.Client, error) {

	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := os.ReadFile("/usr/share/elasticsearch/config/certs/ca/ca.crt")
	if err != nil {
		fmt.Printf("Failed to append cert to RootCAs: %v", err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		fmt.Printf("No certs appended, using system certs only")
	}

	esConfig := elasticsearch.Config{
		MaxRetries: 3,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: *insecure,
				RootCAs:            rootCAs,
			},
		},
		Addresses: []string{
			"https://es01:9200",
		},
		Username: config.ESUserName,
		Password: config.ESPassword,
	}

	return elasticsearch.NewClient(esConfig)
}
