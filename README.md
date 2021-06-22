# Tagger

Tagger will cut a new tag for a github repo. It accepts a command in the form:

```irc
bump [major|minor|patch] version for jspc-bots/tagger
```

If there are no tags, the bot starts from '0.0.0'.

This bot does not bother with `rc` tags- it'd need to know which primary field to bump, and then to discard that when un-rc'ing a tag. Feels too complicated. Let rc tags be done manually.

This bot also works for github only- we use the github API, rather than cloning and parsing tags. This is for a couple of reasons:

1. The github API handles enumerating tags and parsing for us
2. We can limit the scope of API tokens better than tokens for cloning/pushing

This bot prepends the lower case character `v` to versions, but will happily parse versions which don't have this. Tagger prepends the `v` purely to make all tagger tags consistent, at the cost of repo consistency.

This bot uses v3 of the github API- the graphql endpoint is pretty cool, but I think the general go client is a bit much for what I'm trying to do- I'm not sure I completely followed all of the nested tags. Prototyping with rest is a little easier.

This bot will not work well where mutliple versions are still receiving patches- consider a project with two active versions:

* v1, and
* v2

Where v2 is the latest and greatest, and v1 still receives bugfixes/ security updates. If the latest version by date is on the v1 line, then the next request to bump a version will also bump the v1 line.

This bot makes a couple of assumptions:

1. You've a SASL account for this bot to use
2. You've enabled actions notifications in github for failed/successful runs

This bot requires the following env vars:

* `$SASL_USER` - the user to connect with
* `$SASL_PASSWORD` - the password to connect with
* `$SERVER` - IRC connection details, as `irc://server:6667` or `ircs://server:6697` (`ircs` implies irc-over-tls)
* `$VERIFY_TLS` - Verify TLS, or sack it off. This is of interest to people, like me, running an ircd on localhost with a self-signed cert. Matches "true" as true, and anything else as false
* `$GITHUB_TOKEN` - Token to use in order to tag repos. Requires at least the `public_repo` scope for tagging public repos. See: https://docs.github.com/en/developers/apps/building-oauth-apps/scopes-for-oauth-apps

The SASL mechanism is hardcoded to PLAIN.

## Building

This bot can be built using pretty standard go tools:

```bash
$ go build
```

Or via docker:

```bash
$ docker build -t foo .
```

## Running

If you've built the app yourself, then happy day- there's your binary!

Otherwise I suggest via docker:

```bash
$ docker build -t foo .
$ docker run foo
```

(Setting the above environment variables accordingly)
