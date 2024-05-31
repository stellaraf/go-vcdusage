package vcdusage

import (
	"fmt"
	"strings"

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

// CoreCount retrieves the allocated CPU MHz for a VDC and calculates the number of cores allocated
// to the VDC by dividing the total allocated CPU MHz by the CPU speed.
//
// For example, if the speed is 3.1 GHz and the allocation amount is 49.6, the core count is 16.
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
		c := uint64(capacity.CPU.Allocated)
		cores += c
	}
	return cores / speed
}

// Memory retrieves the amount of allocated memory to an oVDC, represented as a DataStorage type.
func (vdc *VDC) Memory() DataStorage {
	avdc, err := vdc.AdminOrg.GetAdminVDCById(vdc.Obj.Vdc.ID, false)
	if err != nil {
		return 0
	}
	bm := float64(0)
	for _, capacity := range avdc.AdminVdc.ComputeCapacity {
		switch capacity.Memory.Units {
		case "KB":
			bm += float64(capacity.Memory.Allocated * kb)
		case "MB":
			bm += float64(capacity.Memory.Allocated * mb)
		case "GB":
			bm += float64(capacity.Memory.Allocated * gb)
		case "TB":
			bm += float64(capacity.Memory.Allocated * tb)
		default:
			bm += float64(capacity.Memory.Allocated)
		}
	}
	return DataStorage(bm)
}

// Storage retrieves the total amount of allocated storage for an oVDC. If multiple storage
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
		sb := sp.StorageTotalMB * mb
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

// VDCs retrieves all VDCs associated with an organization and provides a wrapper for utilization
// functions for each VDC.
func (client *Client) VDCs(orgID string) ([]VDC, error) {
	if !strings.HasPrefix(orgID, "urn:vcloud:org:") {
		orgID = fmt.Sprintf("urn:vcloud:org:%s", orgID)
	}
	org, err := client.VCD.GetAdminOrgById(orgID)
	if err != nil {
		err = errorx.Decorate(err, "failed to retrieve org '%s'", orgID)
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
