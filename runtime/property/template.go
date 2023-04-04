package property

// Do not use this property struct directly. This is an example code for
// representing how to implement `AppProperties` interface.
type HelloAppProperties struct {
	Input Input `sodas_prop:"input"`
	Head  int   `sodas_prop:"head"`
}

func (p *HelloAppProperties) RootFieldTag() string {
	return "sodas_prop"
}

type Input struct {
	Type         string `json:"type"`
	UserName     string `json:"user_name"`
	BaseUrl      string `json:"base_url"`
	RefreshToken string `json:"refresh_token"`
	Endpoint     string `json:"end_point"`
	ObjectName   string `json:"object_name"`
}

type Output struct {
	Type         string `json:"type"`
	UserName     string `json:"user_name"`
	BaseUrl      string `json:"base_url"`
	RefreshToken string `json:"refresh_token"`
	Endpoint     string `json:"end_point"`
	ObjectName   string `json:"object_name"`
}
