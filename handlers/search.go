package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CodeWithSameera/Vehicles/helpers"
	"github.com/olivere/elastic/v7"
)

type Vehicle struct {
	VNo         string  `json:"vno"`
	OwnerName          string   `json:"name"`
}
func GetESClient() (*elastic.Client, error) {

client, err :=  elastic.NewClient(elastic.SetURL(helpers.GoDotEnvVariable("ELASTIC_HOST")),
elastic.SetSniff(false),
elastic.SetHealthcheck(false))

fmt.Println("ES initialized...")

return client, err

}

func SaveEntity(vNo, owner string) {

	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	//creating vecicle object
	newVehicle := Vehicle{
		VNo:         vNo,
		OwnerName:          owner,
	}

	dataJSON, err := json.Marshal(newVehicle)
	js := string(dataJSON)
	_, err = esclient.Index().
		Index("vehicles").
		BodyJson(js).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println("[Elastic][InsertProduct]Insertion Successful")

}

func GetResults(owner string){
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	var vehicles []Vehicle

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("owner", owner))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	/* until this block */

	searchService := esclient.Search().Index("students").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var vehicle Vehicle
		err := json.Unmarshal(hit.Source, &vehicle)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err != nil {
		fmt.Println("Fetching vehicle fail: ", err)
	} else {
		for _, s := range vehicles {
			fmt.Printf("Student found Name: %s, Age: %d, Score: %f \n", s.VNo, s.OwnerName)
		}
	}
}