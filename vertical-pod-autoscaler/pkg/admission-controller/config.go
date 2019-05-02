package main

import (
 "crypto/tls"
 "crypto/x509"
 "fmt"
 "time"
 "k8s.io/api/admissionregistration/v1beta1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/client-go/kubernetes"
 "k8s.io/client-go/rest"
 "github.com/golang/glog"
)

const (
 webhookConfigName = "vpa-webhook-config"
)

func getClient() *kubernetes.Clientset {
 _logClusterCodePath()
 defer _logClusterCodePath()
 config, err := rest.InClusterConfig()
 if err != nil {
  glog.Fatal(err)
 }
 clientset, err := kubernetes.NewForConfig(config)
 if err != nil {
  glog.Fatal(err)
 }
 return clientset
}
func getAPIServerCert(clientset *kubernetes.Clientset) []byte {
 _logClusterCodePath()
 defer _logClusterCodePath()
 c, err := clientset.CoreV1().ConfigMaps("kube-system").Get("extension-apiserver-authentication", metav1.GetOptions{})
 if err != nil {
  glog.Fatal(err)
 }
 pem, ok := c.Data["requestheader-client-ca-file"]
 if !ok {
  glog.Fatalf(fmt.Sprintf("cannot find the ca.crt in the configmap, configMap.Data is %#v", c.Data))
 }
 glog.V(4).Info("client-ca-file=", pem)
 return []byte(pem)
}
func configTLS(clientset *kubernetes.Clientset, serverCert, serverKey []byte) *tls.Config {
 _logClusterCodePath()
 defer _logClusterCodePath()
 cert := getAPIServerCert(clientset)
 apiserverCA := x509.NewCertPool()
 apiserverCA.AppendCertsFromPEM(cert)
 sCert, err := tls.X509KeyPair(serverCert, serverKey)
 if err != nil {
  glog.Fatal(err)
 }
 return &tls.Config{Certificates: []tls.Certificate{sCert}, ClientCAs: apiserverCA, ClientAuth: tls.NoClientCert}
}
func selfRegistration(clientset *kubernetes.Clientset, caCert []byte) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 time.Sleep(10 * time.Second)
 client := clientset.AdmissionregistrationV1beta1().MutatingWebhookConfigurations()
 _, err := client.Get(webhookConfigName, metav1.GetOptions{})
 if err == nil {
  if err2 := client.Delete(webhookConfigName, nil); err2 != nil {
   glog.Fatal(err2)
  }
 }
 webhookConfig := &v1beta1.MutatingWebhookConfiguration{ObjectMeta: metav1.ObjectMeta{Name: webhookConfigName}, Webhooks: []v1beta1.Webhook{{Name: "vpa.k8s.io", Rules: []v1beta1.RuleWithOperations{{Operations: []v1beta1.OperationType{v1beta1.Create}, Rule: v1beta1.Rule{APIGroups: []string{""}, APIVersions: []string{"v1"}, Resources: []string{"pods"}}}, {Operations: []v1beta1.OperationType{v1beta1.Create, v1beta1.Update}, Rule: v1beta1.Rule{APIGroups: []string{"autoscaling.k8s.io"}, APIVersions: []string{"v1beta1"}, Resources: []string{"verticalpodautoscalers"}}}}, ClientConfig: v1beta1.WebhookClientConfig{Service: &v1beta1.ServiceReference{Namespace: "kube-system", Name: "vpa-webhook"}, CABundle: caCert}}}}
 if _, err := client.Create(webhookConfig); err != nil {
  glog.Fatal(err)
 } else {
  glog.V(3).Info("Self registration as MutatingWebhook succeeded.")
 }
}
