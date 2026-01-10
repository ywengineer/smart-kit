package types

type Env string

func (e Env) String() string {
	return string(e)
}

const Development Env = "dev"
const UAT Env = "uat"
const Production Env = "prod"
