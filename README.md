# okta-service-app-creator

Creates a [service-to-service](https://developer.okta.com/docs/guides/implement-oauth-for-okta/create-oauth-app/) Okta Oauth App.

The command generates a YAML file in the [required format](https://github.com/okta/okta-sdk-golang#yaml-configuration) for the Okta SDK for Go.

    $ okta-service-app-creator -name ExampleApp -org dev-123456
    Wrote config to file: dev-123446-ExampleApp.yaml
    Created Okta App with Client ID: 0oa52k2qtEzyJdfF54x6

## Requirements

You need an [Okta API token](https://help.okta.com/en/prod/Content/Topics/Security/API.htm) to run the command, as well as an Okta organization to add the application to.

## Implementation

`okta-sdk-golang` does not currently support the parameters needed to actually create a service to service app, so this tool creates the app using an HTTP POST.

The POST request uses [these API docs](https://developer.okta.com/docs/reference/api/apps/#add-oauth-2-0-client-application) for reference.
