module github.com/openshift/oauth-server

go 1.16

require (
	github.com/RangelReale/osin v0.0.0
	github.com/RangelReale/osincli v0.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/gophercloud/gophercloud v0.1.0
	github.com/gorilla/context v0.0.0-20190627024605-8559d4a6b87e // indirect
	github.com/gorilla/securecookie v0.0.0-20190707033817-86450627d8e6 // indirect
	github.com/gorilla/sessions v0.0.0-20171008214740-a3acf13e802c
	github.com/openshift/api v0.0.0-20210915110300-3cd8091317c4
	github.com/openshift/build-machinery-go v0.0.0-20210806203541-4ea9b6da3a37
	github.com/openshift/client-go v0.0.0-20210916133943-9acee1a0fb83
	github.com/openshift/library-go v0.0.0-20210414082648-6e767630a0dc
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/text v0.3.6
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/square/go-jose.v2 v2.2.2
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/apiserver v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/component-base v0.22.1
	k8s.io/klog/v2 v2.9.0
)

replace (
	github.com/RangelReale/osin => github.com/openshift/osin v1.0.1-0.20180202150137-2dc1b4316769
	github.com/RangelReale/osincli => github.com/openshift/osincli v0.0.0-20160924135400-fababb0555f2
	github.com/openshift/api => github.com/stlaz/api v0.0.0-20210923111930-9bf0ff72d1a0
	// k8s.io/apiserver => github.com/openshift/kubernetes-apiserver v0.0.0-20210416115049-61c04108f2c7 // points to openshift-apiserver-4.8-kubernetes-1.21.0
	github.com/openshift/library-go => github.com/stlaz/library-go v0.0.0-20210922080348-60f745231bd6
)
