package main

import (
	"bufio"
	"crypto/tls"
	"net"
	"strconv"
)

type ServerOpts struct {
	// Server Name, Default is Go Ftp Server
	Name string

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string

	// Public IP of the server
	PublicIp string

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int

	// A logger implementation, if nil the StdLogger is used
	Logger Logger
}
type Server struct {
	*ServerOpts
	listenTo  string
	logger    Logger
	listener  net.Listener
	tlsConfig *tls.Config
	feats     string
}

func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
	}
	if opts.Port == 0 {
		newOpts.Port = 9000
	} else {
		newOpts.Port = opts.Port
	}
	if opts.Name == "" {
		newOpts.Name = "Monitor Server"
	} else {
		newOpts.Name = opts.Name
	}

	newOpts.Logger = &StdLogger{}
	if opts.Logger != nil {
		newOpts.Logger = opts.Logger
	}

	newOpts.PublicIp = opts.PublicIp

	return &newOpts
}

func NewServer(opts *ServerOpts) *Server {
	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = net.JoinHostPort(opts.Hostname, strconv.Itoa(opts.Port))
	s.logger = opts.Logger
	return s
}

func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", server.listenTo)

	if err != nil {
		return err
	}

	sessionID := server.Name
	server.logger.Printf(sessionID, "%s listening on %d", server.Name, server.Port)

	return server.Serve(listener)
}

func (server *Server) Serve(l net.Listener) error {
	server.listener = l
	sessionID := ""
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			server.logger.Printf(sessionID, "listening error: %v", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		} else {
			server.logger.Printf(sessionID, "accept new connection", tcpConn.RemoteAddr())
		}

		dconn := server.newConn(tcpConn)
		go dconn.Serve()
	}
}

func (server *Server) newConn(tcpConn net.Conn) *Conn {
	c := new(Conn)
	c.sessionID = tcpConn.RemoteAddr().String()
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.server = server
	c.logger = server.logger

	return c
}
