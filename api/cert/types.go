package cert

import "time"

type SignRequest struct {
	DomainList  []string `json:"domain_list"`
	IpList      []string `json:"ip_list"`
	ValidDays   string   `json:"valid_days"`
	CsrText     string   `json:"csr_text"`
	RequestedBy string   `json:"requested_by"`
}

type SignRequestResponse struct {
	Id           string   `json:"id"`
	Organization string   `json:"organization"`
	CommonName   string   `json:"common_name"`
	DomainList   []string `json:"domain_list"`
	IpList       []string `json:"ip_list"`
	ValidDays    int      `json:"valid_days"`
	State        string   `json:"state"`
	RequestedIp  string   `json:"requested_ip"`
	RequestedBy  string   `json:"requested_by"`
	SubmitURL    string   `json:"submit_url"`
}

type AuthorityRequest struct {
	Name             string `json:"name"`
	Organization     string `json:"organization"`
	Domain           string `json:"domain"`
	RootValidDays    int    `json:"root_valid_days"`
	DefaultValidDays int    `json:"default_valid_days"`
	MaxValidDays     int    `json:"max_valid_days"`
	Agent            string `json:"agent"`
	Owner            string `json:"owner"`
	Install          bool   `json:"install"`
}

type AuthorityCreateResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Organization     string    `json:"organization"`
	Domain           string    `json:"domain"`
	RootValidDays    int       `json:"root_valid_days"`
	DefaultValidDays int       `json:"default_valid_days"`
	MaxValidDays     int       `json:"max_valid_days"`
	Agent            string    `json:"agent"`
	Owner            string    `json:"owner"`
	Instruction      string    `json:"instruction"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AuthorityListResponse struct {
	Count    int                 `json:"count"`
	Current  int                 `json:"current"`
	Next     int                 `json:"next"`
	Previous string              `json:"previous"`
	Last     int                 `json:"last"`
	Results  []AuthorityResponse `json:"results"`
}

type AuthorityResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Organization     string    `json:"organization"`
	Domain           string    `json:"domain"`
	RootValidDays    int       `json:"root_valid_days"`
	DefaultValidDays int       `json:"default_valid_days"`
	MaxValidDays     int       `json:"max_valid_days"`
	Agent            string    `json:"agent"`
	AgentName        string    `json:"agent_name"`
	Owner            string    `json:"owner"`
	OwnerName        string    `json:"owner_name"`
	UpdatedAt        time.Time `json:"updated_at"`
	SignedAt         time.Time `json:"signed_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type AuthorityAttributes struct {
	Name             string `json:"name"`
	Organization     string `json:"organization"`
	Domain           string `json:"domain"`
	RootValidDays    int    `json:"root_valid_days"`
	DefaultValidDays int    `json:"default_valid_days"`
	MaxValidDays     int    `json:"max_valid_days"`
	Server           string `json:"server"` // agentName -> server
	Owner            string `json:"owner"`  // Owner name -> owner
	SignedAt         string `json:"signed_at"`
}

type CSRSubmit struct {
	CsrText string `json:"csr_text"`
}

type CSRListResponse struct {
	Count    int           `json:"count"`
	Current  int           `json:"current"`
	Next     int           `json:"next"`
	Previous string        `json:"previous"`
	Last     int           `json:"last"`
	Results  []CSRResponse `json:"results"`
}

type CSRResponse struct {
	Id              string    `json:"id"`
	Authority       string    `json:"authority"`
	AuthorityName   string    `json:"authority_name"`
	CommonName      string    `json:"common_name"`
	DomainList      []string  `json:"domain_list"`
	IpList          []string  `json:"ip_list"`
	ValidDays       int       `json:"valid_days"`
	State           string    `json:"state"`
	RequestedIp     string    `json:"requested_ip"`
	RequestedBy     string    `json:"requested_by"`
	RequestedByName string    `json:"requested_by_name"`
	AddedAt         time.Time `json:"added_at"`
}

type CSRAttributes struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"` // Derived from the first domain in the CSR domain list
	Authority     string   `json:"authority"`
	DomainList    []string `json:"domain_list"`
	IpList        []string `json:"ip_list"`
	State         string   `json:"state"`
	RequestedIp   string   `json:"requested_ip"`
	RequestedBy   string   `json:"requested_by"`
	RequestedDate string   `json:"requested_date"`
}

type CSRDetails struct {
	Id              string    `json:"id"`
	Authority       string    `json:"authority"`
	AuthorityName   string    `json:"authority_name"`
	CommonName      string    `json:"common_name"`
	DomainList      []string  `json:"domain_list"`
	IpList          []string  `json:"ip_list"`
	ValidDays       int       `json:"valid_days"`
	State           string    `json:"state"`
	CsrText         string    `json:"csr_text"`
	RequestedIp     string    `json:"requested_ip"`
	RequestedBy     string    `json:"requested_by"`
	RequestedByName string    `json:"requested_by_name"`
	AddedAt         time.Time `json:"added_at"`
}

type CertificateListResponse struct {
	Count    int           `json:"count"`
	Current  int           `json:"current"`
	Next     int           `json:"next"`
	Previous string        `json:"previous"`
	Last     int           `json:"last"`
	Results  []Certificate `json:"results"`
}

type Certificate struct {
	Id           string    `json:"id"`
	Authority    string    `json:"authority"`
	Csr          string    `json:"csr"`
	CrtText      string    `json:"crt_text"`
	ValidDays    int       `json:"valid_days"`
	SignedAt     time.Time `json:"signed_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	SignedBy     string    `json:"signed_by"`
	SignedByName string    `json:"signed_by_name"`
}

type CertificateAttributes struct {
	Id        string `json:"id"`
	Authority string `json:"authority"`
	Csr       string `json:"csr"`
	ValidDays int    `json:"valid_days"`
	SignedAt  string `json:"signed_at"`
	ExpiresAt string `json:"expires_at"`
	SignedBy  string `json:"signed_by"`
}
