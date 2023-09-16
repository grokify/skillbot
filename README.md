# Skillbot

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

 [build-status-svg]: https://github.com/grokify/skillbot/workflows/test/badge.svg
 [build-status-url]: https://github.com/grokify/skillbot/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/skillbot
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/skillbot
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-url]: https://godoc.org/github.com/grokify/skillbot
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/skillbot/blob/master/LICENSE

Example Skillbot for Glip.

This is a work in progress app and not ready for use.

## Configuration

### AWS

* Runtime: Go 1.x
* Handler: `main`

```
CHATBLOX_ENGINE=awslambda
CHATBLOX_REQUEST_FUZZY_AT_MENTION_MATCH=true
CHATBLOX_RESPONSE_AUTO_AT_MENTION=true
RINGCENTRAL_BOT_ID=12345678
RINGCENTRAL_BOT_NAME=Skill Bot
RINGCENTRAL_TOKEN=myToken
RINGCENTRAL_SERVER_URL=https://platform.ringcentral.com
ALGOLIA_APP_CREDENTIALS_JSON={"applicationId": "myApplicationId", "searchOnlyApiKey": "mySearchOnlyApiKey", "adminApiKey": "myAdminApiKey", "analyticsApiKey": "myAnalyticsApiKey", "monitoringApiKey": "myMonitoringApiKey"}
ALGOLIA_INDEX=skillbot
```

#### Test

Use AWS Gateway Proxy POST wiht the following data with the data in [`sample.message.json`](sample.message.json) with same `ownerId` and `groupId` values for `nethttp`.


#### Example

`Glip> Who can help me book a flight with $airline`