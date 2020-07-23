1.1.0

> MAJOR PATCH  

## Major

- Changed command handler so that DM capable commands can be executed in DM to shinpuru. [#49]  
![](https://i.imgur.com/AvS2HrA.png)

- The web frontend now has a section for guild security features like guild backups or invite blocking. [#124, #132]
![](https://i.imgur.com/0wHeE2k.png)  
This change also adds the following REST API endpoints:
  - [`GET /api/guilds/:guildid/backups`](https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs#get-guild-backups)
  - [`GET /api/guilds/:guildid/backups/:backupid/download`](https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs#download-guild-backups)
  - [`POST /api/guilds/:guildid/backups/toggle`](https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs#toggle-guild-backup-enabled)
  - [`POST /api/guilds/:guildid/inviteblock`](https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs#toggle-guild-invite-block-enabled)

- Code execution listener was reworked so that edited messages are also recognized. Also, the implementaiton now used a single event listener for reactions instead of registering one for each execution message. [#53]

## Minor

- Login session keys now also use the JWT implementation. This makes sessions independend from the database, which is more secure when a database leak occurs, and more practical to store session metadata in the session key. The key used for sessions is randomly generated on each startup and periodically after a specified time has elaped. Also it is only held in RAM during runtime for security reasons. [#123]

- Twitch Notification Thumbnails should now be less "static" due to Discord's CDN caching. This bypass attempt was realized by adding an `?rid=` query parameter with a random integer as value generated for each embed. [#129]

- The [`vote close` command](https://github.com/zekroTJA/shinpuru/wiki/Commands#vote) now has another optional parameter `nc` or `nochart` to skip chart generation on close. [#128]

## Bug Fixes

- Fix inviteblock command so that it is not requesting a permission level on enabling [#131]

## Backstage

- Because Discord will shut down their `discordapp.com` domain to switch to the new `discord.com` domain, all endpoints and URLs were changed to `discord.com`. [#130] 

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.1.0
```