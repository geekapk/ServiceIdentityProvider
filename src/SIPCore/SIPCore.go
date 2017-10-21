package SIPCore

import (
	"sync"
	"errors"
	"SIPConfig"
)

type Context struct {
	m sync.Mutex
	Config *SIPConfig.Config
	serviceLookupTableByKey map[string]*SIPConfig.Service
	serviceLookupTableByName map[string]*SIPConfig.Service
	groupLookupTableByName map[string]*SIPConfig.Group
}

type Identity struct {
	CurrentContext *Context
	Service *SIPConfig.Service
	Groups []*SIPConfig.Group
}

func NewContext(config *SIPConfig.Config) *Context {
	c := &Context {
		Config: config,
	}

	c.initLookupTables()
	return c
}

func (c *Context) initLookupTables() {
	c.serviceLookupTableByKey = make(map[string]*SIPConfig.Service)
	c.serviceLookupTableByName = make(map[string]*SIPConfig.Service)
	c.groupLookupTableByName = make(map[string]*SIPConfig.Group)

	for _, svc := range c.Config.Services {
		c.serviceLookupTableByKey[svc.Key] = svc
		c.serviceLookupTableByName[svc.Name] = svc
	}

	for _, grp := range c.Config.Groups {
		c.groupLookupTableByName[grp.Name] = grp
	}
}

func (c *Context) GetIdentityByKey(key string) (*Identity, error) {
	c.m.Lock()
	defer c.m.Unlock()

	svc, ok := c.serviceLookupTableByKey[key]
	if !ok {
		return nil, errors.New("not found")
	}

	id := &Identity {
		CurrentContext: c,
		Service: svc,
	}
	id.initExtInfo()

	return id, nil
}

func (c *Context) GetIdentityByName(name string) (*Identity, error) {
	c.m.Lock()
	defer c.m.Unlock()

	svc, ok := c.serviceLookupTableByName[name]
	if !ok {
		return nil, errors.New("not found")
	}

	id := &Identity {
		CurrentContext: c,
		Service: svc,
	}
	id.initExtInfo()

	return id, nil
}

func (id *Identity) initExtInfo() {
	for _, grp := range id.CurrentContext.Config.Groups {
		for _, m := range grp.Members {
			if m == id.Service.Name {
				id.Groups = append(id.Groups, grp)
				break
			}
		}
	}
}

func (id *Identity) GroupIntersection(other *Identity) []*SIPConfig.Group {
	interMap := make(map[string]bool)
	inter := make([]*SIPConfig.Group, 0)

	for _, grp := range id.Groups {
		interMap[grp.Name] = true
	}

	id.CurrentContext.m.Lock()
	defer id.CurrentContext.m.Unlock()

	for _, grp := range other.Groups {
		if interMap[grp.Name] {
			inter = append(inter, id.CurrentContext.groupLookupTableByName[grp.Name])
		}
	}

	return inter
}
