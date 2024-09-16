package gapi

type ApiHeader struct {
	UserAgent   string `header:"User-Agent" binding:"eq=api" label:"user-agent"`
	XApiVersion string `header:"x-api-version" binding:"eq=1.0" label:"x-api-version"`
	XApiTs      int64  `header:"x-api-ts" binding:"required" label:"x-api-ts"`
	XApiKey     string `header:"x-api-key" binding:"required" label:"x-api-key"`
	XApiSign    string `header:"x-api-sign" binding:"required" label:"x-api-sign"`
}
