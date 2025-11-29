// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package userconfigs

// 璁よ瘉鐘舵€?

type UserIdentityStatus = string

const (
	UserIdentityStatusNone      UserIdentityStatus = "none"
	UserIdentityStatusSubmitted UserIdentityStatus = "submitted"
	UserIdentityStatusRejected  UserIdentityStatus = "rejected"
	UserIdentityStatusVerified  UserIdentityStatus = "verified"
)

// 璁よ瘉绫诲瀷

type UserIdentityType = string

const (
	UserIdentityTypeIDCard            UserIdentityType = "idCard"
	UserIdentityTypeEnterpriseLicense UserIdentityType = "enterpriseLicense"
)

// 缁勭粐绫诲瀷

type UserIdentityOrgType = string

const (
	UserIdentityOrgTypeEnterprise UserIdentityOrgType = "enterprise"
	UserIdentityOrgTypeIndividual UserIdentityOrgType = "individual"
)
