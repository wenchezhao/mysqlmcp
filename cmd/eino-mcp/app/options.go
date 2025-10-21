package app

import (
	pflag "github.com/spf13/pflag"
)

type ServerRunOptions struct {
	listenAddr     string
	Debug          bool
	ConfigFilePath string
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{}
}

func (s *ServerRunOptions) Flags() (fs *pflag.FlagSet) {
	fs = pflag.NewFlagSet("agent", pflag.ExitOnError)
	//fs.AddGoFlagSet(pflag.CommandLine)
	fs.BoolVar(&s.Debug, "debug", false, "Enable debug mode,default is false.")
	fs.StringVar(&s.listenAddr, "listenAddr", ":8888", "Server listening addr,default is :8888.")
	fs.StringVar(&s.ConfigFilePath, "config", "./config.yaml", "Config file path,default is ./config.yaml.")
	return fs
}
