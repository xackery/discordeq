# discordeq
This plugin allows Everquest to communicate with Discord in a bidirectional manner.

##How to install

### Note for previous versions
* If you were running a previous version of DiscordEQ, you can now remove the quest file that I had you create, it is no longer needed!

###Set up eqemu_config.xml
* Add to eqemu_config.xml these lines:
```xml
<!-- Discord Configuration -->
	<discord>
		<username>YourDiscordUsername</username>
		<password>YourDiscordPassword</password>
        <telnetusername>YourTelnetAccountLSNAme</telnetusername>
        <telnetpassword>YourTelnetPassword</telnetpassword>        
		<serverid>ServerIDFromDiscord</serverid>
		<channelid>ChannelIDFromDiscord</channelid>
		<refreshrate>5</refreshrate>
        <itemurl>http://yourallaclone.com/alla/?a=item&amp;id=</itemurl>
	</discord>
```

###Set ServerID and ChannelID
* Click the sprocket on the bottom left area to go into user settings.
* Inside User Settings Pop up, go to the Appearance tab on left.
* Inside the Appearance tab, Enable on the Developer Mode option
* Hit done to exit the user settings pop up.
* Inside Discord, create a channel called #ooc (or another name, whichever you prefer)
* Right click the channel's name, and choose the copy link option. When you paste it, you'll get a number that represents a channelid noted above.
* Right click the server's name, and copy link. When you paste it, you'll get a number that represents serverid noted above.

###Set ItemURL
* this is optional, if you leave the itemurl field blank (or remove it), it will default to showing item links in italics in chat, e.g. *Arrow* when someone itemlinks an Arrow.

###Enable Telnet
* Look through eqemu_config and you will find an option for telnet="enabled".

###Set a telnet account password
* Go to your Accounts table on the DB, and set a password for one of your GM accounts. You can manually enter a password in there, plain text, and copy/paste it into the <telnetpassword> field in your eqemu_config.xml.

###Run Discord.
* Run discordeq from the same directory that eqemu_config.xml exists. If any settings are incorrect, you will be prompted on what is incorrect and how to fix it.

###Enabling Players to talk from Discord to EQ
* Admin-level accounts can only do the following steps.
* To allow this, inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. The `IGN:` prefix is required for DiscordEQ to detect a player and is used to identify the player in game, For example, to identify the discord user `Xackery` as `Shin`, I would create a role named `IGN: Shin`, right click the user Xackery, and assign the role to them.
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`
