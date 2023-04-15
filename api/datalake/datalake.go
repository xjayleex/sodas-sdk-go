package datalake

type ListCredentialParams struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
type ListCredentialResult struct {
	Total   int `json:"total"`
	Results []struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
	} `json:"results"`
}

type CreateCredentialResult struct {
	Result string `json:"result"`
}

type RemoveCredentialResult struct {
	AccessKey string `json:"accessKey"`
}
