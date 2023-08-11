package version

import (
	"net/http"

	http_util "github.com/udmire/observability-operator/pkg/utils/http"
)

type BuildInfoResponse struct {
	Status    string    `json:"status"`
	BuildInfo BuildInfo `json:"data"`
}

type BuildInfo struct {
	Application string `json:"application"`
	Version     string `json:"version"`
	Revision    string `json:"revision"`
	Branch      string `json:"branch"`
	GoVersion   string `json:"goVersion"`
}

func BuildInfoHandler(application string, features interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := BuildInfoResponse{
			Status: "success",
			BuildInfo: BuildInfo{
				Application: application,
				Version:     Version,
				Revision:    Revision,
				Branch:      Branch,
				GoVersion:   GoVersion,
			},
		}

		http_util.WriteJSONResponse(w, response)
	})
}
