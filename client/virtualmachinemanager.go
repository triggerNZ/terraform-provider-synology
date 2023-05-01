package client

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type readReponse struct {
	Data    Guest `json:"data"`
	Success bool  `json:"success"`
}

type Guest struct {
	Autorun     int     `json:"autorun"`
	Description string  `json:"description"`
	GuestId     string  `json:"guest_id"`
	GuestName   string  `json:"guest_name"`
	Status      string  `json:"status"`
	StorageId   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
	VcpuNum     int     `json:"vcpu_num"`
	Vdisks      []VDisk `json:"vdisks"`
	Vnics       []VNic  `json:"vnics"`
	VramSize    int     `json:"vram_size"`
}

type VNic struct {
	Mac         string `json:"mac"`
	Model       int    `json:"model"`
	NetworkID   string `json:"network_id"`
	NetworkName string `json:"network_name"`
	VnicID      string `json:"vnic_id"`
}

type VDisk struct {
	Controller int    `json:"controller"`
	Unmap      bool   `json:"unmap"`
	VdiskId    string `json:"vdisk_id"`
	VdiskSize  int    `json:"vdisk_size"`
}

type GuestInfo struct {
	auto_clean_task bool
	guest_id        string
	progress        int
	status          string
}

type CreateGuestResponse struct {
	task_id string
}

func CreateGuest(apiInfo map[string]InfoData, host string, sid string, name string, storage_id string, storage_name string, vnics []interface{}, vdisks []interface{}) (CreateGuestResponse, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	vnicsString, _ := json.Marshal(vnics)
	vdisksString, _ := json.Marshal(vdisks)

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "create"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	if storage_id != "" {
		queryString["storage_id"] = storage_id
	}
	if storage_name != "" {
		queryString["storage_name"] = storage_name
	}
	queryString["vnics"] = string(vnicsString)
	queryString["vdisks"] = string(vdisksString)

	wsUrl := host + "/webapi/entry.cgi"
	_, body, err := HttpCall(wsUrl, queryString)

	if err != nil {
		return CreateGuestResponse{}, err
	}

	log.Println("Create VMM Guest body" + string(body))

	var CreateGuestResponse CreateGuestResponse
	json.Unmarshal(body, &CreateGuestResponse)

	return CreateGuestResponse, nil
}

func SetGuest(apiInfo map[string]InfoData, host string, sid string, name string, autorun int, description string, vcpu_num int, vram_size int) ([]byte, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "set"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	queryString["autorun"] = strconv.Itoa(autorun)
	queryString["description"] = description
	if vcpu_num != 0 {
		queryString["vcpu_num"] = strconv.Itoa(vcpu_num)
	}
	if vram_size != 0 {
		queryString["vram_size"] = strconv.Itoa(vram_size)
	}

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func ReadGuest(apiInfo map[string]InfoData, host string, sid string, name string) (Guest, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "get"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	queryString["additional"] = "true"

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return Guest{}, err
	}

	response := readReponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err.Error())
		return Guest{}, err
	}

	return response.Data, nil
}

func UpdateGuest(apiInfo map[string]InfoData, host string, sid string, name string, new_name string) (int, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "set"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	queryString["new_guest_name"] = new_name

	wsUrl := host + "/webapi/entry.cgi"

	statusCode, _, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func DeleteGuest(apiInfo map[string]InfoData, host string, sid string, name string) (int, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "delete"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name

	wsUrl := host + "/webapi/entry.cgi"

	statusCode, _, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func (g Guest) String() string {
	str := fmt.Sprintf("Guest:\n\tGuestName: %s\n\tGuestId: %s\n\tAutorun: %d\n\tDescription: %s\n\tStatus: %s\n\tStorageName: %s\n\tStorageId: %s\n\tVcpuNum: %d\n\tVramSize: %d\n\tVdisks: [\n", g.GuestName, g.GuestId, g.Autorun, g.Description, g.Status, g.StorageName, g.StorageId, g.VcpuNum, g.VramSize)
	for _, vdisk := range g.Vdisks {
		str += fmt.Sprintf("\t\t%s\n", vdisk.String())
	}
	str += "\t]\n\tVnics: [\n"
	for _, vnic := range g.Vnics {
		str += fmt.Sprintf("\t\t%s\n", vnic.String())
	}
	str += "\t]\n"
	return str
}

func (vnic VNic) String() string {
	return fmt.Sprintf("VNic:\n\tMac: %s\n\tModel: %s\n\tNetworkID: %s\n\tNetworkName: %s\n\tVnicID: %s", vnic.Mac, vnic.Model, vnic.NetworkID, vnic.NetworkName, vnic.VnicID)
}

func (vdisk VDisk) String() string {
	return fmt.Sprintf("VDisk:\n\tController: %d\n\tUnmap: %t\n\tVdiskId: %s\n\tVdiskSize: %d", vdisk.Controller, vdisk.Unmap, vdisk.VdiskId, vdisk.VdiskSize)
}
