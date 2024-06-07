package vcdusage

import (
	"fmt"
	"strings"

	"github.com/destel/rill"
	"github.com/joomcode/errorx"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

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

// Storage retrieves the total amount of used storage for an oVDC. If multiple storage
// polices are attached an oVDC, the amount will include the sum of all storage policies.
func (vdc *VDC) Storage() DataStorage {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	bs := float64(0)
	for _, stor := range avdc.AdminVdc.VdcStorageProfiles.VdcStorageProfile {
		sp, err := vdc.Client.VCD.QueryProviderVdcStorageProfileByName(stor.Name, avdc.AdminVdc.ProviderVdcReference.HREF)
		if err != nil {
			return 0
		}
		sb := sp.StorageUsedMB * mb
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
		if vm.Deployed {
			count++
		}
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
