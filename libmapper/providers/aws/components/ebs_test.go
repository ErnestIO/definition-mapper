package components

// Basic imports
import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	graph "gopkg.in/r3labs/graph.v2"
)

// EBSTestSuite : Test suite for ebs component
type EBSTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
	List                          map[string]EBSVolume
}

// SetupTest : Setup test suite
func (suite *EBSTestSuite) SetupTest() {
	var iops int64
	key := "test"
	iops = 2
	suite.List = make(map[string]EBSVolume)
	suite.List[EBSErrNilName] = EBSVolume{Name: ""}
	suite.List[EBSErrAvailabilityNameNil] = EBSVolume{Name: "test"}
	suite.List[EBSErrNilType] = EBSVolume{Name: "test", AvailabilityZone: "test"}
	suite.List[EBSErrNilEncryption] = EBSVolume{Name: "test", AvailabilityZone: "test", VolumeType: "test", Encrypted: true}
	suite.List[EBSErrInvalidType] = EBSVolume{Name: "test", AvailabilityZone: "test", VolumeType: "test", Encrypted: true, EncryptionKeyID: &key, Iops: &iops}
	iops = 16385
	suite.List[EBSErrInvalidSize] = EBSVolume{Name: "test", AvailabilityZone: "test", VolumeType: "io1", Encrypted: true, EncryptionKeyID: &key, Size: &iops}
}

// TestValidate : Testing validate method
func (suite *EBSTestSuite) TestValidate() {
	for key, obj := range suite.List {
		err := obj.Validate()
		fmt.Println(" - Validating error message : '" + key + "'")
		suite.Equal(err.Error(), key)
	}
}

// TestUpdate : Testing update method
func (suite *EBSTestSuite) TestUpdate() {
	g := EBSVolume{VolumeAWSID: "lol"}
	e := EBSVolume{Name: "test"}
	e.Update(&g)
	suite.Equal(e.VolumeAWSID, "lol")
	suite.Equal(e.ComponentType, TYPEEBSVOLUME)
	suite.Equal(e.ComponentID, TYPEEBSVOLUME+TYPEDELIMITER+e.Name)
	suite.Equal(e.ProviderType, PROVIDERTYPE)
	suite.Equal(e.DatacenterName, DATACENTERNAME)
	suite.Equal(e.DatacenterType, DATACENTERTYPE)
	suite.Equal(e.DatacenterRegion, DATACENTERREGION)
	suite.Equal(e.AccessKeyID, ACCESSKEYID)
	suite.Equal(e.SecretAccessKey, SECRETACCESSKEY)
}

// TestRebuild : Testing rebuild method
func (suite *EBSTestSuite) TestRebuild() {
	g := graph.Graph{}
	e := EBSVolume{Name: "test"}
	e.Rebuild(&g)
	suite.Equal(e.ComponentType, TYPEEBSVOLUME)
	suite.Equal(e.ComponentID, TYPEEBSVOLUME+TYPEDELIMITER+e.Name)
	suite.Equal(e.ProviderType, PROVIDERTYPE)
	suite.Equal(e.DatacenterName, DATACENTERNAME)
	suite.Equal(e.DatacenterType, DATACENTERTYPE)
	suite.Equal(e.DatacenterRegion, DATACENTERREGION)
	suite.Equal(e.AccessKeyID, ACCESSKEYID)
	suite.Equal(e.SecretAccessKey, SECRETACCESSKEY)
}

// TestEbsTestSuite : tests for ebs component
func TestEbsTestSuite(t *testing.T) {
	suite.Run(t, new(EBSTestSuite))
}
