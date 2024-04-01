package nsip

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Record struct {
	Next  plugin.Handler
	Rules []rule
}

type rule struct {
	zones    []string
	policies []policy
}

type policy struct {
	ns string
	ip string
}

func (a Record) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	querySourceIP := state.IP()
	namespace, err := findPodNamespaceByIP(querySourceIP)
	if err != nil {
		log.Errorf("error searching for namespace: %v", err)
		return plugin.NextOrFailure(a.Name(), a.Next, ctx, w, r)
	}
	log.Info(fmt.Sprintf("query from namespace: %s", namespace))

	m := new(dns.Msg)
	m.SetReply(r)

	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}

	for _, p := range a.Rules[0].policies {
		if p.ns == namespace {
			rr.A = net.ParseIP(p.ip)
			m.Answer = append(m.Answer, rr)
		}
	}

	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func (a Record) Name() string { return "any" }

func findPodNamespaceByIP(ip string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get in-cluster configuration: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve list of pods: %w", err)
	}

	for _, pod := range pods.Items {
		if pod.Status.PodIP == ip {
			return pod.Namespace, nil
		}
	}

	return "", fmt.Errorf("pod with IP %s not found", ip)
}
