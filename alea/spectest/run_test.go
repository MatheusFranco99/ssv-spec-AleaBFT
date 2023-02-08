package spectest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	tests2 "github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	for _, test := range AllTests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func TestJson(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate")
	fileName := "tests.json"
	untypedTests := map[string]interface{}{}
	byteValue, err := os.ReadFile(path + "/" + fileName)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(byteValue, &untypedTests); err != nil {
		panic(err.Error())
	}

	tests := make(map[string]SpecTest)
	for name, test := range untypedTests {
		testName := test.(map[string]interface{})["Name"].(string)
		t.Run(testName, func(t *testing.T) {
			testType := strings.Split(name, "_")[0]
			switch testType {
			case reflect.TypeOf(&tests2.MsgProcessingSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.MsgProcessingSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				// a little trick we do to instantiate all the internal instance params
				preByts, _ := typedTest.Pre.Encode()
				pre := alea.NewInstance(
					testingutils.TestingConfigAlea(testingutils.KeySetForShare(typedTest.Pre.State.Share)),
					typedTest.Pre.State.Share,
					typedTest.Pre.State.ID,
					typedTest.Pre.State.Height,
				)
				err = pre.Decode(preByts)
				require.NoError(t, err)
				typedTest.Pre = pre

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.MsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.MsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.CreateMsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.CreateMsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			default:
				panic("unsupported test type " + testType)
			}
		})
	}
}
