# discordeq
This plugin allows Everquest to communicate with Discord in a bidirectional manner.

##How to install
###eqemu_config initial setup
Add to eqemu_config.xml these lines:
```xml
<!-- Discord Configuration -->
	<discord>
		<username>YourDiscordUsername</username>
		<password>YourDiscordPassword</password>
		<serverid>ServerIDFromDiscordEQ</serverid>
		<channelid>ChannelIDFromDiscordEQ</channelid>
	</discord>
```
###Copy and Run Executable
Copy `discordeq.exe` to your eqemu directory where eqemu_config resides. Run it, and it should detect your above settings.

###Configure Step 2: Discord
* If option 1) shows `Status: Good`, then select 2) and configure Discord settings.
* follow the procedure to obtain the server and channel ids, and place them inside your eqemu_config.xml
* Once option 2) shows `Status: Good`, you can proceed to next step.
 
###Configure Step 3: Enable Server Rule for Chat Logging
* If option 3) shows `Status: Bad`, then you may choose to allow discordeq to enable chat logging.
 
###Configure Step 4: Copy Quest File
* If option 4) shows `Status: Bad`, then you may choose to allow discordeq to copy a file into your quests directory.
 
###Start DiscordEQ
* If all options are in good status, option 5) will appear and can be selected to start DiscordEQ.
* If all settings are good on startup of `discordeq.exe`, it will skip the menu and start immediately.
 
###Enabling Players to talk from Discord to EQ
* Admin-level accounts must grant players the ability to talk in game. 
* To allow this, inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. IGN: prefix is required for DiscordEQ to detect a player. For example, to identify the discord user `Xackery` as `Shin`, add a role named `IGN: Shin` and assign it to the user `Xackery`
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`
