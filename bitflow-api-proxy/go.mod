module github.com/bitflow-stream/bitflow-k8s-operator/bitflow-api-proxy

go 1.14

require (
	github.com/antongulenko/golib v0.0.25
	github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller v0.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.5.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	k8s.io/api v0.0.0-20191115135540-bbc9463b57e5
	k8s.io/apimachinery v0.0.0-20191115015347-3c7067801da2
	k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	sigs.k8s.io/controller-runtime v0.1.12
)

//k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
replace github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller v0.0.0 => ../bitflow-controller

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

replace (
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)

// Fix issue with the Bitbucket repository being unreachable (404)
replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
