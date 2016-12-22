# cmdotp 
Like Google Authenticator, but on the command line. 

```
$ cmdotp add example-secret PUTTHESECRETHERE
Password:

$ cmdotp add another-secret HEREISANOTHERSUPERSECRET
Password:

$ cmdotp
Password:

example-secret		340271
another-secret		452889
```

## Install

`go get github.com/geotho/cmdotp`

## Usage

### Adding secrets

`$ cmdotp add <NAME> <SECRET KEY>`

You should probably not store this command in your shell history.
A common way of doing this is to have `HISTCONTROL=ignorespace` in your shell, and then prepending your command with a space.
See [here](http://stackoverflow.com/questions/8473121/execute-command-without-keeping-it-in-history) for details.

You'll have to use the same password for all secrets.

Secrets are stored encrypted with your password. By default, they live in `~/.cmdotp/secrets`.

### Viewing OTPs

`$ cmdotp`

Use the same password you used to add them.
