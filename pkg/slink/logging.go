package slink

import "charm.land/log/v2"

func SetLogLevel(level string) error {
	var err error
	ll, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(ll)
	return nil
}
