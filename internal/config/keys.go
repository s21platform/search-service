package config

type key string

const KeyUUID key = key("uuid")
const KeyMetrics = key("metrics")
const KeyLogger = key("logger")
const TotalReadUser = int64(5)
