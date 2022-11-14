package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"io"
	"net/http"
)

const (
	fetchBalancesUrl     string = `%s/accounts/%s/balances`
	fetchTransactionsUrl string = `%s/accounts/%s/transactions`
	postMethod           string = "POST"
	getMethod            string = "GET"
)

type TrioHttpClientImpl struct {
	config *config.Config
}

func NewTrioHttpClient(config *config.Config) domain.TrioClient {
	return &TrioHttpClientImpl{
		config: config,
	}
}

func (trioHttpClientImpl *TrioHttpClientImpl) SetBasicAuth(request *http.Request) {
	request.SetBasicAuth(trioHttpClientImpl.config.ClientID, trioHttpClientImpl.config.ClientSecret)
}

func (trioHttpClientImpl *TrioHttpClientImpl) FetchBalancesFromBank(AccountId string) (*restclient.FetchRequestResponse, error) {
	url := fmt.Sprintf(fetchBalancesUrl, trioHttpClientImpl.config.BasePath, AccountId)
	client := &http.Client{}
	request, err := http.NewRequest(postMethod, url, nil)
	if err != nil {
		return nil, err
	}
	trioHttpClientImpl.SetBasicAuth(request)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error calling fetch balances API")
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var responseObj *restclient.FetchRequestResponse
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}

func (trioHttpClientImpl *TrioHttpClientImpl) ListBalance(AccountId string) (*restclient.ListBalanceResponse, error) {
	url := fmt.Sprintf(fetchBalancesUrl, trioHttpClientImpl.config.BasePath, AccountId)
	client := &http.Client{}
	request, err := http.NewRequest(getMethod, url, nil)
	if err != nil {
		return nil, err
	}
	trioHttpClientImpl.SetBasicAuth(request)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error calling list balances API")
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var responseObj *restclient.ListBalanceResponse
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}

func (trioHttpClientImpl *TrioHttpClientImpl) FetchTransactionsFromBank(AccountId string) (*restclient.FetchRequestResponse, error) {
	url := fmt.Sprintf(fetchTransactionsUrl, trioHttpClientImpl.config.BasePath, AccountId)
	client := &http.Client{}
	request, err := http.NewRequest(postMethod, url, nil)
	if err != nil {
		return nil, err
	}
	trioHttpClientImpl.SetBasicAuth(request)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error calling fetch transactions API")
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var responseObj *restclient.FetchRequestResponse
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}

func (trioHttpClientImpl *TrioHttpClientImpl) ListTransactions(AccountId string) (*restclient.ListTransactionsResponse, error) {
	url := fmt.Sprintf(fetchBalancesUrl, trioHttpClientImpl.config.BasePath, AccountId)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	trioHttpClientImpl.SetBasicAuth(request)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error calling fetch balances API")
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var responseObj *restclient.ListTransactionsResponse
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}
