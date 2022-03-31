package manager

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/web_feature"
)

func (m *Manager) extensionWebFeature(configuration *Configuration) (*web_feature.WebFeature, error) {
	if configuration.Extension == nil || configuration.Extension.WebFeature == nil {
		return nil, ErrExtensionConfigurationConflict
	}

	model := caskin.CasbinModelText()

	if configuration.EntryFactory == nil {
		return nil, caskin.ErrInitializationNilEntryFactory
	}
	if configuration.MetaDB == nil {
		return nil, caskin.ErrInitializationNilMetaDB
	}

	options := &web_feature.Options{
		Caskin:        m.caskin,
		DomainFactory: m.caskin.GetOptions().GetSuperadminDomain,
		ObjectFactory: configuration.EntryFactory.NewObject,
		MetaDB:        configuration.MetaDB,
		ModelText:     model,
	}

	return web_feature.New(options)
}
