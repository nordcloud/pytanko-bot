# ⁉️Pytanko Bot

Pytanko Bot is a bot that collects questions and posts them to a GraphQL API. It is best used with the Pytanko UI (https://github.com/nordcloud/pytanko-ui) which can then present the question in an online page.

# How to use it

Register a slash command in your Slack team, e.g. `/pytanko`. 

# Deployment

Deplyments needs the an api key and GraphQL API uri. We assumed it runs on Amazon AppSync.

```
export API_URL=<your gql url>
export API_KEY=<your api key>
sls deploy -s staging 
```

# Authors

Made in Nordcloud Poznań with ♥️

