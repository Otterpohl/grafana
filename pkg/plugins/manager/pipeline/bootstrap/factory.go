package bootstrap

import (
	"fmt"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/log"
	"github.com/grafana/grafana/pkg/plugins/manager/loader/assetpath"
)

type pluginFactoryFunc func(p plugins.FoundPlugin, pluginClass plugins.Class, sig plugins.Signature) (*plugins.Plugin, error)

// DefaultPluginFactory is the default plugin factory used by the Construct step of the Bootstrap stage.
//
// It creates the plugin using plugin information found during the Discovery stage and makes use of the assetPath
// service to set the plugin's BaseURL, Module, Logos and Screenshots fields.
type DefaultPluginFactory struct {
	assetPath *assetpath.Service
}

// NewDefaultPluginFactory returns a new DefaultPluginFactory.
func NewDefaultPluginFactory(assetPath *assetpath.Service) *DefaultPluginFactory {
	return &DefaultPluginFactory{assetPath: assetPath}
}

func (f *DefaultPluginFactory) createPlugin(p plugins.FoundPlugin, class plugins.Class,
	sig plugins.Signature) (*plugins.Plugin, error) {
	baseURL, err := f.assetPath.Base(p.JSONData, class, p.FS.Base())
	if err != nil {
		return nil, fmt.Errorf("base url: %w", err)
	}
	moduleURL, err := f.assetPath.Module(p.JSONData, class, p.FS.Base())
	if err != nil {
		return nil, fmt.Errorf("module url: %w", err)
	}

	plugin := &plugins.Plugin{
		JSONData:      p.JSONData,
		FS:            p.FS,
		BaseURL:       baseURL,
		Module:        moduleURL,
		Class:         class,
		Signature:     sig.Status,
		SignatureType: sig.Type,
		SignatureOrg:  sig.SigningOrg,
	}
	plugin.SetLogger(log.New(fmt.Sprintf("plugin.%s", plugin.ID)))

	if err = setImages(plugin, f.assetPath); err != nil {
		return nil, err
	}

	return plugin, nil
}

func setImages(p *plugins.Plugin, assetPath *assetpath.Service) error {
	var err error
	for _, dst := range []*string{&p.Info.Logos.Small, &p.Info.Logos.Large} {
		*dst, err = assetPath.RelativeURL(p, *dst, defaultLogoPath(p.Type))
		if err != nil {
			return fmt.Errorf("logo: %w", err)
		}
	}
	for i := 0; i < len(p.Info.Screenshots); i++ {
		screenshot := &p.Info.Screenshots[i]
		screenshot.Path, err = assetPath.RelativeURL(p, screenshot.Path, "")
		if err != nil {
			return fmt.Errorf("screenshot %d relative url: %w", i, err)
		}
	}
	return nil
}

func defaultLogoPath(pluginType plugins.Type) string {
	return fmt.Sprintf("public/img/icn-%s.svg", string(pluginType))
}
