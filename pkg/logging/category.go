package logging

type Categoty string
type SubCategoty string
type ExtraKey string

const (
	General         Categoty = "General"
	Internal        Categoty = "Internal"
	Postgres        Categoty = "Postgres"
	Redis           Categoty = "Redis"
	RequestResponse Categoty = "RequestResponse"
)

const (
	//General
	Startup         SubCategoty = "Startub"
	ExternalService SubCategoty = "ExternalService"

	//Postgres
	Migration SubCategoty = "Migration"
	Select    SubCategoty = "Select"
	Rollback  SubCategoty = "Rollback"
	Update    SubCategoty = "Update"
	Delete    SubCategoty = "Delete"
	Insert    SubCategoty = "Insert"

	//Internal
	Api                 SubCategoty = "Api"
	HashPassword        SubCategoty = "HashPassword"
	DefaultRoleNotFound SubCategoty = "DefaultRoleNotFound"

	//Validation
	MobileValidation   SubCategoty = "MobileValidation"
	PasswordValidation SubCategoty = "PasswordValidation"
)

const (
	AppName      ExtraKey = "AppName"
	LoggerName   ExtraKey = "LoggerName"
	ClientIp     ExtraKey = "ClientIp"
	HostIp       ExtraKey = "HostIp"
	Method       ExtraKey = "Method"
	StatusCode   ExtraKey = "StatusCode"
	BodySize     ExtraKey = "BodySize"
	Path         ExtraKey = "Path"
	Latency      ExtraKey = "Latency"
	RequestBody  ExtraKey = "RequestBody"
	ResponseBody ExtraKey = "ResponseBody"
	ErrorMessage ExtraKey = "ErrorMessage"
)
