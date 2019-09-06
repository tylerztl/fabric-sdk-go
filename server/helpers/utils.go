package helpers

import (
	"go/build"
	"path"
	"path/filepath"
)

// goPath returns the current GOPATH. If the system
// has multiple GOPATHs then the first is used.
func goPath() string {
	gpDefault := build.Default.GOPATH
	gps := filepath.SplitList(gpDefault)

	return gps[0]
}

func GetConfigPath(filename string) string {
	const configPath = "conf"
	return path.Join(goPath(), "src", Project, configPath, filename)
}

func GetChannelConfigPath(filename string) string {
	return path.Join(goPath(), "src", Project, ChannelConfigPath, filename)
}

func GetDeployPath() string {
	const ccPath = "artifacts/chaincode"
	return path.Join(goPath(), "src", Project, ccPath)
}
