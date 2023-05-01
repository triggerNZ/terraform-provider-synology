package client

import (
	"encoding/json"
	"log"
	"strconv"
)

type Guest struct {
	autorun      int
	description  string
	guest_id     string
	guest_name   string
	status       string
	storage_id   string
	storage_name string
	vcpu_num     int
	vdisks       []VDisk
	vnics        []VNic
	vram_size    int
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

type VNic struct {
	mac          string
	network_id   string
	network_name string
}

type VDisk struct {
	create_type int
	vdisk_size  int
	image_id    string
	image_name  string
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

	var guest Guest

	_, body, err := HttpCall(wsUrl, queryString)
	errJson := json.Unmarshal(body, &guest)
	if err != nil {
		return Guest{}, err
	}

	if errJson != nil {
		return Guest{}, err
	}

	return guest, nil
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
