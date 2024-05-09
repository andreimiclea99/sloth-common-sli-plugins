package availability

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/prometheus/targets/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
sum(count_over_time((sum by (service) (clamp_min(up{ {{.filter}} }, 0)) > 0) [{{"{{ .window }}"}}:1m]))
/
sum(count_over_time((sum by (service) (clamp_min(up{ {{.filter}} }, 1)) [{{"{{ .window }}"}}:1m]))
`))

// SLIPlugin will return a query that will return the availability of Prometheus registered targets.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	var b bytes.Buffer
	data := map[string]string{
		"filter": getFilter(options),
	}
	err := queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getFilter(options map[string]string) string {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")

	return filter
}
