package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andy-zhangtao/crawlerparam/v1"
)

const (
	APIVERSION = "CRAWLER_API_VERSION"
	CALLBACK   = "CRAWLER_CALL_BACK"
)

var _VERSION_ = "unknown"

func main() {
	if ready, name := isReady(); !ready {
		fmt.Printf("[%s] Env Can't Be Empty!", name)
		os.Exit(-1)
	}

	fmt.Println(getVersion())

	err := callback()
	if err != nil {
		fmt.Printf("CallBack Error [%s]", err.Error())
		os.Exit(-1)
	}

	fmt.Println("PARAM SYNC COMPLETE!")
}

func getVersion() string {
	return fmt.Sprintf("CRAWLER PARAMER VERSION [%s] [%s]", os.Getenv(APIVERSION), _VERSION_)
}

func isReady() (bool, string) {
	if os.Getenv(APIVERSION) == "" {
		return false, APIVERSION
	}

	if os.Getenv(CALLBACK) == "" {
		return false, CALLBACK
	}
	return true, ""
}

func callback() error {
	callURL := os.Getenv(CALLBACK)

	urls := strings.Split(callURL, "|")

	var content string
	var err error
	switch os.Getenv(APIVERSION) {
	case "v1":
		content, err = v1.GetChanMap()
		if err != nil {
			return err
		}
	default:
		return errors.New("API VERSION ERROR! " + os.Getenv(APIVERSION))
	}

	client := &http.Client{}

	for _, u := range urls {
		log.Printf("CALLBACK [%s] ", u)
		req, err := http.NewRequest("POST", u, strings.NewReader(content))
		if err != nil {
			return fmt.Errorf("[%s] error:[%s]", u, err.Error())

		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("[%s] error:[%s]", u, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("[%s] error:[%s]", u, err.Error())
		}
	}

	return nil
}
