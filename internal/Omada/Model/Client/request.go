package Client

import (
	"omada_exporter_go/internal/Log"
	"omada_exporter_go/internal/Omada/HttpClient/ApiClient"
)

// Get fetches all clients from the Omada OpenAPI, paginating automatically.
func Get() (*[]Client, error) {
	Log.Debug("Fetching client data")
	client := ApiClient.GetInstance()

	result, err := ApiClient.Get[Client](client, path_OpenApiClients, nil, nil, true)
	if err != nil {
		return nil, Log.Error(err, "Failed to get client data")
	}

	Log.Debug("Fetched %d clients", len(*result))
	return result, nil
}
