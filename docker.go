package docker

import (
	"time"
)

type APIContainers struct {
	ID         string `json:"Id"`
	Image      string
	Command    string
	Created    int64
	Status     string
	Ports      []APIPort
	SizeRw     int64
	SizeRootFs int64
}

type APIPort struct {
	PrivatePort int64
	PublicPort  int64
	Type        string
}

type Container struct {
	ID string

	Created time.Time

	Path string
	Args []string

	Config *Config
	State  State
	Image  string

	NetworkSettings *NetworkSettings

	SysInitPath    string
	ResolvConfPath string
	HostnamePath   string
	HostsPath      string

	VolumesRW map[string]bool
}

type PortMapping map[string]string

type NetworkSettings struct {
	IPAddress   string
	IPPrefixLen int
	Gateway     string
	Bridge      string
	PortMapping map[string]PortMapping
}

type State struct {
	Running   bool
	Pid       int
	ExitCode  int
	StartedAt time.Time
	Ghost     bool
}

type Config struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int64 // Memory limit (in bytes)
	MemorySwap      int64 // Total memory usage (memory + swap); set `-1' to disable swap
	CpuShares       int64 // CPU shares (relative weight vs. other containers)
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	Tty             bool // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool // Open stdin
	StdinOnce       bool // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Cmd             []string
	Dns             []string
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	VolumesFrom     string
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	Privileged      bool
}

type HostConfig struct {
	Binds           []string
	ContainerIDFile string
	LxcConf         []KeyValuePair
}

type KeyValuePair struct {
	Key   string
	Value string
}

type Image struct {
	ID              string    `json:"id"`
	Parent          string    `json:"parent,omitempty"`
	Comment         string    `json:"comment,omitempty"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container,omitempty"`
	ContainerConfig Config    `json:"container_config,omitempty"`
	DockerVersion   string    `json:"docker_version,omitempty"`
	Author          string    `json:"author,omitempty"`
	Config          *Config   `json:"config,omitempty"`
	Architecture    string    `json:"architecture,omitempty"`
	Size            int64
}

type APIVersion struct {
	Version   string
	GitCommit string `json:",omitempty"`
	GoVersion string `json:",omitempty"`
}

type APIInfo struct {
	Debug              bool
	Containers         int
	Images             int
	NFd                int    `json:",omitempty"`
	NGoroutines        int    `json:",omitempty"`
	MemoryLimit        bool   `json:",omitempty"`
	SwapLimit          bool   `json:",omitempty"`
	IPv4Forwarding     bool   `json:",omitempty"`
	LXCVersion         string `json:",omitempty"`
	NEventsListener    int    `json:",omitempty"`
	KernelVersion      string `json:",omitempty"`
	IndexServerAddress string `json:",omitempty"`
}

type APIImages struct {
	Repository  string `json:",omitempty"`
	Tag         string `json:",omitempty"`
	ID          string `json:"Id"`
	Created     int64
	Size        int64
	VirtualSize int64
}

type EventErr struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Event struct {
	Status   string   `json:"status,omitempty"`
	Progress string   `json:"progress,omitempty"`
	ID       string   `json:"id,omitempty"`
	From     string   `json:"from,omitempty"`
	Time     int64    `json:"time,omitempty"`
	Error    EventErr `json:"errorDetail,omitempty"`
}
