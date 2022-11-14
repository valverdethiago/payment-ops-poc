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
	url2 "net/url"
)

const (
	fetchBalancesUrl     string = `%s/accounts/%s/balances`
	fetchTransactionsUrl string = `%s/accounts/%s/transactions?date_from=2022-01-01`
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

func (trioHttpClientImpl *TrioHttpClientImpl) FetchBalancesFromBank(account domain.Account) (*restclient.FetchRequestResponse, error) {
	url := fmt.Sprintf(fetchBalancesUrl, trioHttpClientImpl.config.BasePath, account.ProviderAccountId)
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
	responseObj := &restclient.FetchRequestResponse{}
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

func (trioHttpClientImpl *TrioHttpClientImpl) FetchTransactionsFromBank(account domain.Account) (*restclient.FetchRequestResponse, error) {
	url, err := url2.Parse(fmt.Sprintf(fetchTransactionsUrl, trioHttpClientImpl.config.BasePath, account.ProviderAccountId))
	if err != nil {
		return nil, err
	}
	values := url.Query()
	values.Del("date_from")
	client := &http.Client{}
	request, err := http.NewRequest(postMethod, url.String(), nil)
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
	responseObj := &restclient.FetchRequestResponse{}
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}

func (trioHttpClientImpl *TrioHttpClientImpl) ListTransactions(account domain.Account) (*restclient.ListTransactionsResponse, error) {

	url, err := url2.Parse(fmt.Sprintf(fetchTransactionsUrl, trioHttpClientImpl.config.BasePath, account.ProviderAccountId))
	if err != nil {
		return nil, err
	}
	values := url.Query()
	if account.LastTransactionsUpdateAt == nil {
		values.Del("date_from")
	} else {
		values.Set("date_from", account.LastTransactionsUpdateAt.Format("2006-01-02"))
	}
	url.RawQuery = values.Encode()
	client := &http.Client{}
	request, err := http.NewRequest(getMethod, url.String(), nil)
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
	responseObj := &restclient.ListTransactionsResponse{}
	err = json.Unmarshal(bodyBytes, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, nil
}
