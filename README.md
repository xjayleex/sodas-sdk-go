# sodas-sdk-go
sodas-sdk-go는 크게 두가지 API를 제공하고 있습니다.
1. sodas+ rest api에 접근하기 위한 go언어 클라이언트
2. sodas+ 컨테이너 런타임 os 환경변수 해석을 위한 API
    ***runtime.Retreive(prop runtime.AppProperties)***

## 1. SODAS rest api 클라이언트
[sodas rest api](https://api.221.154.134.31.traefik.me:10012/)에 대한 go언어 래퍼 클라이언트입니다. 레포지토리 패턴 등 여러 디자인 패턴을 적용하여, 각각의 API 요청에 필요한 boilerplates 코드들을 제거한 것이 컨트르뷰션이라고 할 수 있습니다. 예를 들어 다음과 같은 메소드 체인으로 Login API(/api/gateway/authentication/user/login)를 호출할 수 있습니다.
```go
clientset, err := sodas.NewForConfig(&config)
...
result, err = clientset.Gateway().Auth().User().Login(context.Background(), gateway.LoginBody{
		Id:       "Johndoe",
		Password: "johndoe",
		Offline:  false,
})
```
이 API는 23'.09.18 ETRI 간 회의를 통해 확인한 바, deprecated 예정입니다. 해당 래퍼 클라이언트는 이 deprecated 예정인 API를 대상으로 구성된 프레임워크입니다. 따라서 [변경된 API](http://datalake-api.221.154.134.31.traefik.me:10017/swagger-ui/swagger-ui/index.html)에 대한 go wrapper가 필요하다면, 수정이 필요한 상태입니다.  
API 레퍼런스 중에서 필요한 API Path에 대한 정의 그리고 GET Call Header, POST Call Body와 같은 파라미트 스펙 정의 등을 수정하시면 됩니다. 대게는 인증과 Object Storage에 관한 API 래퍼만 필요하실거라 생각이 됩니다.

## 2. SODAS+ 컨테이너 런타임 OS 환경변수 해석을 위한 API
SODAS+ 워크플로우 도구의 알고리즘 작성이나, 하베스팅 환경에서 플랫폼으로부터 들어오는 입력 파라미터, 혹은 템플릿은 컨테이너 런타임 OS 환경변수의 json 스트링으로 입력됩니다. 예를 들어, 프로그램에서 사용자가 API 인증을 위해서 **계정 정보(id,pw)** 를 통해, AccessToken을 발급받아야만 하는 상황이라고 해봅니다. id:string, pw:string을 입력받아야 하는 상황인 것이며, 이는 사용자 입장에서는 워크플로우 실행 시, 사용자가 입력하게 됩니다. 입력된 두 key:value 쌍이 플랫폼에 의해 컨테이너 런타임 OS 환경변수로 주입되는 것입니다. 이는 json string 형태로 되는데, 이를 파싱해야 하는 것이고요. Golang으로 SODAS+ 알고리즘 작성 시, 이 파싱을 도와주는 것이 sodas.runtime 패키지입니다(sodas-sdk-go).
```go
import "sodas-sdk-go/runtime"
// runtime/property/property.go의 AppProperties 인터페이스 구현체
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

func main() {
    prop := HelloAppProperties{}
    err := runtime.Retrieve(&prop)
    // use {prop}
    ...
}
```
