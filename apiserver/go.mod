module github.com/zhaozf-zhiming/suneee/apiserver

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.3.0
	github.com/spf13/viper v1.3.2
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20190415214537-1da14a5a36f2 // indirect
	golang.org/x/oauth2 v0.0.0-20170412232759-a6bd8cefa181 // indirect
	golang.org/x/time v0.0.0-20161028155119-f51c12702a4d // indirect
	google.golang.org/appengine v1.5.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20190221042446-c2654d5206da // indirect
)

replace (
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab => github.com/kubernetes/api v0.0.0-20190416052506-9eb4726e83e4 // indirect
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1 => github.com/kubernetes/apimachinery v0.0.0-20190416052411-7dcd37fca543 // indirect
	k8s.io/client-go v11.0.0+incompatible => github.com/kubernetes/client-go v11.0.0+incompatible // indirect
	k8s.io/klog v0.0.0-20190306015804-8e90cee79f82 => github.com/kubernetes/klog v0.3.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30 => github.com/kubernetes/kube-openapi v0.0.0-20190401085232-94e1e7b7574c // indirect
	k8s.io/utils v0.0.0-20190221042446-c2654d5206da => github.com/kubernetes/utils v0.0.0-20190308190857-21c4ce38f2a7 // indirect

)
