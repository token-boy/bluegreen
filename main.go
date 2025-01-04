package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"
)

type Server struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type TLS struct {
	CertResolver string `yaml:"certResolver"`
}

type Router struct {
	Rule    string `yaml:"rule"`
	TLS     TLS    `yaml:"tls"`
	Service string `yaml:"service"`
}

type Http struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

type Config struct {
	Http Http `yaml:"http"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		serviceName := query.Get("service")
		host := query.Get("host")
		port := query.Get("port")
		if serviceName == "" || host == "" || port == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required parameters (service, host, port)"))
			return
		}

		containers, err := cli.ContainerList(ctx, container.ListOptions{
			Filters: filters.NewArgs(filters.Arg("label", serviceName)),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if len(containers) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No running containers"))
			return
		}

		// Sort containers by creation time
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Created > containers[j].Created
		})

		configFilePath := "/etc/traefik/conf/" + serviceName + ".yml"
		if !fileExists(configFilePath) {
			config := Config{}

			config.Http.Routers = make(map[string]Router)
			config.Http.Routers[serviceName] = Router{
				Rule: fmt.Sprintf("Host(`%s`)", host),
				TLS: TLS{
					CertResolver: "letsencrypt",
				},
				Service: serviceName + "@file",
			}

			config.Http.Services = make(map[string]Service)
			config.Http.Services[serviceName] = Service{
				LoadBalancer{
					Servers: []Server{
						{URL: fmt.Sprintf("http://%s:%s", containers[0].NetworkSettings.Networks["main"].IPAddress, port), Weight: 100},
					},
				},
			}

			configData, err := yaml.Marshal(&config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			os.WriteFile(configFilePath, configData, 0644)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("First container deployed"))
		} else {
			configFile, err := os.ReadFile(configFilePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			var config Config
			if err := yaml.Unmarshal(configFile, &config); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			servers := []Server{}
			for i := range len(containers) {
				weight := 0
				if i == 0 {
					weight = 100
				}
				servers = append(servers, Server{
					URL:    fmt.Sprintf("http://%s:%s", containers[i].NetworkSettings.Networks["main"].IPAddress, port),
					Weight: weight,
				})
			}
			service := config.Http.Services[serviceName]
			service.LoadBalancer.Servers = servers
			config.Http.Services[serviceName] = service

			configData, err := yaml.Marshal(&config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			os.WriteFile(configFilePath, configData, 0644)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Container joined the load balancer"))
		}
	})

	fmt.Println("Server is running on http://localhost:10208")
	if err := http.ListenAndServe(":10208", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
