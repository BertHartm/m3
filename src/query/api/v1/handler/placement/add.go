// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package placement

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/m3db/m3/src/query/api/v1/handler"
	"github.com/m3db/m3/src/query/generated/proto/admin"
	"github.com/m3db/m3/src/query/util/logging"
	"github.com/m3db/m3/src/x/net/http"
	"github.com/m3db/m3cluster/placement"

	"github.com/gogo/protobuf/jsonpb"
	"go.uber.org/zap"
)

const (
	// AddHTTPMethod is the HTTP method used with this resource.
	AddHTTPMethod = http.MethodPost
)

var (
	// DeprecatedM3DBAddURL is the old url for the placement add handler, maintained for
	// backwards compatibility.
	DeprecatedM3DBAddURL = path.Join(handler.RoutePrefixV1, PlacementPathName)

	// M3DBAddURL is the url for the placement add handler (with the POST method)
	// for the M3DB service.
	M3DBAddURL = path.Join(handler.RoutePrefixV1, M3DBServicePlacementPathName)

	// M3AggAddURL is the url for the placement add handler (with the POST method)
	// for the M3Agg service.
	M3AggAddURL = path.Join(handler.RoutePrefixV1, M3AggServicePlacementPathName)

	// M3CoordinatorAddURL is the url for the placement add handler (with the POST method)
	// for the M3Coordinator service.
	M3CoordinatorAddURL = path.Join(handler.RoutePrefixV1, M3CoordinatorServicePlacementPathName)
)

// AddHandler is the handler for placement adds.
type AddHandler Handler

type unsafeAddError struct {
	hosts string
}

func (e unsafeAddError) Error() string {
	return fmt.Sprintf("instances [%s] do not have all shards available", e.hosts)
}

// NewAddHandler returns a new instance of AddHandler.
func NewAddHandler(opts HandlerOptions) *AddHandler {
	return &AddHandler{HandlerOptions: opts, nowFn: time.Now}
}

func (h *AddHandler) ServeHTTP(serviceName string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.WithContext(ctx)

	req, rErr := h.parseRequest(r)
	if rErr != nil {
		xhttp.Error(w, rErr.Inner(), rErr.Code())
		return
	}

	placement, err := h.Add(serviceName, r, req)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(unsafeAddError); ok {
			status = http.StatusBadRequest
		}
		logger.Error("unable to add placement", zap.Any("error", err))
		xhttp.Error(w, err, status)
		return
	}

	placementProto, err := placement.Proto()
	if err != nil {
		logger.Error("unable to get placement protobuf", zap.Any("error", err))
		xhttp.Error(w, err, http.StatusInternalServerError)
		return
	}

	resp := &admin.PlacementGetResponse{
		Placement: placementProto,
		Version:   int32(placement.GetVersion()),
	}

	xhttp.WriteProtoMsgJSONResponse(w, resp, logger)
}

func (h *AddHandler) parseRequest(r *http.Request) (*admin.PlacementAddRequest, *xhttp.ParseError) {
	defer r.Body.Close()
	addReq := new(admin.PlacementAddRequest)
	if err := jsonpb.Unmarshal(r.Body, addReq); err != nil {
		return nil, xhttp.NewParseError(err, http.StatusBadRequest)
	}

	return addReq, nil
}

// Add adds a placement.
func (h *AddHandler) Add(
	serviceName string,
	httpReq *http.Request,
	req *admin.PlacementAddRequest,
) (placement.Placement, error) {
	instances, err := ConvertInstancesProto(req.Instances)
	if err != nil {
		return nil, err
	}

	serviceOpts := NewServiceOptions(
		serviceName, httpReq.Header, h.M3AggServiceOptions)
	service, algo, err := ServiceWithAlgo(h.ClusterClient, serviceOpts, h.nowFn())
	if err != nil {
		return nil, err
	}

	if req.Force {
		newPlacement, _, err := service.AddInstances(instances)
		if err != nil {
			return nil, err
		}

		return newPlacement, nil
	}

	curPlacement, version, err := service.Placement()
	if err != nil {
		return nil, err
	}

	if err := validateAllAvailable(curPlacement); err != nil {
		return nil, err
	}

	newPlacement, err := algo.AddInstances(curPlacement, instances)
	if err != nil {
		return nil, err
	}

	// Ensure the placement we're updating is still the one on which we validated
	// all shards are available.
	if err := service.CheckAndSet(newPlacement, version); err != nil {
		return nil, err
	}

	return newPlacement.SetVersion(version + 1), nil
}
