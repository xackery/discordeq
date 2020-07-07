# DiscordEQ
This plugin allows Everquest to communicate with Discord in a bidirectional manner. Check out [TalkEQ](https://github.com/xackery/talkeq) for an alternative project with additional features, planned to replace discordEQ


## For newbies to github
A break down of the steps you will need to follow:
* Download a file that represents your environment here: https://github.com/xackery/discordeq/releases most likely discordeq-0.51-windows-x64.exe if you aren't sure.
* Copy the downloaded file to your everquest directory. The same place folder world.exe, zone.exe, and eqemu_config.json files are found.

## How to Configure

### Edit eqemu_config
* Add to eqemu_config.json these lines under the server {} section:
```json
{
	"server": {
		"discord" : {
			"channelid" : "11111",
			"itemurl" : "",
			"refreshrate" : "15",
			"clientid": "22222",
			"serverid" : "44444",
			"username" : "55555abc"
		},
		/* ignore this, but other configuration options will be here, like host, chatserver, etc */
        }
}
```
* 

### Prepare App/Bot in Discord
* Go to https://discordapp.com/developers/
* Click My Apps on the top left area.
* Click + New App
* Write anything you wish for the app name, click Create App
* Copy the clientid into your eqemu_config as shown above, changing the clientid 22222 to this new number
* Scroll down to the bot section, and click Create Bot User
* Confirm creating a bot with Yes, Do it!
* Make sure Public Bot is unchecked, as well as oauth2 grant.
* Save, click to reveal the token. Copy the bot token, changing the username 55555abc to this new token
* Copy this link https://discordapp.com/oauth2/authorize?&client_id=22222&scope=bot&permissions=2146958591 , change the 22222 inside the link to your client id you obtained earlier. Visit it, and you will get a discord prompt to authorize the bot to access your server. If you do this step correctly, an offline version of your bot should appear in the members list of your discord server.

### Prepare server, channel IDs
* Inside Discord, right click your server's circular icon and on the bottom choose Copy ID.
* Paste the serverID into your eqemu_config's serverid's 44444 field
* Create a new channel for OOC chat. Right click the channels' name and copy ID.
* Paste the channelID into your eqemu_config's 11111 field

### Prepare itemURL (optional)
* If you have a website that has item links, you can place it into the itemurl field e.g. "http://yoursite.com/item.php?id=". It is assumed your item id's are appended to the end of the url link.
* If you don't have a website with itemlinks, you can safely ignore this step.

### Enable Telnet (optional)
* If you ran Akka's windows installer and did not change the telnet options, you do not have to follow this.
* Inside your eqemu_config.json, look for a section that reads:
```json
"telnet": {
     "ip": "0.0.0.0",
     "port": "9000",
     "enabled": "true"
},
```
* Make sure the `"enabled": "true"` field is set to true. If you don't have any mentions of telnet like the above, it will fall under the Server > World > Telnet nodes, you can see an example of this here: https://github.com/Akkadius/EQEmuInstall/blob/master/eqemu_config.json

### Enabling Players to talk from Discord to EQ
* (Admin-level accounts on Discord can only do the following steps.)
* Inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. The `IGN:` prefix is required for DiscordEQ to detect a player and is used to identify the player in game, For example, to identify the discord user `Xackery` as `Shin`, Create a role named `IGN: Shin`, right click the user Xackery, and assign the role to them.
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`

### Troubleshooting and info
* Your offline bot should go to "online" mode on the member list when you start discordeq. if it isn't, you likely are having an authorization issue (copied the user token incorrectly). Visit the link  again: https://discordapp.com/oauth2/authorize?&client_id={CLIENTID}&scope=bot&permissions=2146958591 changing {CLIENTID} to your client ID.
* At this point you should be seeing bidirectional chat functioning, but if not here's some details that may be helpful if it doesn't work first run.
* Firewall doesn't matter, as by default, discordeq's settings will try to connect to localhost telnet.
