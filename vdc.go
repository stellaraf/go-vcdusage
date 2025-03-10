package vcdusage

import (
	"fmt"
	"strings"

	"github.com/destel/rill"
	"github.com/joomcode/errorx"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

func isValidStatus(status string) bool {
	valid := []int{2, 3, 4, 8}
	for _, v := range valid {
		if status == types.VAppStatuses[v] {
			return true
		}
	}
	return false
}

// VDC is a wrapper around a govcd.Vdc object with references to the corresponding Admin Org and
// vcdusage.Client.
type VDC struct {
	Obj      *govcd.Vdc
	AdminOrg *govcd.AdminOrg
	Client   *Client
}

// VDCs is a slice of vcdusage.VDC wrappers.
type VDCs []VDC

// CoreCount retrieves the used CPU MHz for a VDC and calculates the number of cores used
// by all VDCs by dividing the total used CPU MHz by the CPU speed.
//
// For example, if the speed is 3.1 GHz and the used amount is 49.6, the core count is 16.
func (vdcs VDCs) CoreCount() uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	count := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		count += vdc.CoreCount()
		return nil
	})
	return count
}

// Memory retrieves the amount of used memory to all VDCs, represented as a DataStorage type.
func (vdcs VDCs) Memory() DataStorage {
	vdcSlice := rill.FromSlice(vdcs, nil)
	mem := DataStorage(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		mem += vdc.Memory()
		return nil
	})
	return mem
}

// Memory retrieves the amount of used storage to all VDCs, represented as a DataStorage type.
func (vdcs VDCs) Storage() DataStorage {
	vdcSlice := rill.FromSlice(vdcs, nil)
	stor := DataStorage(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		stor += vdc.Storage()
		return nil
	})
	return stor
}

// VMCount retrieves the number of VMs deployed in all VDCs.
func (vdcs VDCs) VMCount() uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	count := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		count += vdc.VMCount()
		return nil
	})
	return count
}

// PoweredOnVMCount retrieves the number of powered on VMs deployed in all VDCs.
func (vdcs VDCs) PoweredOnVMCount() uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	count := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		count += vdc.PoweredOnVMCount()
		return nil
	})
	return count
}

// VMCountWithQuery retrieves the number of VMs matching all of the provided queries in all VDCs.
// If PoweredOn is false (default), VMs that are both powered on or off will be included.
func (vdcs VDCs) VMCountWithQuery(queries ...VMQuerySetter) uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	count := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		count += vdc.VMCountWithQuery(queries...)
		return nil
	})
	return count
}

// VMCountWithQuery retrieves the number of cores on VMs matching all of the provided queries in
// all VDCs. If PoweredOn is false (default), VMs that are both powered on or off will be included.
func (vdcs VDCs) VMCoreCountWithQuery(queries ...VMQuerySetter) uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	count := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		count += vdc.VMCoreCountWithQuery(queries...)
		return nil
	})
	return count
}

// Speed retrieves the max CPU speed of all VDCs in MHz. This is required for calculating core count.
func (vdcs VDCs) Speed() uint64 {
	vdcSlice := rill.FromSlice(vdcs, nil)
	speed := uint64(0)
	rill.ForEach(vdcSlice, len(vdcs), func(vdc VDC) error {
		vdcSpeed := vdc.Speed()
		if vdcSpeed > speed {
			speed = vdcSpeed
		}
		return nil
	})
	return speed
}

// Speed retrieves the CPU speed of a VDC in MHz. This is required for calculating core count.
func (vdc *VDC) Speed() uint64 {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	if avdc.AdminVdc.VCpuInMhz2 == nil {
		return 0
	}
	return uint64(*avdc.AdminVdc.VCpuInMhz2)
}

