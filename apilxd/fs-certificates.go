package apilxd

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"os"
	"fmt"
)

type fsCertInfo struct {
	ID			int 	`json:"ID"`
	Fingerprint	string	`json:"Fingerprint"`
	Type 		int		`json:"Type"`
	Name 		string	`json:"Name"`
	Certificate	string	`json:"Certificate"`
}

func fsCertsGet() (certs []*fsCertInfo, err error ) {
	certs = []*fsCertInfo{}
	if _, err := os.Stat("./lxdcerts"); err == nil {
		files, errread := ioutil.ReadDir("./lxdcerts")
		if errread != nil {
			fmt.Println(errread)
		}
		for _, file := range files {
			aux,err2 := fsLoadCert(file.Name())
			if err2 != nil {
				return nil,err2
			}
			certs = append(certs,&aux)
		}
		return certs,nil
	}
	return nil,err
}

func fsCertsGetSingle() (certs []*fsCertInfo) {
	certs = []*fsCertInfo{}
	if _, err := os.Stat("./lxdcerts"); err == nil {
		files, errread := ioutil.ReadDir("./lxdcerts")
		if errread != nil {
			fmt.Println(errread)
		}
		for _, file := range files {
			aux,err2 := fsLoadCert(file.Name())
			if err2 != nil {
				return nil
			}
			certs = append(certs,&aux)
		}
		return certs
	}
	return nil
}
/*
func fsCertsGetByFingerprint(fingerprint string) (certs []*fsCertInfo, err error) {
	certs = []fsCertInfo
	if _, err := os.Stat("./lxdcerts"); err == nil {
		bytes, erread := ioutil.ReadFile("./lxdcerts"+fingerprint)
		if erread != nil {
			fmt.Println(erread)
		}
		cert := fsCertInfo{}
		json.Unmarshal(bytes,&cert)
		append(certs,cert)
		return certs
	}
	return nil
}
*/
func fsCertSave(cert *fsCertInfo) error{
	buf, err := json.Marshal(cert)

	if err != nil {
		return err
	}

	if _, err := os.Stat("./lxdcerts/"+cert.Fingerprint); err == nil {
		return errors.New("File already exists")
	}

	ioutil.WriteFile(".lxdcerts/"+cert.Fingerprint, buf, 0644)


	return nil
}

func fsCertDelete(path string) error {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	} else { return errors.New("File doesn't exists.")}
	
	return nil
}

func fsLoadCert(name string) (certinfo fsCertInfo,err error) {
	if _, err := os.Stat("./lxdcerts/"+name); err == nil {
		content,_ := ioutil.ReadFile("./lxdcerts/"+name)
		cert := fsCertInfo{}
		json.Unmarshal(content,&cert)
		return cert,nil
	}
	return certinfo,err
}

