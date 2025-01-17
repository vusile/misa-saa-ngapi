package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/csrf"
	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

type AdminHandler struct {
	Client   *gorm.DB
	ESClient *elasticsearch.Client
}

func (adminHandler *AdminHandler) Home(w http.ResponseWriter, r *http.Request) {
	var response struct {
		SearchResults []model.Parokia
		Token         map[string]interface{}
		Title         string
	}

	response.Token = SetupToken(r)
	response.Title = "Muda wa Ibada na Huduma Tanzania"

	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/backend/home/index.html",
		"/go/src/app/views/backend/template.html"))

	tmpl.Execute(w, response)
}

func (adminHandler *AdminHandler) Search(w http.ResponseWriter, r *http.Request) {
	var response struct {
		SearchResults []model.Parokia
		Token         map[string]interface{}
		Title         string
	}

	response.Token = make(map[string]interface{})
	response.Token[csrf.TemplateTag] = csrf.TemplateField(r)
	response.Title = "Muda wa Ibada na Huduma Tanzania"

	tmpl := template.Must(template.ParseFiles(
		"/go/src/app/views/frontend/home/index.html",
		"/go/src/app/views/frontend/template.html"))

	if r.PostFormValue("search") != "" {

		var buf bytes.Buffer

		// query := map[string]interface{}{
		// 	"query": map[string]interface{}{
		// 		"multi_match": map[string]interface{}{
		// 			"query":  r.PostFormValue("search"),
		// 			"fields": []string{"name", "jimbo", "location"},
		// 		},
		// 	},
		// }

		query := map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query": r.PostFormValue("search"),
					"type":  "bool_prefix",
					"fields": []string{
						"name",
						"name._2gram",
						"name._3gram",
						"jimbo",
						"jimbo._2gram",
						"jimbo._3gram",
						"location",
						"location._2gram",
						"location._3gram",
					},
				},
			},
		}

		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			fmt.Println("Error encoding", err)
		}

		esClient := adminHandler.ESClient
		res, err := esClient.Search(
			esClient.Search.WithIndex("parokia"),
			esClient.Search.WithBody(&buf),
		)

		defer res.Body.Close()

		if err != nil || res.IsError() {
			fmt.Println("searc error", err, res.Body.Close().Error())
			tmpl.Execute(w, response)
		} else {
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				tmpl.Execute(w, response)
			}

			var ids []uint
			if hits, ok := r["hits"].(map[string]interface{}); ok {
				if hitsHits, ok := hits["hits"].([]interface{}); ok {
					for _, hit := range hitsHits {
						if hitMap, ok := hit.(map[string]interface{}); ok {
							if idStr, ok := hitMap["_id"].(string); ok {
								id, _ := strconv.Atoi(idStr)
								ids = append(ids, uint(id))
							}
						}
					}
				}
			}

			var parokia []model.Parokia
			adminHandler.Client.Preload("Jimbo").Find(&parokia, "id in ?", ids)
			response.SearchResults = parokia

			tmpl := template.Must(template.ParseFiles(
				"/go/src/app/views/frontend/home/index.html",
				"/go/src/app/views/frontend/template.html"))

			tmpl.ExecuteTemplate(w, "search-results", response)

		}
	}
}
