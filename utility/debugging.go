package utility

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

var preRequestMiddleware = func(c *resty.Client, req *resty.Request) error {
	logrus.WithFields(logrus.Fields{
		"url":        c.HostURL,
		"authscheme": c.AuthScheme,
		"requestUrl": req.URL,
	}).Traceln("made api request")
	return nil
}

var postResponseMiddleware = func(c *resty.Client, resp *resty.Response) error {
	logrus.WithFields(logrus.Fields{
		"url":            c.HostURL,
		"authscheme":     c.AuthScheme,
		"responseStatus": resp.Status(),
	}).Traceln("received api response")
	return nil
}

var requestErrorMiddleware = func(req *resty.Request, err error) {
	if v, ok := err.(*resty.ResponseError); ok {
		logrus.WithFields(logrus.Fields{
			"lastResponse":  v.Response,
			"originalError": v.Err,
		}).Errorln("api request error")
	}
	logrus.Errorln("inextendable api request error")
}

// AddWrappersToOverseerr ensures that all overseerr requests will be logged
func AddWrappersToOverseerr(o *goverseerr.Overseerr) {
	o.RegisterPostResponseMiddleware(postResponseMiddleware)
	o.RegisterPreRequestMiddleware(preRequestMiddleware)
	o.RegisterRequestErrorMiddleware(requestErrorMiddleware)
}
