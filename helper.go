package main

import (
	"net/mail"
)

func extractEmails(addr []*mail.Address, _ ...error) []string {
	ret := []string{}

	for _, e := range addr {
		if e.Name != "" {
			ret = append(ret, e.Name+" <"+e.Address+">")
		} else {
			ret = append(ret, e.Address)
		}
	}

	return ret
}

func transformStdAddressToEmailAddress(addr []*mail.Address) []*EmailAddress {
	ret := []*EmailAddress{}

	for _, e := range addr {
		ret = append(ret, &EmailAddress{
			Address: e.Address,
			Name:    e.Name,
		})
	}

	return ret
}
