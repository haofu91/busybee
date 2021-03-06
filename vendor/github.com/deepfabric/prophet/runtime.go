package prophet

import (
	"sort"
	"sync"
)

const (
	batchLimit = 10000
)

// Runtime runtime info
type Runtime struct {
	sync.RWMutex

	p          *defaultProphet
	containers map[uint64]*ContainerRuntime
	resources  map[uint64]*ResourceRuntime

	leaders   map[uint64]map[uint64]*ResourceRuntime // container -> resource -> ResourceRuntime
	followers map[uint64]map[uint64]*ResourceRuntime // container -> resource -> ResourceRuntime
}

func newRuntime(p *defaultProphet) *Runtime {
	return &Runtime{
		p:          p,
		containers: make(map[uint64]*ContainerRuntime),
		resources:  make(map[uint64]*ResourceRuntime),
		leaders:    make(map[uint64]map[uint64]*ResourceRuntime),
		followers:  make(map[uint64]map[uint64]*ResourceRuntime),
	}
}

func (rc *Runtime) load() {
	err := rc.p.store.LoadResources(batchLimit, func(meta Resource) {
		rc.Lock()
		defer rc.Unlock()

		rc.resources[meta.ID()] = newResourceRuntime(meta, nil)
	})
	if err != nil {
		log.Fatalf("load resources failed, errors:%+v", err)
	}

	err = rc.p.store.LoadContainers(batchLimit, func(meta Container) {
		rc.Lock()
		defer rc.Unlock()

		rc.containers[meta.ID()] = newContainerRuntime(meta)
	})
	if err != nil {
		log.Fatalf("load containers failed, errors:%+v", err)
	}
}

// Containers returns the containers, using clone
func (rc *Runtime) Containers() []*ContainerRuntime {
	rc.RLock()
	defer rc.RUnlock()

	value := make([]*ContainerRuntime, len(rc.containers), len(rc.containers))
	idx := 0
	for _, cr := range rc.containers {
		value[idx] = cr.Clone()
		idx++
	}

	return value
}

// Resources returns the resources, using clone
func (rc *Runtime) Resources() []*ResourceRuntime {
	rc.RLock()
	defer rc.RUnlock()

	value := make([]*ResourceRuntime, len(rc.resources), len(rc.resources))
	idx := 0
	for _, cr := range rc.resources {
		value[idx] = cr.Clone()
		idx++
	}

	return value
}

// Container returns a cloned value of container runtime info
func (rc *Runtime) Container(id uint64) *ContainerRuntime {
	rc.RLock()
	defer rc.RUnlock()

	return rc.getContainerWithoutLock(id)
}

// Resource returns a cloned value of resource runtime info
func (rc *Runtime) Resource(id uint64) *ResourceRuntime {
	rc.RLock()
	defer rc.RUnlock()

	return rc.getResourceWithoutLock(id)
}

// ResourceContainers returns containers that has the resource's peer
func (rc *Runtime) ResourceContainers(target *ResourceRuntime) []*ContainerRuntime {
	rc.RLock()
	defer rc.RUnlock()

	var containers []*ContainerRuntime
	for id := range target.GetContainerIDs() {
		if container := rc.getContainerWithoutLock(id); container != nil {
			containers = append(containers, container.Clone())
		}
	}
	return containers
}

// ResourceFollowerContainers returns all containers for peers exclude leader
func (rc *Runtime) ResourceFollowerContainers(res *ResourceRuntime) []*ContainerRuntime {
	rc.RLock()
	defer rc.RUnlock()

	var containers []*ContainerRuntime
	for id := range res.getFollowers() {
		if container := rc.getContainerWithoutLock(id); container != nil {
			containers = append(containers, container)
		}
	}
	return containers
}

// RandLeaderResource returns the random leader resource
func (rc *Runtime) RandLeaderResource(id uint64, kind ResourceKind) *ResourceRuntime {
	rc.RLock()
	defer rc.RUnlock()

	return randResource(rc.leaders[id], kind, rc.p.cfg.ResourceSortCompareFunc)
}

// RandFollowerResource returns the random follower resource
func (rc *Runtime) RandFollowerResource(id uint64, kind ResourceKind) *ResourceRuntime {
	rc.RLock()
	defer rc.RUnlock()

	return randResource(rc.followers[id], kind, rc.p.cfg.ResourceSortCompareFunc)
}

func (rc *Runtime) getContainerWithoutLock(id uint64) *ContainerRuntime {
	container, ok := rc.containers[id]
	if !ok {
		return nil
	}

	return container.Clone()
}

func (rc *Runtime) getResourceWithoutLock(id uint64) *ResourceRuntime {
	resource, ok := rc.resources[id]
	if !ok {
		return nil
	}

	return resource.Clone()
}

func randResource(resources map[uint64]*ResourceRuntime, kind ResourceKind, compare func(Resource, Resource) int) *ResourceRuntime {
	var values []*ResourceRuntime
	for _, res := range resources {
		values = append(values, res)
	}

	if compare != nil {
		sort.Slice(values, func(i, j int) bool {
			return compare(values[i].meta, values[j].meta) < 0
		})
	}

	for _, res := range values {
		if res.leaderPeer == nil {
			log.Fatalf("rand resource %d without leader", res.meta.ID())
		}

		if len(res.downPeers) > 0 {
			continue
		}

		if len(res.pendingPeers) > 0 {
			continue
		}

		if kind == LeaderKind && !res.meta.SupportTransferLeader() {
			continue
		}

		if kind == ReplicaKind && !res.meta.SupportRebalance() {
			continue
		}

		return res.Clone()
	}

	return nil
}