// CoreCount retrieves the used CPU MHz for a VDC and calculates the number of cores used
// by the VDC by dividing the total used CPU MHz by the CPU speed.
//
// For example, if the speed is 3.1 GHz and the used amount is 49.6, the core count is 16.
func (vdc *VDC) CoreCount() uint64 {
	speed := vdc.Speed()
	if speed == 0 {
		return 0
	}
	cores := uint64(0)
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	for _, capacity := range avdc.AdminVdc.ComputeCapacity {
		c := uint64(capacity.CPU.Used)
		cores += c
	}
	return cores / speed
}

// Memory retrieves the amount of used memory to an oVDC, represented as a DataStorage type.
func (vdc *VDC) Memory() DataStorage {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	bm := float64(0)
	for _, capacity := range avdc.AdminVdc.ComputeCapacity {
		switch capacity.Memory.Units {
		case "KB":
			bm += float64(capacity.Memory.Used * kb)
		case "MB":
			bm += float64(capacity.Memory.Used * mb)
		case "GB":
			bm += float64(capacity.Memory.Used * gb)
		case "TB":
			bm += float64(capacity.Memory.Used * tb)
		default:
			bm += float64(capacity.Memory.Used)
		}
	}
	return DataStorage(bm)
}

// allStorageProfiles retrieves all storage profiles for a VDC and filters out partial duplicates, for
// example, when Veeam CDP creates a storage profile for the datastore.
func (vdc *VDC) allStorageProfiles() ([]*types.VdcStorageProfile, error) {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return nil, err
	}
	profiles := make([]*types.VdcStorageProfile, 0)
	for _, stor := range avdc.AdminVdc.VdcStorageProfiles.VdcStorageProfile {
		sp, err := vdc.Client.VCD.GetStorageProfileById(stor.ID)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, sp)
	}
	return profiles, nil
}

// storageProfile retrieves the default storage profile for the VDC.
func (vdc *VDC) defaultStorageProfile() (*types.VdcStorageProfile, error) {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return nil, err
	}
	ref, err := avdc.GetDefaultStorageProfileReference()
	if err != nil {
		return nil, err
	}
	sp, err := vdc.Client.VCD.GetStorageProfileById(ref.ID)
	if err != nil {
		return nil, err
	}
	return sp, nil
}

// Storage retrieves the total amount of 'requested' storage for an oVDC using the oVDC default
// storage policy.
func (vdc *VDC) Storage() DataStorage {
	profile, err := vdc.defaultStorageProfile()
	if err != nil {
		return 0
	}
	sb := profile.StorageUsedMB * mb
	return DataStorage(sb)
}

// StorageAll retrieves the total amount of used storage for an oVDC, totaling the 'requested'
// storage for all storage policies.
func (vdc *VDC) StorageAll() DataStorage {
	profiles, err := vdc.allStorageProfiles()
	if err != nil {
		return 0
	}
	bs := float64(0)
	for _, prof := range profiles {
		sb := prof.StorageUsedMB * mb
		bs += float64(sb)
	}
	return DataStorage(bs)
}

// VMCount retrieves the number of VMs deployed in the VDC.
func (vdc *VDC) VMCount() uint64 {
	ovdc, err := vdc.AdminOrg.GetVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	vms, err := ovdc.QueryVmList(types.VmQueryFilterAll)
	if err != nil {
		return 0
	}

	count := uint64(0)
	for _, vm := range vms {
		vm := vm
		if isValidStatus(vm.Status) && !vm.VAppTemplate && !vm.Deleted {
			count++
		}
	}
	return count
}

// PoweredOnVMCount retrieves the number of powered-on VMs deployed in the VDC.
func (vdc *VDC) PoweredOnVMCount() uint64 {
	ovdc, err := vdc.AdminOrg.GetVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	vms, err := ovdc.QueryVmList(types.VmQueryFilterAll)
	if err != nil {
		return 0
	}
	count := uint64(0)
	for _, vm := range vms {
		if vm.Status == types.VAppStatuses[4] && !vm.VAppTemplate && !vm.Deleted {
			count++
		}
	}
	return count
}

