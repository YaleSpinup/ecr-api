/*
Copyright © 2020 Yale University

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {
	api := s.router.PathPrefix("/v1/ecr").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// ECS Image Registry Repository handlers
	api.HandleFunc("/{account}/repositories", s.RepositoriesListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/repositories/{group}", s.RepositoriesCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/{account}/repositories/{group}", s.RepositoriesListHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/repositories/{group}/{name}", s.RepositoriesShowHandler).Methods(http.MethodGet)
	api.HandleFunc("/{account}/repositories/{group}/{name}", s.RepositoriesUpdateHandler).Methods(http.MethodPut)
	api.HandleFunc("/{account}/repositories/{group}/{name}", s.RepositoriesDeleteHandler).Methods(http.MethodDelete)

	api.HandleFunc("/{account}/repositories/{group}/{name}/images", s.RepositoriesImageListHandler).Methods(http.MethodGet)
}
