package server

import (
	"fmt"
	"github.com/deffusion/boot-server/httpclient"
	"github.com/deffusion/boot-server/nodectrl"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"log"
)

type PeerMeta struct {
	IP   string
	TCP  uint
	HTTP uint
}

type Server struct {
	serverConf ServerConfig
	netConf    NetConfig

	containerPortStartFrom uint
	gin                    *gin.Engine
	containers             map[string]*nodectrl.Container
	metas                  map[string]*PeerMeta
	dockerCli              *nodectrl.DockerCli
}

func New() (*Server, error) {

	serverConf, err := serverConfFromFile()
	if err != nil {
		return nil, err
	}
	netConf, err := netConfFromFile()
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	dockerCli, err := nodectrl.New()
	if err != nil {
		return nil, err
	}
	s := &Server{
		serverConf: serverConf,
		netConf:    netConf,

		containerPortStartFrom: serverConf.NATStartFromPort,
		containers:             map[string]*nodectrl.Container{},
		metas:                  map[string]*PeerMeta{},
		dockerCli:              dockerCli,
	}
	rg := r.Group("/", func(ctx *gin.Context) {
		ctx.Set("server", s)
	})
	Route(rg)
	s.gin = r
	return s, nil
}

func (s *Server) Run() error {
	s.StartContainers()
	s.SetContainerMetas()
	s.ConnectContainers()
	return s.gin.Run(fmt.Sprintf(":%d", s.serverConf.HTTP))
}

func (s *Server) StartContainers() {
	for i := 0; i < s.netConf.Size; i++ {
		tcpPort := s.containerPortStartFrom + uint(2*i)
		httpPort := tcpPort + 1
		var containerName = fmt.Sprintf("netfusion-%d", tcpPort)
		log.Printf("%s tcp:%d, http:%d", containerName, tcpPort, httpPort)
		portBindings := nat.PortMap{
			"3721/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprint(tcpPort),
				},
			},
			"3824/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprint(httpPort),
				},
			},
		}
		c, err := s.dockerCli.NewContainer(
			&container.Config{
				Image: "netfusion",
				ExposedPorts: nat.PortSet{
					"3721/tcp": {},
					"3824/tcp": {},
				},
			},
			&container.HostConfig{
				PortBindings: portBindings,
			},
			containerName,
		)
		if err != nil {
			log.Fatal(err)
		}
		s.containers[containerName] = &c
		s.metas[containerName] = &PeerMeta{IP: s.serverConf.IP, TCP: tcpPort, HTTP: httpPort}
	}
	for name, c := range s.containers {
		c.Start()
		log.Println(name, "started")
	}
}

func (s *Server) SetContainerMetas() {
	for _, meta := range s.metas {
		url := fmt.Sprintf("http://%s:%d/meta", meta.IP, meta.HTTP)
		httpclient.Post(url, *meta)
	}
}

func (s *Server) ConnectContainers() {
	n := s.netConf.NPeerToConnect
	for _, meta := range s.metas {
		metas := make([]PeerMeta, 0, n)
		toPrint := make([]string, 0, n)
		for _, m := range s.metas {
			if len(metas) == cap(metas) {
				break
			}
			if *m == *meta {
				continue
			}
			metas = append(metas, *m)
			toPrint = append(toPrint, fmt.Sprintf(":%d", m.TCP))
		}
		url := fmt.Sprintf("http://%s:%d/peer", meta.IP, meta.HTTP)
		log.Println(fmt.Sprintf(":%d", meta.TCP), "connect to:", toPrint)
		_, err := httpclient.Post(url, metas)
		if err != nil {
			log.Println("post err:", err)
		}
	}
}

func (s *Server) StopContainers() {
	for name, n := range s.containers {
		n.Stop()
		log.Println(name, "stopped")
	}
}
