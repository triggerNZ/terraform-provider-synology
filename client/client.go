package client

import (
	"log"
)

type SynologyClient interface {
	Connect(host string, username string, password string) error
	Disconnect() error
	CreateFolder(folderPath string, name string, forceParent bool, additional string) (CreateFolderResponse, error)
	Download(path string) ([]byte, error)
	Delete(path string, recursive bool) error
	Upload(path string, createParents bool, overwrite bool, fileName string, fileContents []byte) error
	CreateGuest(name string, storage_id string, storage_name string, vnics []interface{}, vdisks []interface{}) (CreateGuestResponse, error)
	SetGuest(name string, autorun int, description string, vcpu_num int, vram_size int) (error)
	ReadGuest(name string) ([]byte, error)
	UpdateGuest(name string, new_name string) (error)
	DeleteGuest(name string) (error)
}

type synologyClient struct {
	host    string
	apiInfo map[string]InfoData
	sid     string
}

func (client *synologyClient) Connect(host string, username string, password string) error {
	log.Println("Connect")
	var err error
	client.host = host
	client.apiInfo, err = Info(client.host, "all")
	if err != nil {
		return err
	}
	loginResponse, err := Login(client.apiInfo, client.host, username, password, "sergief-synology-client", "cookie")
	if err != nil {
		return err
	}
	client.sid = loginResponse.Sid

	return nil
}

func (client synologyClient) Disconnect() error {
	log.Println("Disconnect")

	_, err := Logout(client.apiInfo, client.host, client.sid)
	return err
}

func (client synologyClient) CreateFolder(folderPath string, name string, forceParent bool, additional string) (CreateFolderResponse, error) {
	log.Println("CreateFolder")

	return CreateFolder(client.apiInfo, client.host, client.sid, folderPath, name, forceParent, additional)
}

func (client synologyClient) Download(path string) ([]byte, error) {
	log.Println("Download")

	statusCode, body, err := Download(client.apiInfo, client.host, client.sid, path)
	if statusCode != 200 {
		return []byte(""), err
	}

	return body, err
}

func (client synologyClient) Delete(path string, recursive bool) error {
	log.Println("Delete")

	_, err := Delete(client.apiInfo, client.host, client.sid, path, recursive)
	return err
}

func (client synologyClient) Upload(path string, createParents bool, overwrite bool, fileName string, fileContents []byte) error {
	log.Println("Upload path=" + path + " filename=" + fileName)

	statusCode, err := Upload(client.apiInfo, client.host, client.sid, path, createParents, overwrite, fileName, fileContents)
	log.Println(statusCode)

	return err
}

func (client synologyClient) CreateGuest(name string, storage_id string, storage_name string, vnics []interface{}, vdisks []interface{}) (CreateGuestResponse, error) {
	return CreateGuest(client.apiInfo, client.host, client.sid, name, storage_id, storage_name, vnics, vdisks)
}

func (client synologyClient) SetGuest(name string, autorun int, description string, vcpu_num int, vram_size int) (error) {
	body, err := SetGuest(client.apiInfo, client.host, client.sid, name, autorun, description, vcpu_num, vram_size)
	log.Println(body)
	return err
}

func (client synologyClient) ReadGuest(name string) ([]byte, error){
	return ReadGuest(client.apiInfo, client.host, client.sid, name)
}

func (client synologyClient) UpdateGuest(name string, new_name string) (error) {
	statusCode, err := UpdateGuest(client.apiInfo, client.host, client.sid, name, new_name)
	log.Println(statusCode)
	return err
}

func (client synologyClient) DeleteGuest(name string) (error) {
	statusCode, err := DeleteGuest(client.apiInfo, client.host, client.sid, name)
	log.Println(statusCode)
	return err
}

func NewClient() SynologyClient {
	return &synologyClient{}
}
