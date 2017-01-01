package main

func scontains(key string, options ...string) bool {
	for _, item := range options {
		if item == key {
			return true
		}
	}
	return false
}