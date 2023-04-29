package provider

import (
	"log"
	"path/filepath"
	"github.com/arnouthoebreckx/terraform-provider-synology/client"
)

type FileItemService struct {
	synologyClient client.SynologyClient
}

func (service FileItemService) Create(filename string, contents []byte) error {
	log.Println("Create " + string(contents))
	path, filename := getPathAndFilenameFromFullPath(filename)
	return service.synologyClient.Upload(path, true, true, filename, contents)
}

func (service FileItemService) Read(filename string) ([]byte, error) {
	return service.synologyClient.Download(filename)
}

func (service FileItemService) Update(filename string, contents []byte) ([]byte, error) {
	log.Println("Update " + string(contents))

	path, filename := getPathAndFilenameFromFullPath(filename)
	err := service.synologyClient.Upload(path, true, true, filename, contents)
	return contents, err
}

func (service FileItemService) Delete(filename string) error {
	return service.synologyClient.Delete(filename, false)

}

func getPathAndFilenameFromFullPath(fullPath string) (string, string) {
	return filepath.Dir(fullPath), filepath.Base(fullPath)
}

type FolderItemService struct {
	synologyClient client.SynologyClient
}

func (service FolderItemService) Create(path string) error {
	log.Println("Create Folder" + string(path))
	basePath, name := getPathAndFilenameFromFullPath(path)
	_, error := service.synologyClient.CreateFolder(basePath, name, true, "")
	return error
}

func (service FolderItemService) Delete(path string) error {
	log.Println("Delete Folder" + string(path))
	return service.synologyClient.Delete(path, true)

}

type GuestService struct {
	synologyClient client.SynologyClient
}

func (service GuestService) Create(name string, storage_id string, storage_name string, vnics []interface{}, vdisks []interface{}) (error) {
	log.Println("Create VMM Guest " + string(name))
	_, err := service.synologyClient.CreateGuest(name, storage_id, storage_name, vnics, vdisks)
	return err
}

func (service GuestService) Set(name string, autorun int, description string, vcpu_num int, vram_size int) (error) {
	log.Println("Setting values on VMM Guest " + string(name))
	log.Println("Values: " + string(autorun) + " " + description + " " + string(vcpu_num) + " " + string(vram_size))
	err := service.synologyClient.SetGuest(name, autorun, description, vcpu_num, vram_size)
	return err
}

func (service GuestService) Read(name string) ([]byte, error) {
	log.Println("Read VMM Guest " + string(name))
	content, err := service.synologyClient.ReadGuest(name)
	return content, err
}

func (service GuestService) Update(name string, new_name string) (error) {
	log.Println("Update VMM Guest from " + string(name) + " to " + string(new_name))
	err := service.synologyClient.UpdateGuest(name, new_name)
	return err
}

func (service GuestService) Delete(name string) (error) {
	log.Println("Delete VMM Guest" + string(name))
	err := service.synologyClient.DeleteGuest(name)
	return err
}

// func (service GuestService) Read(name string) error {
// 	log.Println("Create VMM Guest" + string(name))
// 	return service.synologyClient.CreateGuest(name, storage_id, storage_name, vnics, vdisks)
// }

// func (service GuestService) Delete() error {
// 	log.Println("Delete VMM Guest")
// 	return service.synologyClient.
// }
