# Skillbot

Example Skillbot for Glip.

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