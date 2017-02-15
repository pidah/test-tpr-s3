package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	client *http.Client
	pool   *x509.CertPool
)

var Logger = logrus.New()

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func check(e error) {
	if e != nil {
		Logger.Panic(e)
	}
}

func init() {
	Logger.Level = logrus.InfoLevel
	Logger.Formatter = &logrus.JSONFormatter{}
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)

	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool}}}
}

func main() {

	bearerToken, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	check(err)

	r := RandStringBytes(10)
	tprName := "test-" + r
	Logger.Info(tprName)

	url := "https://kubernetes:443/apis/extensions/v1beta1/namespaces/default/thirdpartyresources"
	Logger.Info("kubernetes thirdpartyresource endpoint: ", url)

	var jsonStr = []byte(`{"apiVersion": "extensions/v1beta1","kind": "ThirdPartyResource","description": "Experimental ThirdPartyResource","metadata": {"name": "dummy-test.prsn.io","labels": {"type": "ThirdPartyResource"}},"versions": [{"name": "v1"}]}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+string(bearerToken))
	resp, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
	}
	Logger.Info(resp.Status)

	svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-1"), HTTPClient: client, DisableSSL: aws.Bool(true)})

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("kubernetes-bitesize-pidah-a"),
		Key:    aws.String("test-tpr-s3"),
	})
	if err != nil {
		Logger.Error(err)
	}
	Logger.Info(result)
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	responseStr := buf.String()
	Logger.Info(responseStr)
	if strings.TrimSpace(string(responseStr)) == "testing" {
		Logger.Info("working")
	}

	request, err := http.NewRequest("DELETE", url+"/dummy-test.prsn.io", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+string(bearerToken))
	respDelete, err := client.Do(request)
	if err != nil {
		Logger.Error(err)
	}
	Logger.Info("deleted thirdparty resource", respDelete.Status)
}
