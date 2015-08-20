package pkginfo

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

type PackageInfo struct {
	PkgInfo `xml:"pkg-info"`
}

type PkgInfo struct {
	XMLName              xml.Name           `xml:"pkg-info"`
	OverwritePermissions bool               `xml:"overwrite-permissions,attr"`
	Relocatable          bool               `xml:"relocatable,attr"`
	Identifier           string             `xml:"identifier,attr"`
	PostInstallAction    string             `xml:"postinstall-action,attr"`
	Version              string             `xml:"version,attr"`
	FormatVersion        int                `xml:"format-version,attr"`
	GeneratorVersion     string             `xml:"generator-version,attr"`
	InstallLocation      string             `xml:"install-location,attr"`
	Auth                 string             `xml:"auth,attr"`
	Payload              payload            `xml:"payload"`
	BundleVersion        bundleVersion      `xml:"bundle-version"`
	UpgradeBundle        upgradeBundle      `xml:"upgrade-bundle"`
	UpdateBundle         updateBundle       `xml:"update-bundle"`
	AtomicUpdateBundle   atomicUpdateBundle `xml:"atomic-update-bundle"`
	StrictIdentifier     strictIdentifier   `xml:"strict-identifier"`
	Relocate             relocate           `xml:"relocate"`
	Scripts              scripts            `xml:"scripts,omitempty"`
}

type payload struct {
	NumberOfFiles int `xml:"numberOfFiles,attr"`
	InstallKBytes int `xml:"installKBytes,attr"`
}

type bundleVersion struct {
}
type upgradeBundle struct {
}
type updateBundle struct {
}
type atomicUpdateBundle struct {
}
type strictIdentifier struct {
}
type relocate struct {
}

type scripts struct {
	Postinstall struct {
		File string `xml:"file,attr"`
	} `xml:"postinstall"`
}

func (p *PkgInfo) Read(XMLdata []byte) error {
	return xml.Unmarshal(XMLdata, &p)
}
func (p *PkgInfo) Write() ([]byte, error) {
	var output []byte
	output, err := xml.MarshalIndent(p, "", "    ")
	if err != nil {
		return output, err
	}
	output = []byte(xml.Header + string(output))
	return output, err
}

func main() {
	xmlFile, err := os.Open("PackageInfo")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	XMLdata, _ := ioutil.ReadAll(xmlFile)
	var info PackageInfo
	info.Read(XMLdata)
	output, err := info.Write()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(output)
}
