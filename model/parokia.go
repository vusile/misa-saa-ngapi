package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

var parokiaIndexName = "parokia"

type Parokia struct {
	gorm.Model
	ID        uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name      string
	Location  string
	JimboID   uint
	Jimbo     Jimbo
	IsKigango bool
	UserID    uint
	User      User
	Slug      string
	Timings   []Timing
}

func (p *Parokia) AddToIndex(esClient *elasticsearch.Client) {
	document := struct {
		Name     string `json:"name"`
		Jimbo    string `json:"jimbo"`
		Location string `json:"location"`
	}{p.Name, p.Jimbo.Name, p.Location}

	data, err := json.Marshal(document)

	if err != nil {
		fmt.Println("failed to marshal:", err)
	}

	req := esapi.IndexRequest{
		Index:      parokiaIndexName,
		DocumentID: strconv.Itoa(int(p.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), esClient)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("Indexing error", err)
	}
}

func (p *Parokia) GenerateTimings() map[uint][]Timing {
	var timings = map[uint][]Timing{}
	for _, time := range p.Timings {
		timings[time.HudumaID] = append(timings[time.HudumaID], time)
	}

	return timings
}
