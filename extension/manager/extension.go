package manager

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/domain_creator"
	"github.com/awatercolorpen/caskin/extension/web_feature"
)

func (m *Manager) extensionDomainCreator(configuration *Configuration) (*domain_creator.Factory, error) {
	if configuration.Extension == nil || configuration.Extension.DomainCreator == nil {
		return nil, nil
	}

	if configuration.DB == nil {
		return nil, ErrInitializationNilDB
	}
	if configuration.EntryFactory == nil {
		return nil, caskin.ErrInitializationNilEntryFactory
	}
	if configuration.DomainCreator != nil {
		return nil, ErrExtensionConfigurationConflict
	}

	return domain_creator.NewFactory(configuration.DB, configuration.EntryFactory)
}

func (m *Manager) extensionWebFeature(configuration *Configuration) (*web_feature.WebFeature, error) {
	if configuration.Extension == nil || configuration.Extension.WebFeature == nil {
		return nil, nil
	}
	if configuration.SuperadminDisable {
		return nil, ErrExtensionConfigurationConflict
	}

	model, err := caskin.CasbinModelText()
	if err != nil {
		return nil, err
	}

	if configuration.SuperadminDisable {
		return nil, ErrExtensionConfigurationConflict
	}
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