func (vdc *VDC) queryVMs(queries ...VMQuerySetter) ([]*types.QueryResultVMRecordType, error) {
	ovdc, err := vdc.AdminOrg.GetVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return nil, err
	}
	query := &VMQuery{
		Name:      nil,
		GuestOS:   nil,
		PoweredOn: false,
	}
	for _, set := range queries {
		set(query)
	}

	vms, err := ovdc.QueryVmList(types.VmQueryFilterOnlyDeployed)
	if err != nil {
		return nil, err
	}
	if query.PoweredOn {
		_vms := make([]*types.QueryResultVMRecordType, 0, len(vms))
		for _, vm := range vms {
			if vm.Status == types.VAppStatuses[4] {
				_vms = append(_vms, vm)
			}
		}
		vms = _vms
	}
	if query.Name != nil {
		_vms := make([]*types.QueryResultVMRecordType, 0, len(vms))
		for _, vm := range vms {
			if query.Name.MatchString(vm.Name) {
				_vms = append(_vms, vm)
			}
		}
		vms = _vms
	}
	if query.GuestOS != nil {
		_vms := make([]*types.QueryResultVMRecordType, 0, len(vms))
		for _, vm := range vms {
			if query.GuestOS.MatchString(vm.GuestOS) {
				_vms = append(_vms, vm)
			}
		}
		vms = _vms
	}
	return vms, nil
}

// VMCountWithQuery retrieves the number of VMs matching all of the provided queries.
// If PoweredOn is false (default), VMs that are both powered on or off will be included.
func (vdc *VDC) VMCountWithQuery(queries ...VMQuerySetter) uint64 {
	vms, err := vdc.queryVMs(queries...)
	if err != nil {
		return 0
	}
	return uint64(len(vms))
}

// VMCoreCountWithQuery retrieves the number cores on VMs matching all of the provided queries.
// If PoweredOn is false (default), VMs that are both powered on or off will be included.
func (vdc *VDC) VMCoreCountWithQuery(queries ...VMQuerySetter) uint64 {
	vms, err := vdc.queryVMs(queries...)
	if err != nil {
		return 0
	}
	count := uint64(0)
	for _, vm := range vms {
		count += uint64(vm.Cpus)
	}
	return count
}

// Org retrieves a vCloud Organization object.
func (client *Client) Org(orgID string) (*govcd.AdminOrg, error) {
	if !strings.HasPrefix(orgID, "urn:vcloud:org:") {
		orgID = fmt.Sprintf("urn:vcloud:org:%s", orgID)
	}
	org, err := client.VCD.GetAdminOrgById(orgID)
	if err != nil {
		err = errorx.Decorate(err, "failed to retrieve org '%s'", orgID)
		return nil, err
	}
	return org, nil
}

// VDC retrieves a single VDC associated with an organization by its ID and provides a wrapper for
// utilization functions for each VDC.
func (client *Client) VDC(orgID string, id string) (*VDC, error) {
	org, err := client.Org(orgID)
	if err != nil {
		return nil, err
	}
	obj, err := org.GetVDCById(id, false)
	if err != nil {
		err = errorx.Decorate(err, "failed to retrieve VDC '%s' for org '%s'", id, orgID)
		return nil, err
	}
	vdc := &VDC{
		Obj:      obj,
		AdminOrg: org,
		Client:   client,
	}
	return vdc, nil
}

// VDCs retrieves all VDCs associated with an organization and provides a wrapper for utilization
// functions for each VDC.
func (client *Client) VDCs(orgID string) (VDCs, error) {
	org, err := client.Org(orgID)
	if err != nil {
		return nil, err
	}
	vdcs, err := org.GetAllVDCs(false)
	if err != nil {
		err = errorx.Decorate(err, "failed to retrieve VDCs for org '%s'", orgID)
		return nil, err
	}
	vdcObjs := make([]VDC, 0, len(vdcs))
	for _, vdc := range vdcs {
		vdcObjs = append(vdcObjs, VDC{Obj: vdc, AdminOrg: org, Client: client})
	}
	return vdcObjs, nil
}
