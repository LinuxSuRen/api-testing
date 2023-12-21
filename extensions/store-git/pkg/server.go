/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pkg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
)

type gitClient struct {
	writer io.Writer
	remote.UnimplementedLoaderServer
}

// NewRemoteServer returns a new remote server
func NewRemoteServer() remote.LoaderServer {
	return &gitClient{writer: os.Stdout}
}
func (s *gitClient) loadCache(ctx context.Context) (opt *gitOptions, err error) {
	if opt, err = s.getClient(ctx); err != nil {
		return
	}

	remoteName := "origin"
	branch := "master"
	repoAddr := opt.cloneOptions.URL
	configDir := opt.cache
	log.Println("load cache from", repoAddr)

	if ok, _ := util.PathExists(configDir); ok {
		var repo *git.Repository
		if repo, err = git.PlainOpen(configDir); err == nil {
			var wd *git.Worktree

			if wd, err = repo.Worktree(); err == nil {
				if err = repo.FetchContext(ctx, opt.fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
					err = fmt.Errorf("failed to fetch '%s', error: %v", remoteName, err)
					return
				}

				if err = wd.PullContext(ctx, &git.PullOptions{
					RemoteName:    remoteName,
					ReferenceName: plumbing.NewBranchReferenceName(branch),
					Force:         true,
					Auth:          opt.cloneOptions.Auth,
				}); err != nil && err != git.NoErrAlreadyUpToDate {
					err = fmt.Errorf("failed to pull git repository '%s', error: %v", repo, err)
					return
				}
				err = nil
			}
		} else {
			err = fmt.Errorf("failed to open git local repository, error: %v", err)
		}
	} else {
		if _, err = git.PlainCloneContext(ctx, configDir, false, opt.cloneOptions); err != nil {
			err = fmt.Errorf("failed to clone git repository '%s' into '%s', error: %v", repoAddr, configDir, err)
		}
	}
	return
}
func (s *gitClient) pushCache(ctx context.Context) (err error) {
	var opt *gitOptions
	if opt, err = s.getClient(ctx); err != nil {
		return
	}

	configDir := opt.cache

	var repo *git.Repository
	if repo, err = git.PlainOpen(configDir); err == nil {
		var wd *git.Worktree

		if wd, err = repo.Worktree(); err == nil {
			if _, err = wd.Add("."); err != nil {
				return
			}

			if _, err = wd.Commit(`auto commit by api-testing

See also https://github.com/LinuxSuRen/api-testing
`, &git.CommitOptions{
				Author: &object.Signature{
					Name:  opt.name,
					Email: opt.email,
					When:  time.Now(),
				},
			}); err == nil {
				err = repo.Push(opt.pushOptions)
			}
		}
	}
	return
}
func (s *gitClient) ListTestSuite(ctx context.Context, _ *server.Empty) (reply *remote.TestSuites, err error) {
	reply = &remote.TestSuites{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	var suites []testing.TestSuite
	if suites, err = loader.ListTestSuite(); err == nil {
		for _, item := range suites {
			reply.Data = append(reply.Data, remote.ConvertToGRPCTestSuite(&item))
		}
	}
	return
}
func (s *gitClient) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.CreateSuite(testSuite.Name, testSuite.Api); err == nil {
		s.pushCache(ctx)
	}
	return
}
func (s *gitClient) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	var normalSuite testing.TestSuite
	if normalSuite, err = loader.GetTestSuite(suite.Name, true); err == nil {
		reply = remote.ConvertToGRPCTestSuite(&normalSuite)
	}
	return
}
func (s *gitClient) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.UpdateSuite(*remote.ConvertToNormalTestSuite(suite)); err == nil {
		err = s.pushCache(ctx)
	}
	return
}
func (s *gitClient) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.DeleteSuite(suite.Name); err == nil {
		err = s.pushCache(ctx)
	}
	return
}
func (s *gitClient) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *server.TestCases, err error) {
	if suite, err = s.GetTestSuite(ctx, suite); err == nil {
		result = &server.TestCases{
			Data: suite.Items,
		}
	}
	return
}
func (s *gitClient) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.CreateTestCase(testcase.SuiteName, remote.ConvertToNormalTestCase(testcase)); err == nil {
		err = s.pushCache(ctx)
	}
	return
}
func (s *gitClient) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	result = &server.TestCase{}

	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{Name: testcase.SuiteName}); err == nil {
		for _, item := range suite.Items {
			if item.Name == testcase.Name {
				result = item
				break
			}
		}
	}
	return
}
func (s *gitClient) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.UpdateTestCase(testcase.SuiteName, remote.ConvertToNormalTestCase(testcase)); err == nil {
		err = s.pushCache(ctx)
	}
	return
}
func (s *gitClient) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var loader testing.Writer
	if loader, err = s.newLoader(ctx); err != nil {
		return
	}

	if err = loader.DeleteTestCase(testcase.SuiteName, testcase.Name); err == nil {
		err = s.pushCache(ctx)
	}
	return
}
func (s *gitClient) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	_, clientErr := s.ListTestSuite(ctx, in)
	reply = &server.ExtensionStatus{
		Ready:   clientErr == nil,
		Message: util.OKOrErrorMessage(clientErr),
		Version: version.GetVersion(),
	}

	// no git repository allows to write files without authentication
	store := remote.GetStoreFromContext(ctx)
	if store != nil {
		reply.ReadOnly = store.Username == "" || store.Password == ""
	}
	return
}
func (s *gitClient) PProf(ctx context.Context, in *server.PProfRequest) (data *server.PProfData, err error) {
	log.Println("pprof", in.Name)

	data = &server.PProfData{
		Data: extension.LoadPProf(in.Name),
	}
	return
}
func (s *gitClient) getClient(ctx context.Context) (opt *gitOptions, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to git server")
	} else {
		auth := &http.BasicAuth{
			Username: store.Username,
			Password: store.Password,
		}

		insecure := store.Properties["insecure"] == "true"

		opt = &gitOptions{
			cache:      filepath.Join(os.TempDir(), store.Name),
			targetPath: store.Properties["targetpath"],
			name:       store.Properties["name"],
			email:      store.Properties["email"],
			cloneOptions: &git.CloneOptions{
				URL:             store.URL,
				Progress:        s.writer,
				InsecureSkipTLS: insecure,
				Auth:            auth,
			},
			pushOptions: &git.PushOptions{
				Progress:        s.writer,
				InsecureSkipTLS: insecure,
				Auth:            auth,
			},
			fetchOptions: &git.FetchOptions{
				Progress:        s.writer,
				InsecureSkipTLS: insecure,
				Auth:            auth,
			},
		}

		if opt.name == "" {
			opt.name = "LinuxSuRen"
		}
		if opt.email == "" {
			opt.email = "LinuxSuRen@users.noreply.github.com"
		}
	}
	return
}

func (s *gitClient) newLoader(ctx context.Context) (loader testing.Writer, err error) {
	var opt *gitOptions
	if opt, err = s.loadCache(ctx); err == nil {
		loader, err = opt.newLoader()
	}
	return
}

type gitOptions struct {
	cache        string
	targetPath   string
	name         string
	email        string
	cloneOptions *git.CloneOptions
	pushOptions  *git.PushOptions
	fetchOptions *git.FetchOptions
}

func (g *gitOptions) newLoader() (loader testing.Writer, err error) {
	if g.targetPath == "" {
		g.targetPath = "."
	}
	parentDir := path.Join(g.cache, g.targetPath)
	if err = os.MkdirAll(parentDir, 0755); err == nil {
		loader = testing.NewFileWriter(parentDir)
		loader.Put(parentDir + "/*.yaml")
	}
	return
}
