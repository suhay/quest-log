# Quest Log
Modular, serverless game engine utilizing Sandwich Shop.

- [Usage](#usage)
  - [.env file](#env-file)
- [Thread](#thread)
  - [Name](#name)
  - [Hooks](#hooks)
  - [Perspective](#perspective)
  - [Tags](#tags)
- [Entry](#entry)
- [Text Decorators](#text-decorators)
- [Dialog Trees](#dialog-trees)
- [State](#state)
  - [Flags and Pins](#flags-and-pins)
  - [Game](#game)
  - [Bias](#bias)
    - [Anti-Bias](#anti-bias)
- [Development](#development)

## Usage

```bash
$ mkdir quest-log
$ cd quest-log
$ go get github.com/suhay/quest-log
$ go build
```
### .env file
```bash
$ cd quest-log
$ touch .env
```
```
MONGODB_HOST=cluster.mongodb.net
MONGODB_PORT=27017
MONGODB_USERNAME=mongo
MONGODB_PASSWORD=password1234
MONGODB_QUEST_LOG_DB=quest-log
JWT_SECRET=your-jwt-secret

MONGODB_GAME_DB=my-game-mongodb
MONGODB_ENTRY_COLLECTION=entry
MONGODB_THREAD_COLLECTION=thread
```

| Key | Description |
| :-- | :---------- |
| `MONGODB_HOST` | URL to MongoDB host. |
| `MONGODB_PORT` | If using a different port other then the default, especially if db is on `localhost`. |
| `MONGODB_USERNAME` | Username to log into db. |
| `MONGODB_PASSWORD` | Password of the user. |
| `MONGODB_QUEST_LOG_DB` | Quest Log MongoDB name (stores metadata, not game data). |
| `JWT_SECRET` | Secret used by Sandwich Shop to validate request origination. |
| `MONGODB_GAME_DB` | Game database name. |
| `MONGODB_ENTRY_COLLECTION` | Collection name of `entries` types within the Game database. |
| `MONGODB_THREAD_COLLECTION` | Collection name of `thread` types within the Game database. |

## Thread

```json
{
  "name": "Low Power",
  "hooks": [
    {
      "hook": "->torti.Low Power", # continue decorator
      "required" : [
        "{Power 2}", # required expression. this and all siblings must evaluate to true
        {
          "or": [ # an operator based requirement collection. 'or' signifies that at least one child must be true
            "{newStabilizer true}",
            "{newShield true}"
          ]
        }
      ]
    },
    {
      "hook": ["good-vs-evil"], # an array signifies a dynamic pivot to a matching thread or hook
      "trigger": "{bossFight true}" # when the 'bossFight' flag evaluates to true, all sibling hooks are skipped and this one becomes active
    },
    {
      "hook": "That'll do, thanks!"
    }
  ],
  "perspective": "two",
  "tags": ["tutorial"]
}
```
### Name

Human readable name of the entry or thread. Used whenever the thread name appears in written form. These should be unique across all threads and entries.

### Hooks

Hooks are an array of story hooks. Each hook is a self contained part of a bigger story. Hooks are executed in order and only pass onto the next when marked as complete through various means outlined below.

| Key | Description |
| :-- | :---------- |
| `event` | An event happens as soon as this hook becomes active. |
| `hook` | The hook can either be a `string` containing straight flavor text, a continue decorator `->` which will declare where the flavor text can be found, or an `array[string]` which will be used to pivot, dynamically, to another hook or thread. There is a threshold for the number of pivots that may occur to prevent endless threads. |
| `required` | An array of expressions that must be met before the hook can be considered complete. |
| `trigger` | A trigger forces the hook to go next and will skip all other siblings that come before it. These are usually evaluated `flags` or when a certain state is achieved. |

### Perspective

A single value used when dynamically selecting additional hooks. This should represent whichever literary voice the thread or entry was written in.

- `one` - First person
- `two` - Second person
- `three` - Third person

### Tags

An array of strings which describe the events within the thread. This is used when a `hook` is set to an array. All members of that array will be considered when selecting the pivoted hook, but not all will be used in the final decision (see Bias and Anti-bias below).

## Entry

An `entry` is a super-set of a `thread` with the below differences. An entry should be thought of as the base for multiple like threads to branch from.

| Key | Description |
| :-- | :---------- |
| `closing` | An expression that activates after all hooks in the `entry` are either satisfied, skipped, or otherwise completed. |

## Text Decorators

| Decoration | Result |
| :--------- | :----- |
| `->`, `->key` | **Continue** - Go to sibling; either the next non-keyed sibling `->` or sibling with the noted key `->key`. |
| `<-`, `<-key` | **Break** - Go to parent; either return up the tree to the parent and then `continue` to the next non-keyed sibling `<-` or return to the parent and then `continue` the sibling of the noted key `<-key`. |
| `<>`, `<>key` | **Root** - Return to the `root` of the tree and then `continue` to the next sibling `<>` or to the sibling with noted key `<>key`. |
| `{expression value}`, `{exp < 5}` | **Expression** - An evaluated variable plus the `value` it must compare to. Logical operators may be used such as `=`, `<`, or `<=`. If no operator is provided, `=` is assumed. |
| `+flag` | **Flag** - Temporary, named `var` which is removed from `state` after the first time it is broadcast back to Quest Log. These should be transient game states that are read once and then forgotten such as completing a certain task, collecting an item, making a decision, or visiting a certain place. |
| `++pin` | **Pin** - Semi-permanent, named `var`. Like `flags`, but should hold more permanent state data such as the death of a character, the acquisition of certain knowledge such as who the killer is, how much mana a character has, or a change to the world that will affect future threads and entries. These should rarely be deleted. |

## Dialog Trees

```yaml
---
INTRO: # INTRO should always be the first key, this is read out when this file is first triggered
  - |
    O-oh. H-hello there. I wasn't expecting to run into anyone out here.
    I think you saved me. Or -- maybe I saved you...
    A-actually I'm n-not sure it matters come to think of it.
    Oh dear, I detect that I maybe rambling.
    Umm, you -- you can detach me at any time i-if you'd like.
    Just send me out into the void -- all alone, cold -- and alone.
    Or I can just, you know, stick -- around with you.
    It's, I mean, it's completely up to you. You do you, that's what is the most important.
  -
    # an empty continue is shorthand for the next sibling
    - Who are you? ->
    # expressions at the beginning of dialog choices will only appear when true
    - {game.plays >= 1} Ahh! Torti, good to see you! ->I know you
    # a named 'continue' tells the engine to search for a parent
    - Ok, bye! {game.mood -10} ->END
  - |
    W-who, me?
    Oh, I'm just, umm -- I'm just Torti.
    Just little Torti.
    Umm -- hi!
  -
    I know you: 
      - You know me?! How can that be???
      -
        - Because science! ->
        - Because space! ->
        - There's no time to explain! ->
      - Oh, umm, ok?
Low Power: |
  So. I know you're still deciding if you should detach me.
  And that's fine! For reals, if you want me gone, boom, done, cya!
  But, i-if you wouldn't mind -- I could use a little extra power.
  I-if it's not too much trouble!
  Oh geez, I'm sweating I'm so nervous. Or, at -- at least I think I'm sweating.
  Can modules sweat in space?
END: It was nice getting to know you! # END signifies the terminating key in a dialog tree
```

Dialog is stored in `yaml` files for ease of writing. Everything will be stored as `json` one it gets up to MongoDB. Each speaking character should have their own dialog set and all text decorators are valid. Text can be `continued` to from a `hook` by using `fileName.key`.

## State

Quest Log is inherently stateless. This allows it to be deployed across as many serverless environments as you would like. Game state is held within a JWT token to protect against state tampering or cheating. State should be sent along with each request for a `thread` or `entry` and then updated, by the server, with any additions.

```json
{
  "state": {
    "flags": {},
    "pins": {},
    "game": {},
    "bias": {}
  }
}
```

### Flags and Pins

`flags` and `pins` hold the states of those particular types. These can be accessed via expressions such as `{flagName true}` or `{pins.pinName > 10}`. Prefixing with `flags` or `pins` is only needed to prevent name collisions. If a prefix is not provided, all collections will the searched until the first instance is found.

### Game

`game` state is up to the developer as to what is stored. These will be different between games and should hold information for recreating UI assets, such as inventory, and placement of world objects.

### Bias

Each use or reinforcement of a tag will increase its `bias` to be used more heavily while selecting the next thread. When multiple tags are used in a hook pivot, not all tags need to be used, but the more heavily used tags (ones with more bias) should be used ahead of others.

#### Anti-Bias

When a bias is being enforced, there is a possibility that its subsequent anti-bias will be used instead. This randomness is accumulated in a flat value `int` from 1 - 100. Once it is activated, it will take on one of two possibilities:

- The provided tags will flip to become their opposites
- The lesser biased tag will become the strongest and the strongest the weakest when selecting the next pivot.

## Development

- Thread and Entry level flavor descriptions, such as setting and perpetual antagonists.
- Conditional logic in dialog trees.
- Text suppression while selecting dialog replies (more narrative).