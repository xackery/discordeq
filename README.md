# discordeq
This plugin allows Everquest to communicate with Discord in a bidirectional manner.

##How to install
###Make an OOC channel
* Inside Discord, create a channel called #ooc (or whatever you prefer)
* Right click the channel's name, and Copy Link. When you paste it, you'll get a url that looks similar to ``
* You want to note the serverid and channelid #'s and put them into your eqemu_config.xml following up.

###Set up eqemu_config.xml
* Add to eqemu_config.xml these lines:
```xml
<!-- Discord Configuration -->
	<discord>
		<username>YourDiscordUsername</username>
		<password>YourDiscordPassword</password>
		<serverid>ServerIDFromDiscordEQ</serverid>
		<channelid>ChannelIDFromDiscordEQ</channelid>
		<refreshrate>5</refreshrate>
	</discord>
```
* Note: Take a peek at the first line of your eqemu_config.xml, read the `<?xml version="1.0"?>` line, note that if the ?> ending is missing, namely the ?, my exe will not parse your config.

###Enable Player Chat Logging
* Chat Logging allows all player-chat events to be stored in the DB. If you do not have this enabled, run this in a MySQL client:
* ```sql REPLACE ```
* In game, /say #reloadrules
* If you go into your DB client and query the data inside the `qs_player_speech` table after reloadrules, entries should pop up, like OOC events.


###Copy Quest
* You **must** have a static zone and an npc that is unkillable for the next step, but essentially, copy this quest to an NPC to allow discord conversations to be enabled in game:
```perl 
my $lastId = 0;

sub EVENT_SPAWN {
    #Get last ID
    $connect = plugin::LoadMysql();
    $query = "SELECT `id` FROM qs_player_speech ORDER BY `id` DESC LIMIT 1";
    $query_handle = $connect->prepare($query);
    $query_handle->execute();
    while (@row = $query_handle->fetchrow_array()){
        $lastId = $row[0];
    }
      quest::settimer("discord", 1);
}

sub EVENT_TIMER {
      $connect = plugin::LoadMysql();
    $query = "SELECT `from`, `message`, `id` FROM qs_player_speech WHERE `id` > ? AND `type` = 5 AND `to` = '!discord' LIMIT 1";
    $query_handle = $connect->prepare($query);
    $query_handle->execute($lastId);
    while (@row = $query_handle->fetchrow_array()){
        quest::we(2, $row[0]." says from discord, '".$row[1]."'");
        $lastId = $row[2];
    }
    return
}
```
* You will need to #reloadqst and #repop a zone for it to activate.

###Run EXE.
* Run discordeq.exe from the same directory that eqemu_config.xml exists. On success, it will show [ooc] Listening and [discord] Listening.
 

###Enabling Players to talk from Discord to EQ
* Admin-level accounts must grant players the ability to talk in game. 
* To allow this, inside discord go to Server Settings.
* Go to Roles.
* Create a new role, with the name: `IGN: <username>`. IGN: prefix is required for DiscordEQ to detect a player. For example, to identify the discord user `Xackery` as `Shin`, add a role named `IGN: Shin` and assign it to the user `Xackery`
* If the above user chats inside the assigned channel, their message will appear in game as `Shin says from discord, 'Their Message Here'`

## Troubleshooting

### eqemu_config: Error decoding config: XML syntax error on line 131: unexpected EOF
* The first line of your eqemu_config.xml should look like this: `<?xml version="1.0"?>`, note the ?> on ending, if there's no ? on ending it will fail to parse.
 
### Tracing where the problem is
There's 3 parts to make this work:
* 1 is discordeq.exe, it handles all discord chat, read/writing to DB
* 2 is the DB, it uses the qs_player_speech table to provide info for the other two steps
* 3 is the quest script, it handles parsing DB entries that discord puts in (they have a to field of !discord) and broadcasts it in game.
* Trail the steps 1 to 2 to 3 or 3 to 2 to 1 based on which one is faulty.


