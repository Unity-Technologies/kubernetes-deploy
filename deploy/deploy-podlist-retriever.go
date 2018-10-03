package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// KubernetesPodListRetriever retrieves pods via Kubernetes API
type KubernetesPodListRetriever struct {
	Client             *http.Client
	Endpoint           string
	Namespace          string
	BearerTokenService BearerTokenRetriever
}

// PodInformation retrieved from Kubernetes API
func (p *KubernetesPodListRetriever) PodInformation() (*PodList, error) {
	if p.Endpoint == "" || p.Namespace == "" {
		return nil, fmt.Errorf("missing Endpoint or Namespace information")
	}

	url := fmt.Sprintf("https://%s/api/v1/namespaces/%s/pods", p.Endpoint, p.Namespace)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.BearerTokenService.RetrieveToken()))
	res, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received %v", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	podList := &PodList{}
	err = json.Unmarshal([]byte(body), podList)
	if err != nil {
		return nil, err
	}

	return podList, nil
}
