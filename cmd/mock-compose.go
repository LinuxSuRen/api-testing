package cmd

import (
    "github.com/linuxsuren/api-testing/pkg/mock"
    "github.com/spf13/cobra"
    "os"
    "os/signal"
    "syscall"
)

func createMockComposeCmd() (c *cobra.Command) {
    c = &cobra.Command{
        Use:   "mock-compose",
        Short: "Mock a server",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) (err error) {
            reader := mock.NewLocalFileReader(args[0])

            var server *mock.Server
            if server, err = reader.Parse(); err != nil {
                return
            }

            var subServers []mock.DynamicServer
            for _, proxy := range server.Proxies {
                subProxy := &mock.Server{
                    Proxies: []mock.Proxy{proxy},
                }

                subReader := mock.NewObjectReader(subProxy)
                subServer := mock.NewInMemoryServer(c.Context(), proxy.Port)
                if err = subServer.Start(subReader, proxy.Prefix); err != nil {
                    return
                }
                subServers = append(subServers, subServer)
            }

            clean := make(chan os.Signal, 1)
            signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
            select {
            case <-c.Context().Done():
            case <-clean:
            }
            for _, server := range subServers {
                server.Stop()
            }
            return
        },
    }
    return
}
