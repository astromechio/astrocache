package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"net/http"
	"time"

	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/send"
	"github.com/astromechio/astrocache/transport"

	acrypto "github.com/astromechio/astrocache/crypto"
)

func main() {
	valSet := loadVals()

	nodes := []*model.Node{
		&model.Node{
			Address: "localhost:3003",
		},
		&model.Node{
			Address: "localhost:3005",
		},
		&model.Node{
			Address: "localhost:3006",
		},
		// &model.Node{
		// 	Address: "localhost:3008",
		// },
		// &model.Node{
		// 	Address: "localhost:3009",
		// },
	}

	start := time.Now()

	for i := 0; i < 40; i++ {
		valSet = reloadVals(valSet)

		setAllVals(valSet, nodes)
	}

	checkAllVals(valSet, nodes)

	for i := 0; i < 10; i++ {
		valSet = reloadVals(valSet)

		setAllVals(valSet, nodes)
	}

	checkAllVals(valSet, nodes)

	finish := time.Now()

	duration := finish.Sub(start)
	fmt.Printf("Test took %f s", duration.Seconds())
}

func setAllVals(valSet map[string]string, nodes []*model.Node) {
	for key, val := range valSet {
		setValRequest := &requests.SetValueRequest{
			Key:   key,
			Value: val,
		}

		index := mrand.Intn(len(nodes))
		node := nodes[index]

		if err := send.SetValue(setValRequest, node); err != nil {
			log.Fatal(err)
		}
	}
}

func checkAllVals(valSet map[string]string, nodes []*model.Node) {
	for key, actual := range valSet {
		index := mrand.Intn(len(nodes))
		node := nodes[index]

		url := transport.URLFromAddressAndPath(node.Address, "v1/value/"+key)

		req, _ := http.NewRequest(http.MethodGet, url, nil)

		res, err := http.DefaultClient.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		val := string(body)

		fmt.Printf("Got value %s, should be %s\n", val, actual)
	}

}

func randomString() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}

func loadVals() map[string]string {
	newSet := make(map[string]string)

	for i := 0; i < 100; i++ {
		key := randomString()
		val := randomString()

		newSet[key] = val
	}

	return newSet
}

func reloadVals(valSet map[string]string) map[string]string {
	newSet := make(map[string]string)

	for key := range valSet {
		newVal := randomString()

		newSet[key] = newVal
	}

	return newSet
}
