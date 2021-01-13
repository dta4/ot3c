package otc

import (
	"errors"
	"strconv"

	"github.com/dta4/ot3c/config"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = logrus.WithFields(
	logrus.Fields{
		"package": "otc",
	},
)

//Login logs into OTC with provided config. Wrappes OTP request if enabled. Saves current Client
func Login() (*golangsdk.ProviderClient, error) {
	log := log.WithField("method", "login")
	if config.OTCUsingOTP {
		log.Infoln("The current config requires an One-Time-Password to be entered. Please provide the current OTP.")
		prompt := promptui.Prompt{
			Label:    "OTP",
			Validate: validateOTP,
		}
		_, err := prompt.Run()
		if err != nil {
			log.WithError(err).Error("Error encounterd on OTP enter")
			return nil, err
		}

	}

	log.Debugln("Logging in OTC ...")
	client, err := openstack.NewClient(config.OTCCreds.IdentityEndpoint)
	if err != nil {
		log.WithError(err).Error("Error on OTC client creation")
		return nil, err
	}
	err = openstack.Authenticate(client, config.OTCCreds)
	if err != nil {
		log.WithError(err).Error("Error on OTC login")
		return nil, err
	}
	log.Debug("OTC login successful")

	return client, err
}

func validateOTP(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return err
	}
	if len(input) != 6 {
		return errors.New("OTP not long enough (6)")
	}
	return nil
}
