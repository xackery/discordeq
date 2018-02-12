# DiscordEQ
This plugin allows Everquest to communicate with Discord in a bidirectional manner.

## How to Configure

### Edit eqemu_config
* Add to eqemu_config.json these lines under the server {} section:
```json
     "discord" : {
         "channelid" : "{CHANNELID}",
         "itemurl" : "{ITEMURL}",
         "refreshrate" : "15",
         "clientid": "{CLIENTID}",
         "serverid" : "{SERVERID}",
         "username" : "{YOURTOKEN}",
         "commandchannelid": "{YOURCOMMANDCHANNELID}"
      },
```
* each section above with {WORD} sections will be populated with data as you follow below steps.

### Prepare App/Bot in Discord
* Go to https://discordapp.com/developers/
* Click My Apps on the top left area.
* Click + New App
* Write anything you wish for the app name, click Create App
* Copy the clientid into your discord into the {CHANNELID} field in eqemu_config
* Scroll down to the bot section, and click Create Bot User
* Confirm with Yes, Do it!
* Make sure Public Bot is unchecked, as well as oauth2 grant.
* Save, click to reveal the token. Copy the bot token into your {YOURTOKEN} field in eqemu_config, under username

### Prepare server, channel IDs
* Inside Discord, right click your server's circular icon and on the bottom choose Copy ID.
* Paste the serverID into your eqemu_config's {SERVERID} field
* Create a new channel for OOC chat. Right click the channels' name and copy ID.
* Paste the channelID into your eqemu_config's {CHANNELID} field

### Prepare itemURL (optional)
* If you have a website that has item links, replace the {ITEMURL} field with the website, e.g. "http://yoursite.com/item.php?id=". It is assumed your item id's are appended to the end of the url link.
* If you don't have a website with itemlinks, remove the {ITEMURL} entry all together and keep it empty. Discord will italics itemlinks in game when displayed on Discord then.



### Run DiscordEQ
* Your first run should fail with an unauthorized notification, since you have not given your bot permissions to your server yet.  You will see a link on the bottom you can copy paste into a browser, and give it permission to access your server.


### Enabling Players to talk from Discord to EQ
* Admin-level accounts can only do the following steps.
* To allow this, inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. The `IGN:` prefix is required for DiscordEQ to detect a player and is used to identify the player in game, For example, to identify the discord user `Xackery` as `Shin`, I would create a role named `IGN: Shin`, right click the user Xackery, and assign the role to them.
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`
