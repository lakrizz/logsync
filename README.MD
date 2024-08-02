

# logsync
Logsync is a little tool that automatically synchronizes a configured set of logseq pages with a static page generator (currently only `hugo`). The mapping is given by a configuration file (see section "*Configuration*")

## How does it do it?
By combining your `logseq` and your `hugo` repositories (in different ways: `logseq` is being pulled/cloned automatically, `hugo` should be an already cloned repository), `ngrok` as a reverse proxy and your `openssh` for authentication against github to create a webhook for your `logseq`repository.

## Installation

Install via `go install github.com/lakrizz/logsync@latest`. Other ways to install them are currently not available.

## Configuration
Place a file called `config.json` in your [`XDG_CONFIG_HOME`](https://wiki.archlinux.org/title/XDG_Base_Directory)  directory, you can find a skeleton in  `/examples/config.json` in this very repository. For `logsync` to work properly, you currently need to create a [GitHub Access Token](https://github.com/settings/tokens) with the following Permissions: `admin:repo_hook, repo`. Additionally you (currently) need to provide an Auth Token for the reverse proxy service [`ngrok`](https://ngrok.com/), if you're already logged in, click [this link](https://dashboard.ngrok.com/tunnels/authtokens). You can replace the Placeholder values in the given `config.json`. The following values need to be set:
- `logseq_repo_url` is the github repository url (other SCM-services are currently not supported) of your `logseq` repository 
- `hugo_repo_path` is the *file system* path of your hugo repository, this path will be used when executing the `hugo` command after an update
- `hugo_exec_params` should be filled with all params that the `hugo` command should be called with (e.g., `--buildDrafts` if you want to include drafts, see [this](https://gohugo.io/commands/hugo/) for available commands)
- `private_key_path` is the path to your ssh key (this will probably be automated soon(tm))
- `private_key_password` needs to be provided if your key is password protected (as we're not making use of your operating system git implementation but a pure go one (which is incredibly awesome, imo))
- `username` is the username that's used for all `logsync` git commits
- `email` is the email that's used for all `logsync` git commits

### Frontmatter
You can add specific frontmatter per mapping, e.g. post-types. Some values are added automatically:
- `title` is the slugged name of the input file for each mapping (includes files added by the `recursive` option)
- `date` is set to `time.Now().Format(time.RFC3339)`

> **Note**: Frontmatter Entries with the key `title` or `date` in your config will be omitted 

### Options
You can choose of the following Options for your Mapping.

#### Remove Internal Links

This option removes all internal links from mapped `logseq` pages to prevent invalid links in your target static page generator.
```json
{
    "mappings": [
        {
            "options": {
                "remove_internal_links": true
            }
        }
    ]
}
```

> **Note**: This will interfere with the *recursion* Option.

#### Remove Empty Trails

This option, if enabled, will remove empty list items from the tail of the mapped page. 
```json
{
    "mappings": [
        {
            "options": {
                "remove_empty_trails": true
            }
        }
    ]
}
```

#### Recursion
> **Attention**: This option is currently not available 
This option, if enabled, will also map `n` levels of links (given by the option `recursion_depth` [defaults to 1]), counting from the given `source` page, if this option is enabled, it requires a `recursion_target`, which needs to be a directory. Additionally you can set `recursion_skip_source` to `true` if you don't want the `source` page to be mapped to your target static page generator (e.g., for blog posts), this setting defaults to `false`

```json
{
    "mappings": [
        {
            "options": {
                "recursive": true,
                "recursion_depth": 1,
                "recursion_target": "content/books/notes/",
                "recursion_skip_source": false
            }
        }
    ]
}
```
For each file added with this method a slug based on the input page name is created and used for Hugo.

#### Unindent First Level

This option unindents the first level of a mapped logseq page, effectively converting each bullet point into a paragraph. This is useful for rich texts (e.g., for blog posts). If there's no indentation levels, this option is a NOP.

```json
{
    "mappings": [
        {
            "options": {
                "unindent_first_level": true,
            }
        }
    ]
}
```

#### Include Attachments
This option toggles whether logsync should copy all visible attachments of a mapping. 


```json
{
    "mappings": [
        {
            "options": {
                "include_attachments": true,
            }
        }
    ]
}
```

# Contributing
Feel free to create Pull Requests. I'm happy for anyone to improve this little tool. You can also open (or work on) Issues here on GitHub. <3
