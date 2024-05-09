## alpacon token acl

Manages command access for API tokens.

### Synopsis


	The acl command allows you to configure access control for API tokens, specifying which commands can be executed by each token. 
	It supports creating, listing, and modifying ACL rules to fine-tune command execution permissions based on your security requirements.
	

```
alpacon token acl [flags]
```

### Options

```
  -h, --help   help for acl
```

### SEE ALSO

* [alpacon token](alpacon_token.md)	 - Commands to manage api tokens
* [alpacon token acl add](alpacon_token_acl_add.md)	 - Add a new command ACL with specific token and command.
* [alpacon token acl delete](alpacon_token_acl_delete.md)	 - Delete the specified command ACL from an API token.
* [alpacon token acl ls](alpacon_token_acl_ls.md)	 - Display all command ACLs for an API token.

###### Auto generated by spf13/cobra on 23-Apr-2024