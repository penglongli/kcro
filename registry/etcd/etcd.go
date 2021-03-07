package etcd

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/penglongli/kcro/registry"
	"github.com/penglongli/kcro/utils"
	"github.com/penglongli/kcro/utils/log"
)

var (
	key = "/service"
)

type Options struct {
	Endpoints   []string
	DialTimeout int64
	TlsConfig   *tls.Config

	LeaseTTL int64

	Name  string
	Ports []int
}

type etcdRegistry struct {
	o      *Options
	client *clientv3.Client

	keepalived <-chan *clientv3.LeaseKeepAliveResponse
}

func NewRegistry(options *Options) registry.Registry {
	r := &etcdRegistry{
		o: options,
	}
	r.configure()
	return r
}

func (r *etcdRegistry) configure() {
	if r.o == nil {
		log.Default().Errorf("Registry options cannot be nil.")
		return
	}
	if len(r.o.Endpoints) == 0 {
		log.Errorf("Require at least one etcd endpoint.")
		return
	}
	if r.o.DialTimeout <= 0 {
		r.o.DialTimeout = 5
	}
	if r.o.LeaseTTL <= 0 {
		r.o.LeaseTTL = 20
	}
	if r.o.TlsConfig == nil {
		r.o.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	for _, ep := range r.o.Endpoints {
		if len(ep) == 0 {
			continue
		}
		_, _, err := net.SplitHostPort(ep)
		if err != nil {
			log.Default().Error(err)
			return
		}
	}

	config := clientv3.Config{
		DialTimeout: time.Duration(r.o.DialTimeout) * time.Second,
		TLS:         r.o.TlsConfig,
		Endpoints:   r.o.Endpoints,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		log.Default().Error(err)
		return
	}
	r.client = cli
	return
}

func (r *etcdRegistry) Register() error {
	if r.client == nil {
		return errors.Errorf("Etcd client not initialized.")
	}

	lease := clientv3.NewLease(r.client)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.o.DialTimeout)*time.Second)
	leaseResp, err := lease.Grant(ctx, r.o.LeaseTTL)
	if err != nil {
		log.Error(err)
		return err
	}

	// Get local ipv4 address.
	ips, err := utils.IFaces()
	if err != nil {
		log.Error(err)
		return err
	}

	kv := clientv3.NewKV(r.client)
	for _, ip := range ips {
		for _, port := range r.o.Ports {
			ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.o.DialTimeout)*time.Second)
			_, err := kv.Put(ctx, key+"/"+r.o.Name, ip+":"+strconv.Itoa(port), clientv3.WithLease(leaseResp.ID))
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	keepalived, err := lease.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Error(err)
		return err
	}
	r.keepalived = keepalived
	return nil
}

func (r *etcdRegistry) Deregister() error {
	return nil
}
