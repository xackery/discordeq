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
         "username" : "{YOURTOKEN}"
      },
```
* each section above with {WORD} sections will be populated with data as you follow below steps.

### Prepare App/Bot in Discord
* Go to https://discordapp.com/developers/
* Click My Apps on the top left area.
* Click + New App
* Write anything you wish for the app name, click Create App
* Copy the clientid into your discord into the {CLIENTID} field in eqemu_config
* Scroll down to the bot section, and click Create Bot User
* Confirm with Yes, Do it!
* Make sure Public Bot is unchecked, as well as oauth2 grant.
* Save, click to reveal the token. Copy the bot token into your {YOURTOKEN} field in eqemu_config, under username
* Visit the link: https://discordapp.com/oauth2/authorize?&client_id={CLIENTID}&scope=bot&permissions=2146958591 changing {CLIENTID} to your client ID.

### Prepare server, channel IDs
* Inside Discord, right click your server's circular icon and on the bottom choose Copy ID.
* Paste the serverID into your eqemu_config's {SERVERID} field
* Create a new channel for OOC chat. Right click the channels' name and copy ID.
* Paste the channelID into your eqemu_config's {CHANNELID} field

### Prepare itemURL (optional)
* If you have a website that has item links, replace the {ITEMURL} field with the website, e.g. "http://yoursite.com/item.php?id=". It is assumed your item id's are appended to the end of the url link.
* If you don't have a website with itemlinks, remove the {ITEMURL} entry all together and keep it empty. Discord will italics itemlinks in game when displayed on Discord then.

### Enable Telnet
* Inside your eqemu_config.json, look for a section that reads:
```json
"telnet": {
     "ip": "0.0.0.0",
     "port": "9000",
     "enabled": "true"
},
```
* Make sure the `"enabled": "true"` field is set to true. If you don't have any mentions of telnet like the above, it will fall under the Server > World > Telnet nodes, you can see an example of this here: https://github.com/Akkadius/EQEmuInstall/blob/master/eqemu_config.json but chances are if you're using Akka's Easy Installer, you had this enabled by default.

### Enabling Players to talk from Discord to EQ
* (Admin-level accounts on Discord can only do the following steps.)
* Inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. The `IGN:` prefix is required for DiscordEQ to detect a player and is used to identify the player in game, For example, to identify the discord user `Xackery` as `Shin`, Create a role named `IGN: Shin`, right click the user Xackery, and assign the role to them.
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`

### Troubleshooting and info
* At this point you should be seeing bidirectional chat functioning, but if not here's some details that may be helpful if it doesn't work first run.
* If you get messages noting Discord is Unauthorized, then your bot is likely not properly authorized via the Visit the link: https://discordapp.com/oauth2/authorize?&client_id={CLIENTID}&scope=bot&permissions=2146958591 changing {CLIENTID} to your client ID. link.
* Firewall doesn't matter, as by default, discordeq's settings will try to connect to localhost 
