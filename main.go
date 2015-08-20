package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/groob/go-pkgbuild/pkginfo"
)

var (
	fRoot       = flag.String("root", "", "<root-path>")
	fIdentifier = flag.String("identifier", "", "identifier")
	fVersion    = flag.String("version", "", "version")
	pkgName     string
)

func init() {
	flag.Parse()
	pkgName = flag.Args()[0]
	err := os.MkdirAll("/tmp/"+pkgName, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := payload()
	if err != nil {
		log.Fatal(err)
	}
	// xmlFile, err := os.Open("PackageInfo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer xmlFile.Close()
	// XMLdata, _ := ioutil.ReadAll(xmlFile)
	// var info pkginfo.PackageInfo
	// info.Read(XMLdata)
	// info.Identifier = *fIdentifier
	// info.Version = *fVersion
	// output, err := info.Write()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// os.Stdout.Write(output)
	//
	bom()
	pkgInfo()
	pkg()
}

func bom() error {
	bomCmd := exec.Command("/usr/bin/mkbom", *fRoot, "/tmp/"+pkgName+"/Bom")
	return bomCmd.Run()
}

func payload() error {
	cmd := fmt.Sprintf("( cd %v && find . | cpio -o --format odc --owner 0:80 | gzip -c ) > %v/Payload", *fRoot, "/tmp/"+pkgName)
	payloadCmd := exec.Command("bash", "-c", cmd)
	return payloadCmd.Run()
}

func pkgInfo() error {
	var info pkginfo.PackageInfo
	info.Identifier = *fIdentifier
	info.Version = *fVersion
	info.OverwritePermissions = true
	info.FormatVersion = 2
	info.GeneratorVersion = "InstallCmds-502 (14F27)"
	info.Auth = "root"
	info.PostInstallAction = "none"
	info.Payload.NumberOfFiles = 5
	info.Payload.InstallKBytes = 2
	output, err := info.Write()
	if err != nil {
		return err
	}
	os.Stdout.Write(output)
	return nil
}

func pkg() error {
	cmd := fmt.Sprintf("( cd %v && /usr/bin/xar --compression none -cf %v *)", "/tmp/"+pkgName, "/tmp/new-foo.pkg")
	xarCmd := exec.Command("bash", "-c", cmd)
	// xarCmd := exec.Command("/usr/bin/xar", "--compression", "none", "-cf", pkgName, *fRoot)
	return xarCmd.Run()
}
