package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
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
			Address: "localhost:3005",
		},
		&model.Node{
			Address: "localhost:3006",
		},
		// &model.Node{
		// 	Address: "localhost:3007",
		// },
		// &model.Node{
		// 	Address: "localhost:3008",
		// },
		// &model.Node{
		// 	Address: "localhost:3009",
		// },
	}

	start := time.Now()

	for i := 0; i < 5; i++ {
		valSet = reloadVals(valSet)

		setAllVals(valSet, nodes)
	}

	//checkAllVals(valSet, nodes)

	for i := 0; i < 5; i++ {
		valSet = reloadVals(valSet)

		setAllVals(valSet, nodes)
	}

	//checkAllVals(valSet, nodes)

	finish := time.Now()

	duration := finish.Sub(start)
	fmt.Printf("Test took %f s\n", duration.Seconds())
}

func setAllVals(valSet map[string]string, nodes []*model.Node) {
	resultChan := make(chan error)
	count := 0
	numFailed := 0
	numSucceeded := 0

	for key, val := range valSet {
		index := mrand.Intn(len(nodes))
		node := nodes[index]

		count++
		go setVal(key, val, node, resultChan)
	}

	for true {
		if err := <-resultChan; err != nil {
			fmt.Println(err.Error())
			numFailed++
		} else {
			numSucceeded++
		}

		if numFailed+numSucceeded == count {
			break
		}
	}

	fmt.Printf("Set %d values with %d successes and %d failures\n", count, numSucceeded, numFailed)
}

func setVal(key, val string, node *model.Node, result chan error) {
	setValRequest := &requests.SetValueRequest{
		Key:   key,
		Value: val,
	}

	if err := send.SetValue(setValRequest, node); err != nil {
		result <- err
	}

	result <- nil
}

func checkAllVals(valSet map[string]string, nodes []*model.Node) {
	resultChan := make(chan error)
	count := 0
	numFailed := 0
	numSucceeded := 0

	for key, val := range valSet {
		index := mrand.Intn(len(nodes))
		node := nodes[index]

		count++
		go checkVal(key, val, node, resultChan)
	}

	for true {
		if err := <-resultChan; err != nil {
			fmt.Println(err.Error())
			numFailed++
		} else {
			numSucceeded++
		}

		if numFailed+numSucceeded == count {
			break
		}
	}

	fmt.Printf("Got %d values with %d successes and %d failures\n", count, numSucceeded, numFailed)
}

func checkVal(key, val string, node *model.Node, result chan error) {
	url := transport.URLFromAddressAndPath(node.Address, "v1/value/"+key)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	res, err := transport.HttpClient().Do(req)
	if err != nil {
		result <- err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		result <- err
	}
	defer res.Body.Close()

	gotVal := string(body)
	if gotVal != val {
		result <- fmt.Errorf("Got %q, shold be %q", gotVal, val)
	}

	result <- nil
}

func randomString() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}

func loadVals() map[string]string {
	newSet := make(map[string]string)

	for i := 0; i < 10; i++ {
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
