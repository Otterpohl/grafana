package schemaloader

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/schema"
	"github.com/grafana/grafana/pkg/schema/load"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/setting"
)

func init() {
	remotecache.Register(&RenderUser{})
	registry.Register(&registry.Descriptor{
		Name:     ServiceName,
		Instance: &SchemaLoaderService{},
	})
}

const ServiceName = "SchemaLoader"

var baseLoadPath load.BaseLoadPaths = load.BaseLoadPaths{
	BaseCueFS:       grafana.CoreSchema,
	DistPluginCueFS: grafana.PluginSchema,
}

type RenderUser struct {
	OrgID   int64
	UserID  int64
	OrgRole string
}

type SchemaLoaderService struct {
	log        log.Logger
	DashFamily schema.VersionedCueSchema
	Cfg        *setting.Cfg `inject:""`
}

func (rs *SchemaLoaderService) Init() error {
	rs.log = log.New("schemaloader")
	var err error
	rs.DashFamily, err = load.BaseDashboardFamily(baseLoadPath)

	if err != nil {
		return fmt.Errorf("failed to load dashboard cue schema from path %q: %w", baseLoadPath, err)
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<I passed", len(schema.AsArray(rs.DashFamily)))
	// // ensure ImagesDir exists
	// err := os.MkdirAll(rs.Cfg.ImagesDir, 0700)
	// if err != nil {
	// 	return fmt.Errorf("failed to create images directory %q: %w", rs.Cfg.ImagesDir, err)
	// }

	// set value used for domain attribute of renderKey cookie
	// switch {
	// case rs.Cfg.RendererUrl != "":
	// RendererCallbackUrl has already been passed, it won't generate an error.
	// u, _ := url.Parse(rs.Cfg.RendererCallbackUrl)
	// 	rs.domain = u.Hostname()
	// case rs.Cfg.HTTPAddr != setting.DefaultHTTPAddr:
	// 	rs.domain = rs.Cfg.HTTPAddr
	// default:
	// 	rs.domain = "localhost"
	// }

	return nil
}

func (rs *SchemaLoaderService) DashboardTrimDefaults(input simplejson.Json) (simplejson.Json, error) {
	val, _ := input.Map()
	val = removeNils(val)
	data, _ := json.Marshal(val)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", len(schema.AsArray(rs.DashFamily)))

	dsSchema, err := schema.SearchAndValidate(rs.DashFamily, data)
	if err != nil {
		return input, err
	}
	spew.Dump(dsSchema)
	// _, err = dsSchema.TrimDefaults(schema.Resource{Value: data})
	// if err != nil {
	// 	return input, err
	// }

	output, err := simplejson.NewJson(data)
	if err != nil {
		return input, err
	}
	return *output, nil
}

func removeNils(initialMap map[string]interface{}) map[string]interface{} {
	withoutNils := map[string]interface{}{}
	for key, value := range initialMap {
		_, ok := value.(map[string]interface{})
		if ok {
			value = removeNils(value.(map[string]interface{}))
			withoutNils[key] = value
			continue
		}
		_, ok = value.([]interface{})
		if ok {
			value = removeNilArray(value.([]interface{}))
			withoutNils[key] = value
			continue
		}
		if value != nil {
			if val, ok := value.(string); ok {
				if val == "" {
					continue
				}
			}
			withoutNils[key] = value
		}
	}
	return withoutNils
}

func removeNilArray(initialArray []interface{}) []interface{} {
	withoutNils := []interface{}{}
	for _, value := range initialArray {
		_, ok := value.(map[string]interface{})
		if ok {
			value = removeNils(value.(map[string]interface{}))
			withoutNils = append(withoutNils, value)
			continue
		}
		_, ok = value.([]interface{})
		if ok {
			value = removeNilArray(value.([]interface{}))
			withoutNils = append(withoutNils, value)
			continue
		}
		if value != nil {
			if val, ok := value.(string); ok {
				if val == "" {
					continue
				}
			}
			withoutNils = append(withoutNils, value)
		}
	}
	return withoutNils
}
