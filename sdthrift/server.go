package sdthrift

import (
	"crypto/tls"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
)

type ProcessorFactory func() (thrift.TProcessor, error)

type ServerOptions struct {
	Secure            bool
	CertFile, KeyFile string
	TF                thrift.TTransportFactory
	PF                thrift.TProtocolFactory
	Factory           ProcessorFactory
}

func RunServer(addr string, opts ServerOptions) error {
	if opts.Factory == nil {
		return sderr.New("nil factory")
	}
	if opts.TF == nil {
		return sderr.New("nil transport factory")
	}
	if opts.PF == nil {
		return sderr.New("nil protocol factory")
	}

	var transport thrift.TServerTransport
	var err error
	if opts.Secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return sderr.WithStack(err)
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}
	if err != nil {
		return sderr.WithStack(err)
	}
	processor, err := opts.Factory()
	if err != nil {
		return sderr.WithStack(err)
	}
	if processor == nil {
		return sderr.New("nil processor")
	}
	server := thrift.NewTSimpleServer4(processor, transport, opts.TF, opts.PF)
	sdlog.WithField("addr", addr).Info("Run thrift server")
	return sderr.WithStack(server.Serve())
}
