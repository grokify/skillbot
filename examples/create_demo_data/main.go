package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/algoliautil"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/fmt/fmtutil"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type Person struct {
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Title  string   `json:"title"`
	Phone  string   `json:"phone"`
	Avatar string   `json:"avatar"`
	Skills []string `json:"skills"`
}

func (p *Person) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Person) ToAlgoliaObject() algoliasearch.Object {
	return algoliasearch.Object{
		"objectID": p.Email,
		"email":    p.Email,
		"name":     p.Name,
		"title":    p.Title,
		"phone":    p.Phone,
		"avatar":   p.Avatar,
		"skills":   p.Skills}
}

func GetTestData() []Person {
	return []Person{
		/*
						{
							Name:   "John Wang",
							Email:  "john.wang@ringcentral.com",
							Title:  "Sr. Director, Platform",
							Phone:  "1 (650) 555-0101",
							Avatar: "https://glip-vault-1.s3.amazonaws.com/web/customer_files/89317449740/avatar_jwang.jpg?Expires=2075494478\u0026AWSAccessKeyId=AKIAJROPQDFTIHBTLJJQ\u0026Signature=q20iqGOhF73kHSiJclQqdKcBdzE%3D",
							Skills: []string{"platform", "API", "SDKs", "Chat", "product management", "bots", "go", "golang", "ruby"},
						},
						{
							Name:   "Tyler Liu",
							Email:  "tyler.liu@ringcentral.com",
							Title:  "Sr. Developer Relations Engineer",
							Phone:  "1 (650) 555-0102",
							Skills: []string{"platform", "API", "bots", "react", "vue", "SDKs", "C#", "Ruby", "Python", "JavaScript"},
						},
			{
				Name:   "David Theis",
				Title:  "Airport Coordinator",
				Skills: []string{"FlightOps", "Airport GPS", "Airport Navigation", "Flight Reservations"},
			},

						Name: David Theis
			Title: Delta Flight Specialist
			Skills: FlightOps, Delta, Flight Reservations

		*/
		{
			Name:   "David Theis",
			Email:  "david.theis@ringcentral.com",
			Title:  "Delta Flight Specialist",
			Skills: []string{"FlightOps", "Delta", "Flight Reservations"},
		},
	}
}

func PersonsToObjects(persons []Person) []algoliasearch.Object {
	objects := []algoliasearch.Object{}
	for _, p := range persons {
		objects = append(objects, p.ToAlgoliaObject())
	}
	return objects
}

func GetIndex(config []byte, indexName string) (*search.Index, error) {
	client, err := algoliautil.NewClientJSON(config)
	if err != nil {
		return nil, err
	}
	index := client.InitIndex(indexName)
	return index, nil
}

func IndexAlgoliaPersons(index *search.Index, persons []Person) error {
	if index == nil {
		return errors.New("algolia index cannot be nil")
	}
	/*
		index, err := GetIndex(config, indexName)
		client, err := algoliautil.NewClientFromJSONAdmin([]byte(os.Getenv("ALGOLIA_APP_CREDENTIALS_JSON")))
		if err != nil {
			return err
		}
		index := client.InitIndex("expertskills")
	*/
	if 1 == 0 {
		res, err := index.SaveObjects(PersonsToObjects(persons))
		if err != nil {
			return err
		}
		fmtutil.PrintJSON(res)
	}
	if 1 == 1 {
		// res, err := index.UpdateObjects(PersonsToObjects(persons))
		res, err := index.ReplaceAllObjects(PersonsToObjects(persons))
		if err != nil {
			return err
		}
		fmtutil.PrintJSON(res)
	}
	return nil
}

func WriteCsvPersons(persons []Person) error {
	w, err := csvutil.NewWriter("users.csv", ",", false, "")
	if err != nil {
		return err
	}

	for _, p := range persons {
		j, err := p.ToJSON()
		if err == nil {
			w.AddLine([]interface{}{p.Name, string(j)})
		}
	}
	w.Close()

	f, err := os.Create("users.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	c2 := csv.NewWriter(f)
	for _, p := range persons {
		skills := strings.Join(p.Skills, ", ")
		j, err := p.ToJSON()
		if err == nil {
			err := c2.Write([]string{p.Name, p.Name + "; " + skills + "; " + string(j)})
			if err != nil {
				return err
			}
		}
	}
	c2.Flush()
	return nil
}

func SearchAndDelete(index *search.Index, qry string) error {
	res, err := index.Search(qry, nil)
	if err != nil {
		return err
	}
	fmtutil.PrintJSON(res)
	fmt.Printf("NUM_HITS [%v]\n", len(res.Hits))
	for _, hit := range res.Hits {
		objectID := hit["objectID"]
		fmt.Printf("ObjectId [%v]\n", objectID)
		if 1 == 0 {
			res, err := index.DeleteObject(objectID.(string))
			if err != nil {
				return err
			}
			fmtutil.PrintJSON(res)
		}
	}

	return nil
}

func main() {
	_, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "../.env"}, 1)
	//err := config.LoadDotEnvSkipEmpty("../.env")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(os.Getenv("ALGOLIA_APP_CREDENTIALS_JSON"))

	indexPtr, err := GetIndex(
		[]byte(os.Getenv("ALGOLIA_APP_CREDENTIALS_JSON")),
		os.Getenv("ALGOLIA_INDEX"))
	if err != nil {
		log.Fatal(err)
	}
	// index := *indexPtr

	persons := GetTestData()

	indexAlgoliaPersons := false
	writeCsvPersons := false
	searchAndDelete := true

	if indexAlgoliaPersons {
		if err := IndexAlgoliaPersons(indexPtr, persons); err != nil {
			log.Fatal(err)
		}
	}

	if writeCsvPersons {
		if err := WriteCsvPersons(persons); err != nil {
			log.Fatal(err)
		}
	}

	if searchAndDelete {
		err := SearchAndDelete(indexPtr, "flight")
		if err != nil {
			log.Fatal(err)
		}
	}
}
